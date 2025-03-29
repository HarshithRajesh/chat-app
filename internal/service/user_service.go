// package service
//
// import (
//   "errors"
//   "github.com/HarshithRajesh/chat-app/internal/domain"
//   "github.com/HarshithRajesh/chat-app/internal/repository"
//   "github.com/HarshithRajesh/chat-app/internal/common"
// )
//
// type UserService interface{
//   Signup(user *domain.User) error 
//   Login(email, password string)(string,error)
// }
//
// type userService struct{
//   repo repository.UserRepository
// }
//
// func NewUserService(repo repository.UserRepository) UserService{
//   return &userService{repo}
// }
//
// func (s *userService) Signup(user *domain.User) error {
// 	hashedPassword, err := common.HashPassword(user.Password)
// 	if err != nil {
// 		return err
// 	}
// 	user.Password = hashedPassword
// 	return s.repo.CreateUser(user)
// }
//
// func (s *userService) Login(email, password string) (string, error) {
// 	user, err := s.repository.GetUserByEmail(email)
// 	if err != nil {
// 		return "", errors.New("invalid credentials")
// 	}
// 	if !common.CheckPasswordHash(password, user.Password) {
// 		return "", errors.New("invalid credentials")
// 	}
// 	return common.GenerateJWT(user.ID, user.Email)
// }
