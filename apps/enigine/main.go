package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"net/http"
	// "time"

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

	http.HandleFunc("/corr/ws", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		corrid, err := rdb.Get(ctx, "corrtoken:" + token).Result()
		if err != nil {
			http.Error(w, "invalid token", 400)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		fmt.Println(corrid)
		achan <- Msg{CorrectorJoined, "", mchan, CorrectorJoinedReqMsg{corrid, conn}}
	})

	// http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
	// 	token := r.URL.Query().Get("token")
	// 	_, err := rdb.Get(ctx, "playtoken:" + token).Result()
	// 	if err != nil {
	// 		http.Error(w, "invalid token", 400)
	// 		return
	// 	}
	// 	w.WriteHeader(200)
	// })
	//
	// http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
	// 	state, err := getState(rdb)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	stime, err := getStartTime(rdb)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	etime, err := getEndTime(rdb)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(200)
	// 	json.NewEncoder(w).Encode(struct{
	// 		State string `json:"state"`
	// 		StartTime int `json:"starttime"`
	// 		EndTime int `json:"endtime"`
	// 	}{
	// 		State: state,
	// 		StartTime: int(stime.Sub(time.Now()).Milliseconds()),
	// 		EndTime: int(etime.Sub(time.Now()).Milliseconds()),
	// 	})
	// })

	if err := http.ListenAndServe(":7777", nil); err != nil {
		panic([]any{"ListenAndServe:", err})
	}
}
