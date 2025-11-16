package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type MessageType int
const (
	WsError MessageType = iota // error
	InvalidMessage // string (cmd)

	PlayerJoinRequest // PlayerJoinReqMsg
	TeamRegister // id (string)

	CorrectorJoined
	CorrectorLeft // error

	PlayerLeft // error
	PlayerJoined


	BuyProb // string (diff)
	SellProb // string (probid)
	BoughtProb // BoughtProbMsg

	WriteMsg // WriteMsgMsg

	NotEnoughMoney // required money (int)
	NotAvaiable // nil
	NotRunning // nil

	TeamChanError
	ServerError // error

	TooManyPlayers
)

type BoughtProbMsg struct {
	Id string `json:"id"`
	Diff string `json:"diff"`
	Name string `json:"name"`
	Text string `json:"text"`
	Imgs []string `json:"images"`
	Answer *string `json:"answer"`
}

type WriteMsgMsg struct {
	probid string
	msg string
}

type PlayerJoinedMsg struct {
	id string
	addr chan Msg
}

type InWsMsg struct {
	Name string `json:"name"`
	Data map[string]string `json:"data"`
}

type OutWsMsg struct {
	Name string `json:"name"`
	Data any `json:"data"`
}

type Msg struct {
	tp MessageType
	from string
	callback chan Msg
	data any
}

func isRunning(conn *redis.Client) bool {
	res, err := conn.Get(ctx, "run").Result()
	if err != nil { return false }
	return res == "on"
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

func popProb(conn *redis.Client, team string, diff string) (string, error) {
	return conn.SPop(ctx, "freeprobs:" + team + ":" + diff).Result()
}

func pushBoughtProb(conn *redis.Client, team string, id string) (int64, error) {
	return conn.SAdd(ctx, "boughtprobs:" + team, id).Result()
}

func getProbInfo(conn *redis.Client, id string) (BoughtProbMsg, error) {
	res := BoughtProbMsg{}
	res.Id = id
	diff, err := conn.Get(ctx, "prob:" + id + ":diff").Result()
	if err != nil { return res, err }
	res.Diff = diff
	name, err := conn.Get(ctx, "prob:" + id + ":name").Result()
	if err != nil { return res, err }
	res.Name = name
	text, err := conn.Get(ctx, "prob:" + id + ":text").Result()
	if err != nil { return res, err }
	res.Text = text
	imgs, err := conn.SMembers(ctx, "prob:" + id + ":images").Result()
	if err != nil { return res, err }
	res.Imgs = imgs
	ans, err := conn.Get(ctx, "prob:" + id + ":ans").Result()
	if err != nil { return res, err }
	res.Answer = &ans

	return res, nil
}

func setMoney(conn *redis.Client, id string, money int) error {
  return conn.Set(ctx, "money:" + id, money, time.Duration(0)).Err()
}

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

		case PlayerJoinRequest:
			data, ok := msg.data.(*websocket.Conn)
		  if !ok {
				msg.callback <- Msg{InvalidMessage, id, self, "join"}
				break
			}
			if len(players) >= 5 {
				msg.callback <- Msg{TooManyPlayers, id, self, "join"}
				break
			}
			plid := strconv.Itoa(len(players))
		  plchan := make(chan Msg, 10)
		  go playerManager(data, plchan, self, plid)
		  players[plid] = plchan
			// TODO: initload

		case BuyProb:
			diff, ok := msg.data.(string)
		  if !ok { msg.callback <- Msg{InvalidMessage, id, self, "buy"}; break }
			if !isRunning(conn) {
				msg.callback <- Msg{NotRunning, id, self, "buy"};
				break
			}
			money, err := getMoney(conn, id)
		  if err != nil { msg.callback <- Msg{ServerError, id, self, err}; break }

			price, err := getPrice(conn, BuyCost, diff)
		  if err != nil { msg.callback <- Msg{ServerError, id, self, err}; break }

		  if price > money { msg.callback <- Msg{NotEnoughMoney, id, self, price}; break }

			probid, err := popProb(conn, id, diff)
		  if err != nil { msg.callback <- Msg{NotAvaiable, id, self, err}; break }

			_, err = pushBoughtProb(conn, id, probid)
		  if err != nil { msg.callback <- Msg{ServerError, id, self, nil}; break }

			err = setMoney(conn, id, money - price)
		  if err != nil { msg.callback <- Msg{ServerError, id, self, nil}; break }
			
		  info, err := getProbInfo(conn, probid)
		  if err != nil { msg.callback <- Msg{ServerError, id, self, nil}; break }

		  msg.callback <- Msg{BoughtProb, id, self, info}

		case WriteMsg:
		  admins <- msg
			

		}
	}
}

func playerManager(ws *websocket.Conn, self chan Msg, team chan Msg, id string) {
	go func() {
		for {
			var msg InWsMsg
			err := ws.ReadJSON(&msg)
			if err != nil {
				self <- Msg{WsError, id, self, err}
				break
			}
			switch msg.Name {

			case "buy":
				diff, ok := msg.Data["diff"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "buy"}
					break
				}
				team <- Msg{BuyProb, id, self, diff}

			case "sell":
				probid, ok := msg.Data["probid"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "sell"}
					break
				}
				team <- Msg{SellProb, id, self, probid}

			case "write":
				probid, ok := msg.Data["probid"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				msg, ok := msg.Data["msg"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				if len(msg) > 50 {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				team <- Msg{WriteMsg, id, self, WriteMsgMsg{probid: probid, msg: msg}}

			}
		}
	}()
	chanloop: for {
		msg, ok := <- self
		fmt.Printf("client: %#v\n", msg)
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

		case BoughtProb:
			data, ok := msg.data.(BoughtProbMsg)
			if !ok { break }
			err := ws.WriteJSON(OutWsMsg{Name: "bought", Data: data})
			if err != nil {
				self <- Msg{WsError, id, self, err}
			}
		}
	}
}

type PlayerJoinReqMsg struct {
	Teamid string
	Conn *websocket.Conn
}

type CorrectorJoinedReqMsg struct {
	Id string
	Conn *websocket.Conn
}

func getProbCorrector(conn *redis.Client, teamid string, probid string) (string, error) {
	return conn.Get(ctx, "probcorr:" + probid).Result()
}

func adminManager(self chan Msg) {
	teams := map[string]chan Msg{}
	correctors := map[string]chan Msg{}
	for {
		msg := <- self
		switch msg.tp {

		case TeamRegister:
			id, ok := msg.data.(string)
		  if !ok { break }
			tchan := make(chan Msg, 10)
		  go teamManager(tchan, self, id)
		  teams[id] = tchan

		case TeamChanError:
			tchan := make(chan Msg, 10)
		  go teamManager(tchan, self, msg.from)
		  teams[msg.from] = tchan

		case CorrectorJoined:
			data, ok := msg.data.(CorrectorJoinedReqMsg)
		  if !ok { break }
			cchan := make(chan Msg, 10)
		  go correctorManager(data.Conn, cchan, self, data.Id)
		  correctors[data.Id] = cchan

		case CorrectorLeft:
			_, ok := msg.data.(error)
		  if !ok { break }
		  delete(correctors, msg.from)

		case PlayerJoinRequest:
			data, ok := msg.data.(PlayerJoinReqMsg)
			if !ok { break }
			team, ok := teams[data.Teamid]
			if !ok { break }
			team <- Msg{PlayerJoinRequest, "", self, data.Conn}

		case WriteMsg:
			// data, ok := msg.data.(WriteMsgMsg)
			// if !ok { break }
		  // corrector, err := 

		  
		  

		}
	}
}

func correctorManager(ws *websocket.Conn, self chan Msg, admins chan Msg, id string) {
	go func() {
		for {
			var msg InWsMsg
			err := ws.ReadJSON(&msg)
			if err != nil {
				self <- Msg{WsError, id, self, err}
				break
			}
			switch msg.Name {
			}
		}
	}()
	chanloop: for {
		msg, ok := <- self
		if !ok {
			ws.Close()
			admins <- Msg{WsError, id, self, errors.New("chan closed")}
			break
		}
		switch msg.tp {

		case WsError:
			data, ok := msg.data.(error)
			if !ok { break }
			ws.Close()
			admins <- Msg{CorrectorLeft, id, self, data}
			break chanloop
		}
	}
}
