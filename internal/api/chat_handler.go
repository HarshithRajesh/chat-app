import api

import(
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/service"
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

  if err := h.chatService.SendMessage(msg);err != nil{
    http.Error(w,err.Error(),http.StatusBadRequest)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Message Sent"))
  
}
