package main

import (
	"context"
	"github.com/felipeagger/go-distributed-websocket/internal/delivery/messaging"
	"github.com/felipeagger/go-distributed-websocket/pkg/cache"
	"github.com/felipeagger/go-distributed-websocket/pkg/utils"
	"os"
	"os/signal"
	"syscall"
)

var (
	streamName    string = os.Getenv("TOPIC")
	consumerGroup string = os.Getenv("GROUP")
)

func main() {
	ctx := context.TODO()
	hostname, _ := os.Hostname()

	cacheClient, err := cache.NewRedisClient(ctx, os.Getenv("CACHE_HOST"))
	if err != nil {
		panic(err)
	}

	defer cacheClient.Close()

	streamName, consumerGroup = utils.SetDefaultEnvs(streamName, consumerGroup)

	if err := cache.CreateConsumerGroup(ctx, cacheClient, streamName, consumerGroup); err != nil {
		panic(err)
	}

	chanOS := make(chan os.Signal)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)

	consumer := messaging.NewConsumer(cacheClient, streamName, consumerGroup, hostname)
	consumer.StartConsumer(ctx, chanOS)
}