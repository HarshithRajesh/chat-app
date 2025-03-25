package repository

import (
  "github.com/HarshithRajesh/chat-app/internal/domain"
  "gorm.io/gorm"
)

type UserRepository interface {
  SignUp(user *domain.User) error
  GetUser(email string)(*domain.User,error)
}

type userRepository struct {
  db *gorm.DB
}

