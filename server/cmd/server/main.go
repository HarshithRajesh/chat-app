package main 

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "github.com/HarshithRajesh/app-chat/internal/api"
  "github.com/HarshithRajesh/app-chat/internal/config"
  "github.com/HarshithRajesh/app-chat/internal/repository"
  "github.com/HarshithRajesh/app-chat/internal/realtime"
  "github.com/HarshithRajesh/app-chat/internal/service"
  "context"
  "os"
  "time"
  "strings"
  "os/signal"
  "syscall"
  // "golang.org/x/net/websocket"
  // "github.com/gorilla/websocket"
)
type Response struct{
  Message string `json:"message"`
}
func health(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  response := Response{Message: "Hi Welcome to Chaat"}
  json.NewEncoder(w).Encode(response)
}
// var upgrader = websocket.Upgrader{
//     ReadBufferSize : 1024,
//     WriteBufferSize : 1024,
//     CheckOrigin: func(r *http.Request) bool{return true},
//   }
// func reader(conn *websocket.Conn){
//
//   defer func(){
//     log.Println("Client is disconnected")
//     conn.Close()
//   }()
//
//   for {
//     messageType,p,err := conn.ReadMessage()
//
//     if err !=  nil{
//       log.Printf("Webosocket upgrade error:",err)
//       return
//     }
//     fmt.Println(string(p))
//     if err := conn.WriteMessage(messageType,p);err != nil{
//       log.Printf("Error writing the message: ",err)
//       return
//     }
//   }
// }


func handler(w http.ResponseWriter, r *http.Request){
  // upgrader.CheckOrigin = func(r *http.Request) bool{return true}

  // ws,err := upgrader.Upgrade(w,r,nil)
  // if err != nil{
  //   log.Println(err)
  // }
  // defer ws.Close()
  log.Println("Hi,there, Welocome to my chaat ")
  // reader(ws)
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

  appCtx ,cancelFunc:= context.WithCancel(context.Background())

  hub := realtime.NewHub()
  go hub.Run(appCtx)
  log.Println("websocket hub started and running in Background")

  userRepo := repository.NewUserRepository(db)
  userService := service.NewUserService(userRepo)
  userHandler := api.NewUserHandler(userService)


  chatRepo := repository.NewChatRepository(db)
  chatService := service.NewChatService(chatRepo,redisClient)
  chatHandler := api.NewChatHandler(chatService)

  wsChatHandler := api.NewWsChatHandler(hub,userService,chatService)
  log.Println("websocket chat handler initialized with hub,user service and chat service")

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

  http.HandleFunc("/ws/chat",wsChatHandler.HandleWebSocket)
  log.Println("Http routes registered, including /ws/chat for websocket")
  // http.HandleFunc("/chat",realtime.websocket.Handler(server.handleWs))

  _,err = redisClient.XGroupCreateMkStream(ctx,"chat_stream","chat_processor","$").Result()
  if err != nil {
    if strings.Contains(err.Error(), "BUSYGROUP Consumer Group name already exists") {
        fmt.Println("Redis group already exists")
    } else {
        log.Fatalf("Error creating Redis Consumer Group: %v", err)
    }
  } else {
    fmt.Println("Redis Consumer group created on stream 'chat_stream'")
  }

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
