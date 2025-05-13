package api

import (
  // "encoding/json"
  "net/http"
  // "io/ioutil"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/service"
  "strconv"
  "encoding/json"
)

type ChatHandler struct{
  chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler{
  return &ChatHandler{chatService}
}

func(h *ChatHandler) SendMessage(w http.ResponseWriter,r *http.Request){
  if r.Method != http.MethodPost{
    http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
    return
  }

  var req domain.Message
  err := json.NewDecoder(r.Body).Decode(&req)
  if err != nil{
    http.Error(w,"Invalid request body",http.StatusBadRequest)
    return
  }

  if err := h.chatService.SendMessage(&req);err != nil{
    http.Error(w,err.Error(),http.StatusBadRequest)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Message Sent"))
  
}

func (h *ChatHandler) GetMessage(w http.ResponseWriter,r *http.Request){
  if r.Method != http.MethodGet{
    http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
    return
  }
  user1Str := r.URL.Query().Get("user1")
  user2Str := r.URL.Query().Get("user2")
  if user1Str== "" || user2Str == ""{
    http.Error(w,"users can not be empty",http.StatusBadRequest)
    return
  }
  num1,err := strconv.ParseUint(user1Str,10,64)
  if err != nil{
    http.Error(w,"Invalid user1 id",http.StatusBadRequest)
    return
  }
  num2,err := strconv.ParseUint(user2Str,10,64)
  if err != nil {
    http.Error(w,"Invalid user2 id",http.StatusBadRequest)
    return
  }
  user1 := uint(num1)
  user2 := uint(num2)
  messages,err := h.chatService.GetMessage(user1,user2)
  if err != nil{
    http.Error(w,err.Error(),http.StatusBadRequest)
    return
  }
  json.NewEncoder(w).Encode(messages)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("History"))
  
}
