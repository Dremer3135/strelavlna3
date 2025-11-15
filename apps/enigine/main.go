package main

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

var ctx = context.Background()

func wsHandler(ws *websocket.Conn) {
	// websocket.Message.Receive()
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	teams, err := rdb.SMembers(ctx, "teams").Result()
	if err != nil { panic(err) }



	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic([]any{"ListenAndServe:", err})
	}
}
