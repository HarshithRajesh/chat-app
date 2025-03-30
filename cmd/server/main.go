package main 

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
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
  http.HandleFunc("/health",health)
  http.HandleFunc("/",handler)
  log.Fatal(http.ListenAndServe(":8080",nil))
}
