package domain

import (
  "time"
)

type Message struct{
  Id            uint 
  SenderId      uint
  ReceiverId    uint
  Content       string
  CreatedAt     time.Time
}
