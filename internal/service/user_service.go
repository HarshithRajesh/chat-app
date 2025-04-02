package service

import (
  "errors"
  "github.com/HarshithRajesh/app-chat/internal/domain"
  "github.com/HarshithRajesh/app-chat/internal/repository"
)
type Response struct{
  Message string `json:"message"`
}
type UserService interface {
  SignUp(user *domain.User) error 
  Login(user *domain.User) error
}

type userService struct {
  repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
  return &userService{repo}
}

func (s *userService) SignUp(user *domain.User) error {
  existingUser ,_ := s.repo.GetUserByEmail(user.Email)
  if existingUser != nil{
    return errors.New("email already registered")
  }

  // user.Password = "hashed_"+user.Password
  return s.repo.CreateUser(user)
}

func (s *userService) Login(user *domain.User)error{
  log,err := s.repo.LoginCheck(user.Email,user.Password)
  if err != nil{
    return err
  }
  if log.Password != user.Password{
    return errors.New("Password didnt match")
  }
  *user = *log
  return nil
}
