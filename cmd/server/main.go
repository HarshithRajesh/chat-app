package main 

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "github.com/HarshithRajesh/app-chat/internal/api"
  "github.com/HarshithRajesh/app-chat/internal/config"
  "github.com/HarshithRajesh/app-chat/internal/repository"
  "github.com/HarshithRajesh/app-chat/internal/service"
)
type Response struct{
  Message string `json:"message"`
}
func health(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  response := Response{Message: "Hi Welcome to Chaat"}
  json.NewEncoder(w).Encode(response)
}

func handler(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w,"Hi,there, Welocome to my chaat ")
}

func main(){
  db := config.ConnectDB()
  userRepo := repository.NewUserRepository(db)
  userService := service.NewUserService(userRepo)
  userHandler := api.NewUserHandler(userService)

  http.HandleFunc("/signup",userHandler.SignUp)
  http.HandleFunc("/Login",userHandler.Login)
  http.HandleFunc("/profile",userHandler.Profile)
  http.HandleFunc("/health",health)
  http.HandleFunc("/",handler)
  log.Fatal(http.ListenAndServe(":8080",nil))
}
