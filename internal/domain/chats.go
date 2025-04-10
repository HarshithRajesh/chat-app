package domain

import (
  "time"
)

type Message struct{
  Id            uint 
  SenderId      uint    `json:"sender_id"`
  ReceiverId    uint     `json:"receiver_id"`
  Content       string   `json:"content"`
  CreatedAt     time.Time
}
