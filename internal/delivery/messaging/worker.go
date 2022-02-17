package messaging

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"sync"
)

type Consumer struct {
	waitGrp       sync.WaitGroup
	cacheClient   *redis.Client
	streamName    string
	consumerGroup string
	consumerName  string
}

func NewConsumer(cacheClient *redis.Client, streamName, consumerGroup, consumerName string) *Consumer {
	return &Consumer{
		waitGrp:       sync.WaitGroup{},
		cacheClient:   cacheClient,
		streamName:    streamName,
		consumerGroup: consumerGroup,
		consumerName:  consumerName,
	}
}

func (c *Consumer) StartConsumer(ctx context.Context, doneChan chan os.Signal) {

	log.Printf("\nStarting Consumer - Group: %s - Name: %s", c.consumerGroup, c.consumerName)

	for {
		select {
		case <-doneChan:
			fmt.Println("Stopping consumer (received signal)...")
			c.waitGrp.Wait()
			return

		default:

			streams, err := c.cacheClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Streams:  []string{c.streamName, ">"},
				Group:    c.consumerGroup,
				Consumer: c.consumerName,
				Count:    1,
				Block:    0,
			}).Result()

			if err != nil {
				log.Printf("err on consume events: %+v\n", err)
				return
			}

			c.waitGrp.Add(len(streams[0].Messages))

			for _, stream := range streams[0].Messages {
				go c.processMessage(ctx, stream)
			}
		}
	}
}
