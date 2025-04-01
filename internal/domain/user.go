package domain

import (
  "time"
)

type User struct {
  Id      uint      `gorm:"primaryKey;autoIncrement"`
  Name    string    `gorm:"size:100;not null"`
  Email   string    `gorm:"unique;not null"`
  Password  string  `gorm:"not null"`
  CreatedAt time.Time
  UpdatedAt time.Time
}
