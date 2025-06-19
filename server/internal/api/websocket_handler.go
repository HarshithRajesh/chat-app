package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HarshithRajesh/app-chat/internal/domain"
	"github.com/HarshithRajesh/app-chat/internal/realtime"
	"github.com/HarshithRajesh/app-chat/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type IncomingMessage struct {
	SenderId   string `json:"sender_id"`
	RecieverID string `json:"reciever_id"`
	Content    string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsChatHandler struct {
	Hub         realtime.IHub
	UserService service.UserService
	ChatService service.ChatService
}

func NewWsChatHandler(hub realtime.IHub, userService service.UserService, chatService service.ChatService) *WsChatHandler {
	return &WsChatHandler{
		Hub:         hub,
		UserService: userService,
		ChatService: chatService,
	}
}

func (h *WsChatHandler) HandleWebSocket(c *gin.Context) {
	log.Printf("Incoming websocket connection request from %s", c.Request.RemoteAddr)

	// Get user_id from query parameters
	userID := c.Query("user_id")
	if userID == "" {
		log.Printf("ERROR: Websocket connection denied for %s: Missing UserID in query", c.Request.RemoteAddr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: UserId required"})
		return
	}

	log.Printf("Authenticated Websocket connection for UserId: %s from %s", userID, c.Request.RemoteAddr)

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error: Failed to upgrade connection from http to websocket: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to websocket"})
		return
	}

	client := &realtime.Client{
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	h.Hub.RegisterClient(client)

	defer func() {
		h.Hub.UnregisterClient(client)
		log.Printf("Websocket connection closed for UserID: %s", userID)
	}()

	go h.writePump(client)
	h.readPump(client)
}

func (h *WsChatHandler) readPump(client *realtime.Client) {
	client.Conn.SetReadLimit(512)
	readDeadLine := 60 * time.Second

	client.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
		return nil
	})

	for {
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Client %s (UserID: %s) closed Websocket connection gracefully", client.Conn.RemoteAddr(), client.UserID)
			} else {
				log.Printf("ERROR: Read error from Websocket client %v", err)
			}
			break
		}

		log.Printf("Received message from client %s (UserID: %s, Type: %d): %s", client.Conn.RemoteAddr(), client.UserID, messageType, message)

		if messageType == websocket.TextMessage {
			var incomingMsg IncomingMessage
			if jsonErr := json.Unmarshal(message, &incomingMsg); jsonErr != nil {
				log.Printf("Error: Failed to Unmarshal incoming chat message from UserID %s: %v. Message %s", client.UserID, jsonErr, message)
				continue
			}

			if incomingMsg.SenderId != client.UserID {
				log.Printf("Alert: Incoming message from different user sender id: %s, userid: %s", incomingMsg.SenderId, client.UserID)
				continue
			}

			user_id, err := strconv.ParseUint(client.UserID, 10, 64)
			if err != nil {
				log.Printf("Error parsing user_id: %v", err)
				continue
			}

			sender_id, err := strconv.ParseUint(incomingMsg.SenderId, 10, 64)
			if err != nil {
				log.Printf("Error parsing sender_id: %v", err)
				continue
			}

			receiver_id, err := strconv.ParseUint(incomingMsg.RecieverID, 10, 64)
			if err != nil {
				log.Printf("Error parsing receiver_id: %v", err)
				continue
			}

			msg := domain.Message{
				Id:         uint(user_id),
				SenderId:   uint(sender_id),
				ReceiverId: uint(receiver_id),
				Content:    incomingMsg.Content,
			}

			chatErr := h.ChatService.SendMessage(&msg)
			if chatErr != nil {
				log.Printf("Error: Failed to send the message from chat service for UserID %s: %v", client.UserID, chatErr)
			} else {
				log.Printf("Message from UserID %d to %d successfully pushed to Redis Stream via WebSocket", sender_id, receiver_id)
			}
		}
	}
}

func (h *WsChatHandler) writePump(client *realtime.Client) {
	pingPeriod := (60 * time.Second * 9) / 10
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			// Setting deadline for slow writes
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error: Failed to get Websocket writer for client %s (UserID: %s): %v", client.Conn.RemoteAddr(), client.UserID, err)
				return
			}

			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("Error: Failed to close websocket writer: %v", err)
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error: Ping error for the client: %v", err)
				return
			}
		}
	}
}
