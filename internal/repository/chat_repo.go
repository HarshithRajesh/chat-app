package repository

import (
  "database/sql"
  "errors"
  "github.com/HarshithRajesh/app-chat/internal/domain"
)

type ChatRepository interface{
   SaveMessage(msg *domain.Message) error
   GetMessage(user1,user2 uint)([]domain.Message,error)
}

type chatRepository struct{
  db *sql.DB 
}

func NewChatRepository (db *sql.DB) ChatRepository{
  return &chatRepository{db}
}

func (r *chatRepository) SaveMessage(msg *domain.Message)error{
  query := "INSERT INTO messages(sender_id,receiver_id,content) VALUES ($1,$2,$3)"
  _,err := r.db.Exec(query,msg.SenderId,msg.ReceiverID,Content)
  if err != nil{
    return errors.New("Failed to send the message"+err.Error())
  }
  return nil
}

func (r *chatRepository) GetMessage(user1,user2 uint)([]domain.Message,error){
  query := "SELECT content FROM messages WHERE sender_id = $1 and receiver_id=$2"
  rows,err := r.db.Query(query,user1,user2)
  if err != nil{
    return nil,err
  }
  defer rows.Close()

  var messages []domain.Message
  for rows.Next(){
    var message domain.Message
    err := rows.Scan(
      &message.Content
    )
    if err != nil{
      return nil,err
    }
    messages = append(messages,message)
  }
  return messages,nil
}
