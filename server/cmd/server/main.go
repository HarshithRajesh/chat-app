package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/HarshithRajesh/app-chat/internal/api"
	"github.com/HarshithRajesh/app-chat/internal/config"
	"github.com/HarshithRajesh/app-chat/internal/realtime"
	"github.com/HarshithRajesh/app-chat/internal/repository"
	"github.com/HarshithRajesh/app-chat/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func health(c *gin.Context) {
	response := Response{Message: "Hi Welcome to Chaat"}
	c.JSON(http.StatusOK, response)
}

func handler(c *gin.Context) {
	log.Println("Hi,there, Welcome to my chat")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to chat app"})
}

func main() {
	// Set Gin mode (can be gin.ReleaseMode for production)
	gin.SetMode(gin.DebugMode)

	db := config.ConnectDB()
	ctx := context.Background()
	redisClient, err := config.ConnectRedisDB()
	if err != nil {
		fmt.Println("Redis client is not initialized")
	}
	defer redisClient.Close()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	appCtx, cancelFunc := context.WithCancel(context.Background())

	hub := realtime.NewHub()
	go hub.Run(appCtx)
	log.Println("websocket hub started and running in Background")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	chatRepo := repository.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo, redisClient)
	chatHandler := api.NewChatHandler(chatService)

	wsChatHandler := api.NewWsChatHandler(hub, userService, chatService)
	log.Println("websocket chat handler initialized with hub,user service and chat service")

	// Define routes
	router.GET("/health", health)
	router.GET("/", handler)

	// User routes
	router.POST("/signup", userHandler.SignUp)
	router.POST("/Login", userHandler.Login)
	router.PUT("/profile", userHandler.Profile)

	// Contact routes
	router.POST("/contact", userHandler.Contact)
	router.GET("/contact/listcontacts", userHandler.ViewContact)

	// Message routes
	router.POST("/user/message", chatHandler.SendMessage)
	router.GET("/user/message/history", chatHandler.GetMessage)

	// WebSocket route (needs special handling)
	router.GET("/ws/chat", wsChatHandler.HandleWebSocket)

	log.Println("Http routes registered, including /ws/chat for websocket")

	// Create HTTP server with Gin router
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in goroutine
	go func() {
		log.Println("Server running on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while running the server: %v", err)
		}
		log.Println("Http Server closed")
	}()

	_, err = redisClient.XGroupCreateMkStream(ctx, "chat_stream", "chat_processor", "$").Result()
	if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP Consumer Group name already exists") {
			fmt.Println("Redis group already exists")
		} else {
			log.Fatalf("Error creating Redis Consumer Group: %v", err)
		}
	} else {
		fmt.Println("Redis Consumer group created on stream 'chat_stream'")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v. Initiating graceful shutdown...", sig)
		cancelFunc()
	}()

	log.Println("Signal handling configured. Press Ctrl+C to shut down gracefully.")

	hostname, _ := os.Hostname()
	consumerName := fmt.Sprintf("consumer-%s-%d", hostname, os.Getpid())
	readCount := int64(10)
	blockDuration := time.Duration(0)

	go service.StartMessageConsumer(appCtx, redisClient, "chat_stream", "chat_processor", consumerName, readCount, blockDuration)

	log.Println("Application is running. Waiting for shutdown signal (Press Ctrl+C to stop)...")
	<-appCtx.Done()

	log.Println("Shutdown signal received. Main goroutine unblocked. Application stopping.")
}
