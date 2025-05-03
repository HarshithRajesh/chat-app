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
  redisClient,err := config.ConnectRedisDB()
  if err != nil{
    fmt.Println("Redis client is not initialized")
  }
  defer redisClient.Close()
  userRepo := repository.NewUserRepository(db)
  userService := service.NewUserService(userRepo)
  userHandler := api.NewUserHandler(userService)


  chatRepo := repository.NewChatRepository(db)
  chatService := service.NewChatService(chatRepo,redisClient)
  chatHandler := api.NewChatHandler(chatService)

  http.HandleFunc("/signup",userHandler.SignUp)
  http.HandleFunc("/Login",userHandler.Login)
  http.HandleFunc("/profile",userHandler.Profile)
  http.HandleFunc("/contact",userHandler.Contact)
  http.HandleFunc("/contact/listcontacts",userHandler.ViewContact)
  http.HandleFunc("/user/message",chatHandler.SendMessage)
  http.HandleFunc("/user/message/history",chatHandler.GetMessage)
  http.HandleFunc("/health",health)
  http.HandleFunc("/",handler)
  log.Fatal(http.ListenAndServe(":8080",nil))
}
