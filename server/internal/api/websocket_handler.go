package api 

import (
  // "context"
  "strconv"
  "log"
  "net/http"
  "github.com/gorilla/websocket"
  "github.com/HarshithRajesh/app-chat/internal/realtime"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/service"
  "time"
  "encoding/json"
) 
type IncomingMessage struct{
  SenderId  string  `json:"sender_id"`
  RecieverID string `json:"reciever_id"`
  Content   string   `json:"content"`
}
var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,

  CheckOrigin: func(r *http.Request) bool{
    return true
  },
}

type WsChatHandler struct{
  Hub realtime.IHub
  UserService service.UserService
  ChatService service.ChatService
}

func NewWsChatHandler(hub realtime.IHub, userService service.UserService,chatService service.ChatService) *WsChatHandler{
  return &WsChatHandler{
    Hub:  hub,
    UserService:  userService,
    ChatService:  chatService,
  }
}

func (h *WsChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request){

  log.Printf("Incoming websocket connection request from %s",r.RemoteAddr)
  //future reference add JWT token or session cookie
  userID := r.URL.Query().Get("user_id")
  //upgrade to jwt token verification if there is error then invalid token
  if userID == ""{
    log.Printf("ERROR: Websocket connection denied for %s: Missing UserID in query")
    http.Error(w,"Unauthorized: UserId required",http.StatusUnauthorized)
    return
  }
  log.Printf("Authenticated Websocket connection for UserId: %s from %s",userID,r.RemoteAddr)

  conn,err := upgrader.Upgrade(w,r,nil)
  if err != nil{
    log.Printf("Error: Failed to upgrade connection from http to websocket")
    return
  }

  client := &realtime.Client{
    Conn : conn,
    Send: make(chan []byte,256),
    UserID: userID,
  }

  h.Hub.RegisterClient(client)

  defer func(){
    h.Hub.UnregisterClient(client)
    log.Printf("Websocket connection closed for for UserID: %s",userID)
  }()
  go h.writePump(client)
  h.readPump(client)
}

func(h *WsChatHandler) readPump(client *realtime.Client){

  client.Conn.SetReadLimit(512)
  readDeadLine := 60 * time.Second

  client.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
  client.Conn.SetPongHandler(func(string)error{
    client.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
    return nil
  })

  for {
    messageType,message,err := client.Conn.ReadMessage()
    if err != nil{
      if websocket.IsCloseError(err,websocket.CloseGoingAway,websocket.CloseNormalClosure){
        log.Printf("Client %s (UserID: %s) closed Websocket connection gracefully.")
      } else{
        log.Printf("ERROR: Read error from Websocket client %v",err)
      }
      break
    }
    log.Printf("Recieved message from client %s (UserID: %s ,Type: %d): %s",client.Conn.RemoteAddr,client.UserID,messageType,message)

    if messageType == websocket.TextMessage{
      var incomingMsg IncomingMessage
      if jsonErr := json.Unmarshal(message,&incomingMsg); jsonErr != nil{
        log.Printf("Error : Failed to Unmarshal incoming chat message from USerID %s : %v.Message %s",client.UserID,jsonErr,message)
        continue
      }

      if incomingMsg.SenderId != client.UserID{
        log.Printf("ALert: Incoming message from different user sender id : %s , userid : %s",incomingMsg.SenderId,client.UserID)
        continue
      }
      // h.chatService.SendMessage(message)

      // sendCtx,cancel := context.WithTimeout(context.Background(),5*time.Second)
      // defer cancel()
      user_id,err := strconv.ParseUint(client.UserID,10,64)
      if err != nil{
        log.Println(err)
      }
      msg:= domain.Message{
        Id: uint(user_id),
        SenderId: incomingMsg.SenderId,
        RecieverId:incomingMsg.RecieverId,
        Content:incomingMsg.Content,
      }

      chatErr := h.ChatService.SendMessage(msg)
      if chatErr != nil{
        log.Printf("Error: Failed to send the message from chat service for user USerID %s : %v",client.UserID,err)
      } else{
        log.Printf("Message from UserID %s to %s successfully pushed to Redis Stream via WebSocket.", incomingMsg.SenderID, incomingMsg.ReceiverID)
      }
      
    }
  }
 }

func (h *WsChatHandler) writePump(client *realtime.Client){

  pingPeriod := (60 *time.Second *9)/10
  ticker := time.NewTicker(pingPeriod)

  defer func(){
    ticker.Stop()
    client.Conn.Close()
  }()

  for{
    select {
    case message,ok := <-client.Send:
      //setting deadline for slow wirtes
      client.Conn.SetWriteDeadline(time.Now().Add(10*time.Second))
      if !ok{
        client.Conn.WriteMessage(websocket.CloseMesssage,[]byte{})
        return
      }

      w,err := client.Conn.NextWriter(websocket.TextMessage)
      if err != nil{
        log.Printf("Error: Failed to get Websocket writer for client %s (UserID : %s):%v",client.Conn.RemoteAddr(),client.UserID,err)
        return
      }

      w.Write(message)

      n:= len(client.Send)
      for i:=0;i<n;i++{
        w.Write([]byte{'\n'})
        w.Write(<-client.Send)
      }

      if err := w.Close();err != nil{
        log.Println("Error: Failed to close websocket writer")
        return
      }
    case <-ticker.C:
      client.Conn.SetWriteDeadline(time.Now().Add(10*time.Second))
      if err := client.Conn.WriteMessage(websocket.PingMessage,nil);err != nil{
        log.Println("Error: Ping error for the client")
        return
      }
      
    }
  }
}
