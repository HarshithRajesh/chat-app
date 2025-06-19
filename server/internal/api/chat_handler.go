package api

import (
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/service"
)

type ChatHandler struct{
  chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler{
  return &ChatHandler{chatService}
}

func(h *ChatHandler) SendMessage(c *gin.Context){
  var req domain.Message
  if err := c.ShouldBindJSON(&req); err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
    return
  }

  if err := h.chatService.SendMessage(&req); err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  c.JSON(http.StatusOK, gin.H{"message": "Message Sent"})
}

func (h *ChatHandler) GetMessage(c *gin.Context){
  user1Str := c.Query("user1")
  user2Str := c.Query("user2")
  
  if user1Str == "" || user2Str == ""{
    c.JSON(http.StatusBadRequest, gin.H{"error": "users can not be empty"})
    return
  }
  
  num1, err := strconv.ParseUint(user1Str, 10, 64)
  if err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1 id"})
    return
  }
  
  num2, err := strconv.ParseUint(user2Str, 10, 64)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user2 id"})
    return
  }
  
  user1 := uint(num1)
  user2 := uint(num2)
  
  messages, err := h.chatService.GetMessage(user1, user2)
  if err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  c.JSON(http.StatusOK, gin.H{
    "message": "History",
    "data": messages,
  })
}
