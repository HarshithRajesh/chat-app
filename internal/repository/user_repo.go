package repository

import (
  "database/sql"
  "errors"
  "github.com/HarshithRajesh/app-chat/internal/domain"
)

type UserRepository interface {
  CreateUser(user *domain.User) error 
  GetUserByEmail(email string)(*domain.User,error)
  LoginCheck(email string,password string)(*domain.User,error)
  CreateProfile(profile *domain.Profile) error
  GetProfile(phone_number string)(*domain.Profile,error)
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
  err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, errors.New("user not found")
    }
    return nil, errors.New("failed to fetch user: "+err.Error())
  }
  return &user, nil
}

func (r* userRepository) LoginCheck(email string, password string)(*domain.User,error){
  var user domain.User
  query := "SELECT * FROM users WHERE email=$1"
  row := r.db.QueryRow(query,email)
  err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, errors.New("user not found")
    }
    return nil, errors.New("failed to fetch user: "+err.Error())
  }
    return &user,nil

}

func (r* userRepository) CreateProfile(profile *domain.Profile) error{
  query:="INSERT INTO profiles(id,name,phone_number,profile_picture_url) VALUES ($1,$2,$3,$4)"
  _,err := r.db.Exec(query,profile.Id,profile.Name,profile.Phone_number,profile.Profile_Picture_Url)
  if err != nil{
    return errors.New("failed to create profile"+err.Error())
  }
  return nil
}

func (r* userRepository) GetProfile(phone_number)
