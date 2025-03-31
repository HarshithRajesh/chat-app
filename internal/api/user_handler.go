package api

import (
  "encoding/json"
  "net/http"
  "fmt"
  "log"
  "github.com/HarshithRajesh/chat-app/internal/domain"
  "github.com/HarshithRajesh/chat-app/internal/service"
)


type UserHandler struct {
  userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler{
  return &UserHandler{userService}
}

func (h *UserHandler) SignUp(w http.ResponseWriter,r *http.Request){
  if r.Method != http.MethodPost{
    http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
    return 
  }

  body, _ := ioutil.ReadAll(r.Body)
  var user domain.User 
  if err := json.Unmarshal(body,&user); err != nil{
    http.Error(w,"Invalid JSON body",http.StatusBadRequest)
    return
  }

  if err := h.userService.SignUp(&user); err != nil{
    http.Error(w,err.Error(),http.StatusBadRequest)
    return 
  }

  w.WriteHeader(http.StatusCreated)
  w.Write([]byte("User registered successfully"))

}
