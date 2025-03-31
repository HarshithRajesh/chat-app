package service

import (
  "errors"
  "github.com/HarshithRajesh/chat-app/domain"
  "github.com/HarshithRajesh/chat-app/repository"
)

type UserService interface {
  SignUp(user *domain.User) error 
}

type userService struct {
  repo repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
  return &userService{repo}
}

func (s *userService) Signup(user *domain.User) error {
  existingUser ,_ := s.repo.GetUserByEmail(user.Email)
  if existingUser != nil{
    returnm errors.New("email already registered")
  }

  user.Password = "hashed_"+user.Password
  return s.repo.CreateUser(user)
}
