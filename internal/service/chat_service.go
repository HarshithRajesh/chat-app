package service

import (
  // "errors"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/repository"
)

type ChatService interface {
  SendMessage(msg *domain.Message)error
  GetMessage(user1,user2 uint)([]domain.Message,error)
}

type chatService struct{
 repo repository.ChatRepository 
}

func NewChatService(repo repository.ChatRepository) ChatService{
  return &chatService{repo}
}

func (s *chatService) SendMessage(msg *domain.Message)error{
  err:=s.repo.SaveMessage(msg)
  if err != nil{
    return err
  }
  return nil
}

func (s *chatService) GetMessage(user1,user2 uint)([]domain.Message,error){
  messages,err := s.repo.GetMessage(user1,user2)
  if err != nil{
    return nil,err
  }
  return messages,nil
}
