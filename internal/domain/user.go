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

type Profile struct {
  Id                    uint      `gorm:"primaryKey;refrences:Id"`
  Name                  string    `gorm:"size:100"`
  Phone_Number          string    
  Bio                   string
  Profile_Picture_Url   string
  CreatedAt             time.Time
  UpdatedAt             time.Time  
}

type UpdateProfile struct{
  Id                uint
  Name              *string
  Bio               *string
  ProfilePictureUrl *string
}

type Contact struct {
  Id          uint
  UserId     uint
  ContactId  uint
}
