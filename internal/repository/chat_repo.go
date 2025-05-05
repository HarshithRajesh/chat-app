package repository

import (
  "database/sql"
  "errors"
  // "log"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/redis/go-redis/v9"
  // "fmt"
  "context"
 )

type ChatRepository interface{
   SaveMessage(msg *domain.Message) error
   GetMessage(user1,user2 uint)([]domain.Message,error)
   // ReadMessageFromStream(ctx context.Context,redisClient *redis.Client,streamName string,
   //                          startId string,count int64)([]redis.XStream,error)
}

type chatRepository struct{
  db *sql.DB 

}

func NewChatRepository (db *sql.DB) ChatRepository{
  return &chatRepository{db}
}

func (r *chatRepository) SaveMessage(msg *domain.Message)error{

// log.Printf("Sending message from %d to %d", msg.SenderId, msg.ReceiverId)
  query := "INSERT INTO messages(sender_id,receiver_id,content) VALUES ($1,$2,$3)"
  _,err := r.db.Exec(query,msg.SenderId,msg.ReceiverId,msg.Content)
  if err != nil{
    return errors.New("Failed to send the message ->"+err.Error())
  }
  
  return nil
  
}

func (r *chatRepository) GetMessage(user1,user2 uint)([]domain.Message,error){
  query := `SELECT id,sender_id,receiver_id,content,created_at FROM messages
            WHERE (sender_id = $1 AND receiver_id=$2) or (sender_id=$2 AND receiver_id=$1)
            ORDER BY created_at ASC`
  rows,err := r.db.Query(query,user1,user2)
  if err != nil{
    return nil,err
  }
  defer rows.Close()

  var messages []domain.Message
  for rows.Next(){
    var message domain.Message
    err := rows.Scan(
      &message.Id,
      &message.SenderId,
      &message.ReceiverId,
      &message.Content,
      &message.CreatedAt,
    )
    if err != nil{
      return nil,err
    }
    messages = append(messages,message)
  }
  return messages,nil
}

func ReadMessageFromStream(ctx context.Context,redisClient *redis.Client,streamName string,
                            startID string,count int64)([]redis.XStream,error){
  res,err:= redisClient.XRead(ctx,&redis.XReadArgs{
    Streams: []string{streamName,startID},
    Count: count,
  }).Result()

  if err != nil{
    if err == redis.Nil{
      return []redis.XStream{},nil
    }
    return nil,err
  }
  
  return res,nil
}
