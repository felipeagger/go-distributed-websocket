package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipeagger/go-distributed-websocket/internal/entity"
	"github.com/felipeagger/go-distributed-websocket/pkg/utils"
	"github.com/go-redis/redis/v8"
	"strings"
)

func NewRedisClient(ctx context.Context, hostname string) (*redis.Client, error) {
	cacheClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", hostname),
	})

	if err := cacheClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return cacheClient, nil
}

func Publish(ctx context.Context, cacheClient *redis.Client, msg entity.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return cacheClient.Publish(ctx, utils.GetTopicName(msg.UserID, msg.Origin), payload).Err()
}

func CreateConsumerGroup(ctx context.Context, cacheClient *redis.Client, streamTopicName, consumerGroup string) error {

	if _, err := cacheClient.XGroupCreateMkStream(ctx, streamTopicName, consumerGroup, "0").Result(); err != nil {

		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)

			return err
		}

	}

	return nil
}
