package main

import (
	"errors"
	"fmt"

	"golang.org/x/net/websocket"
)

type MessageType int
const (
	WsError MessageType = iota // error
	PlayerLeft // error
	InvalidMessage // string (cmd)
	BuyProb // string (diff)
	SellProb // string (probid)
	BoughtProb // BoughtProbMsg
)

type BoughtProbMsg struct {
	Id string
	Diff string
	Name string
	Text string
	Img string
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
