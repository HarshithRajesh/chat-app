package realtime

import (
  "log"
  "sync"
)

type Hub struct{
  clients map[string]*Client
  broadcast chan []byte
  register chan *Client 
  unregister chan *Client
  mu sync.RWMutex
}

var _ IHub = (*Hub)(nil)

func Newhub() *Hub{
  return &Hub{
    broadcast: make(chan []byte,256),
    register: make(chan *Client,100),
    unregister: make(chan *Client,100),
    clients: make(map[string]*Client),
  }
}

func (h *Hub) Run(){
  log.Println("Websocket Hub event has started")
  for {
    select {
    case client := <- h.register:
      h.mu.Lock()
      h.clients[clients.UserID] = client
      h.mu.Unlock()
      log.Printf("Client registered: UserID = %s, Addr = %s. Toatl active users: %d",
                  client.UserID,client.Conn.RemoteAddr().String(),len(h.clients))
    
    case client := <- h.unregister:
      h.mu.Lock()
      if existingClient,ok := h.clients[client.UserID]; ok && existingClient == client{
        delete(h.clients,client.UserID)
        close(client.Send)
      }
      h.mu.Unlock()
      log.Printf("Client unregistered:UserId: %s, Addr = %s ,Total Active users = %d",
                  client.UserID,client.Conn.RemoteAddr().String(),len(h.clients))

    case client := <-h.broadcast:
      h.mu.RLock()
      for _,client := range h.clients{
        select {
        case client.Send <- message:
        default:
          close(client.Send)
          h.mu.RUnclock()
          h.mu.RLock()
          log.Printf("Client %s (UserID : %s) send buffer full or connection problematic,disconnecting",client.Conn.RemoteAddr().String(),client.UserID)
          
        }
      }
      h.mu.RUnclock()
    case <-ctx.Done():
      log.Println("Websocket Hub received shutdown")
      h.mu.Lock()
      for _,client := range h.clients{
        client.Conn.WriteMessage(websocket.CloseMessage,[]byte{})
        close(client.Send)
      }
      h.clients = make(map[string]*Client)
      h.mu.Unlock()
      return
    }
  }
}

func(h *Hub) RegisterClient(client *Client){

  select {
  case h.register <- client:
  
  default:
    log.Println("Hub register channel full, client registeration skipped")
      
    }
}

func (h *Hub) UnregisterClient (client *Client){

  select{
  case h.unregister <- client:
  
  default:
    log.Println("Hub unregister channel is full, client unregisteration delayed")
  }
}

func (h *Hub) BroadcastMessage(message[]byte){
  
  select {
  case h.brooadcast <- message:
  default:
    log.Println("Hub broadcast channel full, message dropped")
    
  }
}

func (h *Hub) SendToUser(userID string,message[]byte){
  
  h.mu.RLock()
  client,err := h.clients[UserID]
  h.mu.RUnclock()

  if err != nil{
    log.Println("Client with UserID %s not found in the Hub, message not sent",UserID)
    return
  }
  select {
  case client.Send <- message:
    log.Printf("Message enqued for User %s",userID)
  default:
    log.Printf("Client %s (UserID: %s) send buffer full, message to user dropped.",
                client.Conn.RemoteAddr().String(), userID)
  }
}
