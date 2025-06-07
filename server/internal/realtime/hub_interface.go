package realtime
import "context"

type IHub interface{
  Run(ctx context.Context)
  RegisterClient(client *Client)
  UnregisterClient(client *Client)
  BroadcastMessage(message []byte)
  SendToUser(userID string,message []byte)
}
