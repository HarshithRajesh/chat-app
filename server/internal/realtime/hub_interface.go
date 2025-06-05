package realtime

type IHub interface{
  Run()
  RegisterClient(client *Client)
  UnregisterClient(client *Client)
  BroadcastMessage(message []byte)
  SendToUser(userID string,message []byte)
}
