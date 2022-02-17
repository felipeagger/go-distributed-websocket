package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipeagger/go-distributed-websocket/internal/entity"
	"github.com/felipeagger/go-distributed-websocket/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrader holds the websocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	Hostname string
	StreamName string
	RedisClient *redis.Client
}

// SocketHandler echos websocket messages back to the client.
func (h *Handler) SocketHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	userId := r.URL.Query().Get("userId")

	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		return
	}

	go h.subscribe(r.Context(), conn, origin, userId)

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		var msg entity.Message
		if data != nil && string(data) != "" {
			err = json.Unmarshal(data, &msg)
			if err != nil {
				errMsg := fmt.Sprintf("\nunmarshal.ReadMessage (Invalid Payload): %v", err)
				sendResponse(conn, messageType, errMsg)
				continue
			}
		}

		msg.UserID = userId
		msg.ReceivedBy = h.Hostname

		err = h.publish(r.Context(), msg)
		if err != nil {
			log.Printf("conn.PublishMessage: %v", err)
			continue
		}

		log.Printf("Host: %v - PublishMessage: %v", h.Hostname, msg)
	}
}

func sendResponse(conn *websocket.Conn, msgType int, message string) {
	if err := conn.WriteMessage(msgType, []byte(message)); err != nil {
		fmt.Printf("\nconn.WriteMessage: %v", err)
		return
	}
}

func (h *Handler) publish(ctx context.Context, msg entity.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = h.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: h.StreamName,
		Values: map[string]interface{}{"data": payload},
	}).Result()

	return err
}

func (h *Handler) subscribe(ctx context.Context, conn *websocket.Conn, origin, userID string) {
	subscriber := h.RedisClient.Subscribe(ctx, utils.GetTopicName(userID, origin))

	for {
		select {
		case <-ctx.Done():
			fmt.Println("subscribe loop cancelled")
		default:
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			sendResponse(conn, 1, msg.Payload)
		}
	}
}