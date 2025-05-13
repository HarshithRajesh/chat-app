package service

import (
  // "errors"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/repository"
  // "github.com/HarshithRajesh/app-chat/internal/config"
  "context"
  "github.com/redis/go-redis/v9"
  "fmt"

)

type ChatService interface {
  SendMessage(msg *domain.Message)error
  GetMessage(user1,user2 uint)([]domain.Message,error)
}

type chatService struct{
 repo repository.ChatRepository 
 redisClient *redis.Client
}

func NewChatService(repo repository.ChatRepository,redisClient *redis.Client) ChatService{
  return &chatService{
    repo:repo,
    redisClient : redisClient,
  }
}

func (s *chatService) SendMessage(msg *domain.Message)error{

  err:=s.repo.SaveMessage(msg)
  if err != nil{
    return err
  }
  _,err = s.redisClient.XAdd(context.Background(),&redis.XAddArgs{
    Stream : "chat_stream",
    Values: map[string]interface{}{
      "sender_id":  msg.SenderId,
      "receiver_id": msg.ReceiverId,
      "message":  msg.Content,
    },
  }).Result()

  if err != nil{
    return fmt.Errorf("Failed to publish message to redis stream: %w",err)
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
