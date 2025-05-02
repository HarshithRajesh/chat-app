package config

import (
  "database/sql"
  "fmt"
  "log"
  "context"
  _ "github.com/lib/pq"
  "github.com/redis/go-redis/v9"
)


func ConnectDB() *sql.DB {
  connStr := "postgres://neo:babe@localhost:5432/chat_app?sslmode=disable" 
  db,err := sql.Open("postgres",connStr)
  if err != nil {
    log.Fatal("Unable to connect to database:",err)
  }

  if err = db.Ping();err != nil{
    log.Fatal("Database connection failed:",err)
  }
  fmt.Println("Database connected")
  return db
}

var RedisClient *redis.Client

func ConnectRedisDB(){
  RedisClient = redis.Client(&redis.Options{
    Addr :"localhost:6379",
    Password : 0;
    DB : 0,
    Protocol : 2,
  })

  ctx := context.Background()
  _,err := RedisClient.Ping(ctx).Result()
  if err != nil{
    log.Fatalf("Failed to connect to the Redis Database") 
  }else {
    log.Println("Connected to Redis")
  }
}
