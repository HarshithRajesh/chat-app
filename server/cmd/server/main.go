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
  "context"
  "os"
  "time"
  "strings"
  "os/signal"
  "syscall"
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
  ctx := context.Background()
  redisClient,err := config.ConnectRedisDB()
  if err != nil{
    fmt.Println("Redis client is not initialized")
  }
  defer redisClient.Close()

 server := &http.Server{
    Addr :  ":8080",
    Handler : nil,
    ReadTimeout:  10* time.Second,
    WriteTimeout: 10* time.Second,
    MaxHeaderBytes: 1<<20,
  }

  userRepo := repository.NewUserRepository(db)
  userService := service.NewUserService(userRepo)
  userHandler := api.NewUserHandler(userService)


  chatRepo := repository.NewChatRepository(db)
  chatService := service.NewChatService(chatRepo,redisClient)
  chatHandler := api.NewChatHandler(chatService)
//add
go func(){
    log.Println("Server running in the port :8080")
    err := server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed{
      log.Printf("Error while running the server: %v",err)
    }
    log.Println("Http Server closed")
  }()
  http.HandleFunc("/signup",userHandler.SignUp)
  http.HandleFunc("/Login",userHandler.Login)
  http.HandleFunc("/profile",userHandler.Profile)
  http.HandleFunc("/contact",userHandler.Contact)
  http.HandleFunc("/contact/listcontacts",userHandler.ViewContact)
  http.HandleFunc("/user/message",chatHandler.SendMessage)
  http.HandleFunc("/user/message/history",chatHandler.GetMessage)
  http.HandleFunc("/health",health)
  http.HandleFunc("/",handler)


  _,err = redisClient.XGroupCreate(ctx,"chat_stream","chat_processor","$").Result()
  if err != nil {
    if strings.Contains(err.Error(), "BUSYGROUP Consumer Group name already exists") {
        fmt.Println("Redis group already exists")
    } else {
        log.Fatalf("Error creating Redis Consumer Group: %v", err)
    }
  } else {
    fmt.Println("Redis Consumer group created on stream 'chat_stream'")
  }
  appCtx ,cancelFunc:= context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1) // Create the channel
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // Tell Go to send SIGINT/SIGTERM here

	go func() {
		sig := <-sigChan

		log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

		cancelFunc()

	}()
 
	log.Println("Signal handling configured. Press Ctrl+C to shut down gracefully.")
  // defer cancel()
  hostname,_ := os.Hostname()
  consumerName := fmt.Sprintf("consumer-%s-%d",hostname,os.Getpid())
  readCount := int64(10)
  blockDuration := time.Duration(0)

  go service.StartMessageConsumer(appCtx,redisClient,"chat_stream","chat_processor",consumerName,readCount,blockDuration)
 
  log.Println("Application is running. Waiting for shutdown signal (Press Ctrl+C to stop)...")
	<-appCtx.Done()

  log.Println("Shutdown signal received. Main goroutine unblocked. Application stopping.")
  // messages,err := repository.ReadMessageFromStream(ctx,redisClient,"chat_stream","0",3)
  //   if err != nil{
  //   fmt.Printf("Error reading from the stream: %v\n",err)
  // }else{
  //   for _,stream := range messages{
  //     for _,msg := range stream.Messages{
  //       fmt.Printf("Message Id : %s\n",msg.ID)
  //       fmt.Println("Values:")
  //       for field,value := range msg.Values{
  //         fmt.Printf("  %s: %v\n",field,value)
  //       }
  //     }
  // }
  //     }
  
 

 

}
