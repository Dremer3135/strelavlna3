package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	PlayerInitLoaded


	BuyProb // string (diff)
	SellProb // string (probid)
	BoughtProb // BoughtProbMsg
	SolveProb // string (probid)
	SolvedProb // string (probid)
	SoldProb

	WriteMsg // WriteMsgMsg
	PlayerMsg
	AdminMsg
	GradeProb
	GradedProb

	NotEnoughMoney // required money (int)
	NotAvaiable // nil
	NotRunning // nil

	PlayerError // error
	TeamChanError
	ServerError // error

	TooManyPlayers
)

// type BoughtProbMsg struct {
// 	Id string `json:"id"`
// 	Diff string `json:"diff"`
// 	Name string `json:"name"`
// 	Text string `json:"text"`
// 	Imgs []string `json:"images"`
// 	Answer *string `json:"answer"`
// }

type WriteMsgMsg struct {
	probid string
	teamid string
	mtype string
	msg string
	admin bool
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

type ProbBoughtAdminMsg struct {
	team string
  prob string
}

type SolveProbMsg struct {
	probid string
	text string
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
				msg.callback <- Msg{PlayerError, id, self, "join"}
				break
			}
			if len(players) >= 5 {
				msg.callback <- Msg{PlayerError, id, self, "join"}
				break
			}
			plid := strconv.Itoa(len(players))
		  plchan := make(chan Msg, 10)
		  go playerManager(data, plchan, self, plid)
		  players[plid] = plchan

			iload := initLoad(conn, id, plid)
		  plchan <- Msg{PlayerInitLoaded, id, self, iload}
		
		case BuyProb:
			diff, ok := msg.data.(string)
		  if !ok { msg.callback <- Msg{PlayerError, id, self, "buy"}; break }

			prob, err := buyProb(conn, id, diff)
			if err != nil { msg.callback <- Msg{PlayerError, id, self, err } }

		  msg.callback <- Msg{BoughtProb, id, self, prob}

			admins <- Msg{BoughtProb, id, self, prob}
		

		case SellProb:
			probid, ok := msg.data.(string)
		  if !ok { msg.callback <- Msg{PlayerError, id, self, "sell"}; break }

			err := sellProb(conn, id, probid)
			if err != nil { msg.callback <- Msg{PlayerError, id, self, err } }

			// getTCorr

		  msg.callback <- Msg{SoldProb, id, self, probid}

		case WriteMsg:

			data, ok := msg.data.(WriteMsgMsg)
		  if !ok { msg.callback <- Msg{InvalidMessage, id, self, "write"}; break }

			pushTLine(conn, msg.from, data.probid, TLineAtom{MSidePlayer, MTypeText, data.msg, time.Now()})

			if data.mtype == MTypeGrade {
				for _, pl := range players {
					pl <- Msg{SolvedProb, id, self, data.probid}
				}
			}
			admins <- Msg{PlayerMsg, id, self, msg.data}
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

			case "solve":
				probid, ok := msg.Data["probid"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "sell"}
					break
				}
				solution, ok := msg.Data["solution"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "solve"}
					break
				}
				team <- Msg{SolvedProb, id, self, SolveProbMsg{probid, solution}}

			case "write":
				probid, ok := msg.Data["probid"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				msgd, ok := msg.Data["msg"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				mtype, ok := msg.Data["mtype"]
				if !ok {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				if len(msgd) > 50 {
					self <- Msg{InvalidMessage, id, self, "write"}
					break
				}
				team <- Msg{WriteMsg, id, self, WriteMsgMsg{probid: probid, msg: msgd, mtype: mtype, admin: false}}

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
			data, ok := msg.data.(Prob)
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

func adminManager(self chan Msg) {
	teams := map[string]chan Msg{}
	correctors := map[string]chan Msg{}
	conn := NewRdbConn()
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
			tickets := getCorrTickets(conn, msg.from)
		  for _, t := range tickets {
				ts := strings.Split(t, ":")
				prob := getProb(conn, ts[1])
				adminid := ""
				for _, a := range prob.Queue {
					_, ok := correctors[a]
					if ok {
						adminid = a
						break
					}
				} 
				if adminid == "" { break }
				corr := correctors[adminid]
				corr <- Msg{BoughtProb, "", self, prob}
				setTCorr(conn, msg.from, prob.Id, adminid)
				addCorrTicket(conn, adminid, msg.from, prob.Id)
			}
		
		

		case PlayerJoinRequest:
			data, ok := msg.data.(PlayerJoinReqMsg)
			if !ok { break }
			team, ok := teams[data.Teamid]
			if !ok { break }
			team <- Msg{PlayerJoinRequest, "", self, data.Conn}

		case BoughtProb:
			prob, ok := msg.data.(Prob)
			if !ok { break }
		  adminid := ""
			for _, a := range prob.Queue {
				_, ok := correctors[a]
				if ok {
					adminid = a
					break
				}
			} 
			if adminid == "" { break }
		  corr := correctors[adminid]
			corr <- Msg{BoughtProb, "", self, prob}
			setTCorr(conn, msg.from, prob.Id, adminid)
		  addCorrTicket(conn, adminid, msg.from, prob.Id)

		case SoldProb:
			probid, ok := msg.data.(string)
			if !ok { break }
			corrid := getTCorr(conn, msg.from, probid)
		  correctors[corrid] <- Msg{SoldProb, "", self, probid}

		case SolvedProb:
			probid, ok := msg.data.(string)
			if !ok { break }
			corrid := getTCorr(conn, msg.from, probid)
		  correctors[corrid] <- Msg{SolvedProb, "", self, probid}

		case PlayerMsg:
			data, ok := msg.data.(WriteMsgMsg)
			if !ok { break }
			corrid := getTCorr(conn, msg.from, data.probid)
		  correctors[corrid] <- Msg{PlayerMsg, "", self, data}

		case GradeProb:
			data, ok := msg.data.(ProbBoughtAdminMsg)
			if !ok { break }
			teams[data.team] <- Msg{GradeProb, "", self, data.prob}
			msg.callback <- Msg{GradedProb, "", self, data}

		
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
