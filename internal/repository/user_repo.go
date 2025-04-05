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
  UpdateProfile(profile *domain.Profile) error
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
  query:="INSERT INTO profiles(id,name,phone_number,bio,profile_picture_url) VALUES ($1,$2,$3,$4,$5)"
  _,err := r.db.Exec(query,&profile.Id,&profile.Name,&profile.Phone_number,&profile.Bio,&profile.Profile_Picture_Url)
  if err != nil{
    return errors.New("failed to create profile"+err.Error())
  }
  return nil
}

func (r* userRepository) GetProfile(phone_number ,string)(*domain.Profile,error){
  var profile domain.Profile

  query := "SELECT * FROM profiles WHERE phone_number=$1"
  row:= r.db.QueryRow(query,phone_number)
  err := row.Scan(&profile.Id,&profile.Name,&profile.Phone_Number,&profile.Bio,
        &profile.Profile_Picture_Url,&profile.CreatedAt,&profile.UpdatedAt)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil,errors.New("profile not found")
    }
    return nil,errors.New("failed to fetch the profile of the user: "+err.Error())
  }
  return &profile,nil

}

func (r* userRepository) UpdateProfile(profile *domain.UpdateProfile) error{
  fields := []string{}
  values := []interface{}
  paramIndex := 1

  if profile.Name != nil{
    fields = append(fields,fmt.Sprintf("name=$%d",paramIndex))
    values = append(values,*profile.Name)
    paramIndex++
  }
  if profile.Bio != nil{
    fields := append(fields,fmt.Sprintf("bio=$%d",paramIndex))
    values = append(values,*profile.Bio)
    paramIndex++
  }
  if profile.ProfilePictureUrl != nil{
    fields = append(fields,fmt.Sprintf("profile_picture_url=$%d",paramIndex))
    values = append(values,*profile.ProfilePictureUrl)
    paramIndex++
  }

  if len(fields) == 0 {
    return error.New("No fields to be update")
  }
  fields = append(fields, fmt.Sprintf("updated_at=$%d", paramIndex))
	values = append(values, time.Now())
	paramIndex++
  
  values = append(values,profile.id)

  query := fmt.Sprintf("UPDATE profiles SET %s WHERE id=$%d",strings.Join(fields,","),paramIndex)
  _,err := r.db.Exec(query,values...)
  if err != nil{
    return  fmt.Errorf("failed to update profile :%w",err)
  }
  return nil
}
