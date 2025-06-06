package api 

import (
  "context"
  "log"
  "net/http"
  "github.com/gorilla/websocket"
  "github.com/HarshithRajesh/app-chat/internal/realtime"
  "github.com/HarshithRajesh/app-chat/internal/service"
) 

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
}

func NewWsChatHandler(hub realtime.IHub, userService service.UserService) *WsChatHandler{
  return &WsChatHandler{
    Hub:  hub,
    UserService:  userService,
  }
}

func (h *WsChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request){
  log.Printf("Incoming websocket connection request from %s",r.RemoteAddr())
  //future reference add JWT token or session cookie
  userID := r.URL.Query.Get("user_id")
}
