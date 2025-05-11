package service

import (
  "github.com/HarshithRajesh/app-chat/internal/repository"
  "context"
  "github.com/redis/go-redis/v9"
  "fmt"
  "time"
)

func StartMessageConsumer(ctx context.Context,redisClient *redis.Client,streamName string,groupName string,
                          consumerName string,count int64,block time.Duration){
  for {
    res,err := repository.ReadMessagesFromGroup(ctx,redisClient,streamName,groupName,consumerName,count,block)
}
}

