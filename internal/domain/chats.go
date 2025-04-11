package domain

import (
  "time"
)

type Message struct{
  Id            uint    `json:"id"` 
  SenderId      uint    `json:"sender_id"`
  ReceiverId    uint     `json:"receiver_id"`
  Content       string   `json:"content"`
  CreatedAt     time.Time `json:created_at"`
}
