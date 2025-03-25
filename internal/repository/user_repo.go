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

func NewUserRepository*db *gorm.DB UserRepository{
  return &userRepository{db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
  return r.db.Create(user).Error 
}

func (r *userRepository) GetUser(email string)(*domain.User,error){
  var user domain.User 
  err := r.db.Where('email =?',email).First(&user).Error 
  if err != nil {
    return nil,err 
  } 
  return &user,nil 
}
