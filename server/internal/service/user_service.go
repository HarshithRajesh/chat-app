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
  Profile(profile *domain.Profile) error
  Contact(user_id uint,phone string) error
  ViewContactList(user_id uint)([]domain.Profile,error)
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

func (s *userService) Profile(profile *domain.Profile)error{
  _,err := s.repo.GetProfile(profile.Id)
  if err != nil{
    if err.Error() == "profile not found"{
      return s.repo.CreateProfile(profile)
    }
    return err 
  }
  updateInput := &domain.UpdateProfile{
    Id : profile.Id,
  }
  if profile.Name != ""{
    updateInput.Name = &profile.Name
  }
  if profile.Bio != ""{
    updateInput.Bio = &profile.Bio 
  }
  if profile.Profile_Picture_Url != ""{
    updateInput.ProfilePictureUrl = &profile.Profile_Picture_Url
  }
  return s.repo.UpdateProfile(updateInput)

}

func (s *userService) Contact(user_id uint,phone string) error{
  contact_id,err := s.repo.GetContact(phone)
  if err != nil {
    return err 
  }
  
  if user_id == contact_id{
    return errors.New("cant add userself")
  }

  exists,err := s.repo.ContactAlreadyAdded(user_id,contact_id)
  if err != nil{
    return err 
  }
  if exists{
    return errors.New("contact already exists")
  }

  return s.repo.CreateContact(user_id,contact_id)
}

func (s *userService) ViewContactList(user_id uint)([]domain.Profile,error){
  contact_ids,err := s.repo.GetContactId(user_id)
  if err != nil{
    return nil,err 
  }
  profile,err := s.repo.GetProfileId(contact_ids)
  if err != nil{
    return nil,err 
  }
  return profile,err

}
