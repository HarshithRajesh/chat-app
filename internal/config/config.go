package config

import (
  "database/sql"
  "fmt"
  "log"

  _ "github.com/lib/pq"
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
