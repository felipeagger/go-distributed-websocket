# go-distributed-websocket
Distributed Web Socket with Golang and Redis

## Dependencies

- gorilla/websocket
- go-redis

## Architecture

![Flow](/assets/distributed-websocket.png)

## Running

Compile with:

```
make build && make up
```

## On Browser Access

http://0.0.0.0:8008/

Send:

```
{"userId": "9e67a109-5c55-4e0b-8d5f-31b06ed4bb38", "origin": "web", "data": "Hello World"}
```

## Web Client

![Flow](/assets/web-client.png)