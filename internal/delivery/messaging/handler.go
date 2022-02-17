package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/felipeagger/go-distributed-websocket/internal/entity"
	"github.com/felipeagger/go-distributed-websocket/pkg/cache"
	"github.com/go-redis/redis/v8"
)

func (c *Consumer) processMessage(ctx context.Context, message redis.XMessage) {
	defer c.waitGrp.Done()

	var payload entity.Message
	if err := json.Unmarshal([]byte(message.Values["data"].(string)), &payload); err != nil {
		fmt.Printf("\nunmarhsal.Error: %v\n", err)
		return
	}

	//Process your message here
	response := entity.Message{
		UserID:      payload.UserID,
		Origin:      payload.Origin,
		Data:        fmt.Sprintf("ProcessedData: %s", payload.Data),
		ReceivedBy:  payload.ReceivedBy,
		ProcessedBy: c.consumerName,
	}

	log.Printf("Host: %v - ProcessedMessage: %v", c.consumerName, response)

	err := cache.Publish(ctx, c.cacheClient, response)
	if err != nil {
		fmt.Printf("\nprocessMessage.Error: %v\n", err)
		return
	}

	c.cacheClient.XAck(ctx, c.streamName, c.consumerGroup, message.ID)
}
