package main

import (
	"errors"
	"fmt"

	"golang.org/x/net/websocket"
	"github.com/redis/go-redis/v9"
)

type MessageType int
const (
	WsError MessageType = iota // error
	InvalidMessage // string (cmd)

	PlayerLeft // error
	PlayerJoined


	BuyProb // string (diff)
	SellProb // string (probid)
	BoughtProb // BoughtProbMsg

	NotEnoughMoney // required money (int)

	TeamChanError
	ServerEroor // error
)

type BoughtProbMsg struct {
	Id string
	Diff string
	Name string
	Text string
	Imgs []string
	Answer *string
}

type PlayerJoinedMsg struct {
	id string
	addr chan Msg
}

type WsMsg struct {
	name string
	data map[string]string
}

type Msg struct {
	tp MessageType
	from string
	callback chan Msg
	data any
}

func getMoney(conn *redis.Client, id string) (int, error) {
	return conn.Get(ctx, "money:" + id).Int()
}

const (
	BuyCost = "buy"
	SellCost = "sell"
	SolveCost = "solve"
)

func getPrice(conn *redis.Client, costType string, diff string) (int, error) {
	return conn.Get(ctx, "cost:" + costType + ":" + diff).Int()
}

func popProb(conn *redis.Conn, )

func teamManager(self chan Msg, admins chan Msg, id string) {
	players := map[string]chan Msg{}
	conn := NewRdbConn()
	for {
		msg, ok := <- self
		if !ok {
			admins <- Msg{TeamChanError, id, self, ok}
			break
		}
		switch msg.tp {

		case PlayerLeft:
			delete(players, msg.from)

		case PlayerJoined:
			data, ok := msg.data.(PlayerJoinedMsg)
		  if !ok {
				msg.callback <- Msg{InvalidMessage, id, self, "join"}
				break
			}
			players[data.id] = data.addr
			// TODO: initload

		case BuyProb:
			diff, ok := msg.data.(string)
		  if !ok {
				msg.callback <- Msg{InvalidMessage, id, self, "join"}
				break
			}
			money, err := getMoney(conn, id)
		  if err != nil {
				msg.callback <- Msg{ServerEroor, id, self, err}
				break
			}
			price, err := getPrice(conn, BuyCost, diff)
		  if err != nil {
				msg.callback <- Msg{ServerEroor, id, self, err}
				break
			}
		  if price > money {
				msg.callback <- Msg{NotEnoughMoney, id, self, price}
				break
			}
		  




		}
	}
}

func playerManager(ws *websocket.Conn, self chan Msg, team chan Msg, id string) {
	go func() {
		for {
			var msg WsMsg
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				self <- Msg{WsError, id, self, err}
				break
			}
			switch msg.name {

			case "buy":
				diff, ok := msg.data["diff"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "buy"}
					break
				}
				self <- Msg{BuyProb, id, self, diff}

			case "sell":
				probid, ok := msg.data["probid"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "sell"}
					break
				}
				self <- Msg{SellProb, id, self, probid}


			}
		}
	}()
	chanloop: for {
		msg, ok := <- self
		if !ok {
			ws.Close()
			team <- Msg{PlayerLeft, id, self, errors.New("chan closed")}
			break
		}
		switch msg.tp {

		case WsError:
			data, ok := msg.data.(error)
			if !ok { break }
			ws.Close()
			team <- Msg{PlayerLeft, id, self, data}
			break chanloop

		case BuyProb:
			data, ok := msg.data.(string)
			if !ok { break }
			team <- Msg{BuyProb, id, self, data}

		case BoughtProb:
			data, ok := msg.data.(BoughtProbMsg)
			if !ok { break }
			err := websocket.JSON.Send(ws, data)
			if err != nil {
				self <- Msg{WsError, id, self, err}
			}
		


		}
	}
}
