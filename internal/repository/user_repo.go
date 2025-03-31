package repository

import (
  "database/sql"
  "errors"
  "github.com/HarshtihRajesh/chat-app/internal/domain"
)

type UserRepository interface {
  CreateUser(user *domain.User) error 
  GetUserByEmail(email string)(*domain.User,error)
}

type userRepository struct{
  db *sql.DB 
}

func NewUserRepository (db *sql.DB) UserRepository{
  return &userRepository{db}
}

func (r *userRepository) CreateUser(user *domain.User) error{
  query := "INSERT INTO users(name,email,password) VALUES ($1,$2,$3)"
  _,err := r.db.Exec(query,user.Name,user.Email,user.Password)
  if err != nil{
    return errors.New("failed to create user: "+err.Error())
  }
  return nil 
}

func (r *userRepository) GetUserByEmail(email string)(*domain.User,error){
  var user domain.User
  query := "SELECT * FROM users WHERE email=$1"
  row := r.db.QueryRow(query,email)
  err := row.Scan()
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, errors.New("user not found")
    }
    return nil, errors.New("failed to fetch user: "+err.Error())
  }
  return &user, nil
}

