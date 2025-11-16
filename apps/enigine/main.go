package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var upgrader = websocket.Upgrader{}

func NewRdbConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func main() {
	rdb := NewRdbConn()

	fmt.Println(rdb.Get(ctx, "a").String())

	achan := make(chan Msg, 10)
	go adminManager(achan)

	mchan := make(chan Msg, 10)

	teams, err := rdb.SMembers(ctx, "teams").Result()
	if err != nil { panic(err) }

	for _, team := range teams {
		achan <- Msg{TeamRegister, "", nil, team}
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		teamid, err := rdb.Get(ctx, "playtoken:" + token).Result()
		if err != nil {
			http.Error(w, "invalid token", 400)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		fmt.Println(teamid)
		achan <- Msg{PlayerJoinRequest, "", mchan, PlayerJoinReqMsg{teamid, conn}}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic([]any{"ListenAndServe:", err})
	}
}
