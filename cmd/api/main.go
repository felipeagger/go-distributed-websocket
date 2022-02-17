package main

import (
	"context"
	ws "github.com/felipeagger/go-distributed-websocket/internal/delivery/websocket"
	"github.com/felipeagger/go-distributed-websocket/pkg/cache"
	"github.com/felipeagger/go-distributed-websocket/pkg/utils"
	"log"
	"net/http"
	"os"
)

var (
	streamName    string = os.Getenv("TOPIC")
)

func main() {
	streamName, _ = utils.SetDefaultEnvs(streamName, "")
	hostname, _ := os.Hostname()

	cacheClient, err := cache.NewRedisClient(context.TODO(), os.Getenv("CACHE_HOST"))
	if err != nil {
		panic(err)
	}

	defer cacheClient.Close()

	handler := ws.Handler{
		Hostname: hostname,
		RedisClient: cacheClient,
		StreamName: streamName,
	}

	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/ws", handler.SocketHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	log.Printf("\nListening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}