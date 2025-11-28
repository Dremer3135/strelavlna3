package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	// "strings"
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

	CorrectorInitLoaded

	Focus


	BuyProb // string (diff)
	SellProb // string (probid)
	BoughtProb // BoughtProbMsg
	SolveProb // string (probid)
	SolvedProb // string (probid)
	SoldProb
	SolveProbReq

	CorrWriteMsg
	CorrGrade
	CorrGraded
	WriteMsg // WriteMsgMsg
	RecMsg
	// GradeProb
	// GradedProb

	NotEnoughMoney // required money (int)
	NotAvaiable // nil
	NotRunning // nil

	UserError // error
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

type TeamProbMsg struct {
	team string
  prob string
}

type SolveProbMsg struct {
	probid string
	text string
}

type PlayerFocusMsg struct {
	probid string
	playerid string
}

type BoughtProbMsg struct {
	prob Prob
	money int
	remprobs map[string]int
}

type IdMoneyMsg struct {
	id string
	money int
}

func teamManager(self chan Msg, admins chan Msg, id string, tname string) {
	players := map[string]chan Msg{}
	conn := NewRdbConn()
	correctors := map[string]chan Msg{}
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
				msg.callback <- Msg{UserError, id, self, "join"}
				break
			}
			if len(players) >= 5 {
				msg.callback <- Msg{UserError, id, self, "join"}
				break
			}
			plidi := 0
			plid := "0"
			for {
				_, ok := players[plid]
				if !ok { break }
				plidi++
				plid = strconv.Itoa(plidi)
			}
		  plchan := make(chan Msg, 10)
		  go playerManager(data, plchan, self, plid, id)
		  players[plid] = plchan

			iload := initLoad(conn, id, plid)
		  plchan <- Msg{PlayerInitLoaded, id, self, iload}

		case CorrectorJoined:
			data, ok := msg.data.(CorrectorJoinedMsg)
	    if !ok { break }
			correctors[data.Id] = data.Chan

		case CorrectorLeft:
			data, ok := msg.data.(string)
	    if !ok { break }
		  delete(correctors, data)
		
		case BuyProb:
			diff, ok := msg.data.(string)
		  if !ok { msg.callback <- Msg{UserError, id, self, "buy"}; break }

			prob, money, remprobs, err := buyProb(conn, id, diff)
			if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() } }

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
			corr <- Msg{BoughtProb, "", self, CorrTicket{id, tname, prob}}

			setTCorr(conn, msg.from, prob.Id, adminid)
			addCorrTicket(conn, adminid, msg.from, prob.Id)

			prob.Answer = "<dobrej pokus>"

			for _, pl := range players {
				pl <- Msg{BoughtProb, id, self, BoughtProbMsg{
					prob, money, remprobs,
				}}
			}

		case SellProb:
			probid, ok := msg.data.(string)
		  if !ok { break }

			money, err := sellProb(conn, id, probid)
			if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() } }

			corr := getTCorr(conn, id, probid)
			correctors[corr] <- Msg{SoldProb, id, self, probid}

			for _, pl := range players {
				pl <- Msg{SoldProb, id, self, IdMoneyMsg{probid, money}}
			}

		case SolveProb:
			data, ok := msg.data.(SolveProbMsg)
		  if !ok { break }

			for _, pl := range players {
				pl <- Msg{SolveProbReq, id, self, data}
			}

			prob := getProb(conn, data.probid)
			if strings.TrimSpace(data.text) == strings.TrimSpace(prob.Answer) {
				money, err := solveProb(conn, id, data.probid)
				if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() } }

				for _, pl := range players {
					pl <- Msg{SolveProb, id, self, IdMoneyMsg{data.probid, money}}
				}
			} else {
				corr := getTCorr(conn, id, data.probid)
				correctors[corr] <- Msg{SolveProb, id, self, data}
			}

		case WriteMsg:
			data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }

			pushTLine(conn, msg.from, data.probid, TLineAtom{MSidePlayer, MTypeText, data.msg, time.Now()})

			for _, pl := range players {
				pl <- Msg{WriteMsg, id, self, data}
			}

			corr := getTCorr(conn, id, data.probid)
			correctors[corr] <- Msg{WriteMsg, id, self, data}

		case CorrGrade:
			data, ok := msg.data.(string)
		  if !ok { break }

			money, err := solveProb(conn, id, data)
			if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() } }

			for _, pl := range players {
				pl <- Msg{SolveProb, id, self, IdMoneyMsg{data, money}}
			}

			corr := getTCorr(conn, id, data)
			correctors[corr] <- Msg{CorrGraded, id, self, TeamProbMsg{id, data}}

		}
	}
}

func playerManager(ws *websocket.Conn, self chan Msg, team chan Msg, id string, teamid string) {
	go func() {
		for {
			var msg map[string]string
			err := ws.ReadJSON(&msg)
			if err != nil {
				self <- Msg{WsError, id, self, err}
				break
			}
			name, ok := msg["name"]
			if !ok { self <- Msg{UserError, id, self, "no name"}; break }
			switch name {

			case "buy":
				diff, ok := msg["diff"]
				if !ok { self <- Msg{UserError, id, self, "no diff"}; break }
				team <- Msg{BuyProb, id, self, diff}

			case "sell":
				probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}; break }
				team <- Msg{SellProb, id, self, probid}

			case "solve":
				probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}; break }
				solution, ok := msg["answer"]
				if !ok { self <- Msg{UserError, id, self, "no answer"}; break }
				team <- Msg{SolveProb, id, self, SolveProbMsg{probid, solution}}

			case "write":
				probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}; break }
				text, ok := msg["text"]
				if !ok { self <- Msg{UserError, id, self, "no text"}; break }
				if len(text) > 50 { self <- Msg{UserError, id, self, "too long"}; break }
				team <- Msg{WriteMsg, id, self, WriteMsgMsg{probid, id, MTypeText, text, false}}

			case "focus":
				probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}; break }
			  team <- Msg{Focus, id, self, PlayerFocusMsg{probid, id}}

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

		case UserError:
			data, ok := msg.data.(string)
			if !ok { break }
			err := ws.WriteJSON(map[string]string{
				"name": "error",
				"error": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case PlayerInitLoaded:
		  data, ok := msg.data.(InitLoad)
			if !ok { break }
			err := ws.WriteJSON(map[string]any{
				"name": "initload",
				"data": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Focus:
		  data, ok := msg.data.(PlayerFocusMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]string{
				"name": "focus",
				"playerid": data.playerid,
				"probid": data.probid,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case WriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
			mtype := "sent"
			if data.admin {
				mtype = "recieved"
			}
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"text": data.msg,
				"type": mtype,
				"solve": false,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case BoughtProb:
		  data, ok := msg.data.(BoughtProbMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]any{
				"name": "bought",
				"prob": data.prob,
				"money": data.money,
				"remprobs": data.remprobs,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SoldProb:
		  data, ok := msg.data.(IdMoneyMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]string{
				"name": "sold",
				"probid": data.id,
				"money": strconv.Itoa(data.money),
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SolvedProb:
		  data, ok := msg.data.(IdMoneyMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]string{
				"name": "solved",
				"probid": data.id,
				"money": strconv.Itoa(data.money),
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SolveProbReq:
		  data, ok := msg.data.(SolveProbMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"text": teamid,
				"type": "sent",
				"solve": false,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }
		  
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

type CorrectorJoinedMsg struct{
	Id string
	Chan chan Msg
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
		  go teamManager(tchan, self, id, getTeamName(conn, id))
		  teams[id] = tchan
			for _, ch := range correctors {
				ch <- Msg{TeamRegister, "", self, CorrectorJoinedMsg{id, tchan}}
			}

		case TeamChanError:
			tchan := make(chan Msg, 10)
		  go teamManager(tchan, self, msg.from, getTeamName(conn, msg.from))
		  teams[msg.from] = tchan

		case CorrectorJoined:
			data, ok := msg.data.(CorrectorJoinedReqMsg)
		  if !ok { break }
			cchan := make(chan Msg, 10)
		  go correctorManager(data.Conn, cchan, self, data.Id, teams)
		  correctors[data.Id] = cchan
			cchan <- Msg{CorrectorInitLoaded, "", self, corrInitLoad(conn, data.Id)}
		  for _, ch := range teams {
				ch <- Msg{CorrectorJoined, "", self, CorrectorJoinedMsg{data.Id, cchan}}
			}

		case CorrectorLeft:
			_, ok := msg.data.(error)
		  if !ok { break }
		  delete(correctors, msg.from)
			tickets := getCorrTickets(conn, msg.from)
		  for _, t := range tickets {
				teamid, probid := parseTicketId(t)
				if getTState(conn, teamid, probid) != OwnedBought { continue }
				prob := getProb(conn, probid)
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
				corr <- Msg{BoughtProb, "", self, CorrTicket{teamid, getTeamName(conn, teamid), prob}}
				setTCorr(conn, msg.from, prob.Id, adminid)
				addCorrTicket(conn, adminid, msg.from, prob.Id)
			}
		  for _, ch := range teams {
				ch <- Msg{CorrectorLeft, "", self, msg.from}
			}

		case PlayerJoinRequest:
			data, ok := msg.data.(PlayerJoinReqMsg)
			if !ok { break }
			team, ok := teams[data.Teamid]
			if !ok { break }
			team <- Msg{PlayerJoinRequest, "", self, data.Conn}
		}
	}
}

func correctorManager(ws *websocket.Conn, self chan Msg, admins chan Msg, id string, teams map[string]chan Msg) {
	go func() {
		for {
			var msg map[string]string
			err := ws.ReadJSON(&msg)
			if err != nil {
				self <- Msg{WsError, id, self, err}
				break
			}
			name, ok := msg["name"]
			if !ok { self <- Msg{UserError, id, self, "no name"}; break }
			switch name {

			case "focus":
				tid, ok := msg["id"]
				if !ok { self <- Msg{UserError, id, self, "not focus id"}}
				self <- Msg{Focus, id, self, tid}

			case "write":
			  tid, ok := msg["teamid"]
				if !ok { self <- Msg{UserError, id, self, "no teamid"}}
			  probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}}
			  text, ok := msg["text"]
				if !ok { self <- Msg{UserError, id, self, "no text"}}
			  self <- Msg{CorrWriteMsg, id, self, WriteMsgMsg{probid, tid, MTypeText, text, true}}
			  self <- Msg{WriteMsg, id, self, WriteMsgMsg{probid, tid, MTypeText, text, true}}

			case "grade":
			  tid, ok := msg["teamid"]
				if !ok { self <- Msg{UserError, id, self, "no teamid"}}
			  probid, ok := msg["probid"]
				if !ok { self <- Msg{UserError, id, self, "no probid"}}
			  self <- Msg{CorrGrade, id, self, TeamProbMsg{tid, probid}}

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

		case UserError:
			data, ok := msg.data.(string)
			if !ok { break }
			err := ws.WriteJSON(map[string]string{"name": "error", "error": data})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case TeamRegister:
		  data, ok := msg.data.(CorrectorJoinedMsg)
			if !ok { break }
			teams[data.Id] = data.Chan

		case CorrectorInitLoaded:
		  data, ok := msg.data.(CorrInitLoad)
			if !ok { break }
		  err := ws.WriteJSON(map[string]any{
				"name": "initload",
				"data": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Focus:
			data, ok := msg.data.(string)
			if !ok { break }
			err := ws.WriteJSON(map[string]string{"name": "focus", "id": data})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case CorrWriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
		  teams[data.teamid] <- Msg{WriteMsg, id, self, data}

		case CorrGrade:
		  data, ok := msg.data.(TeamProbMsg)
		  if !ok { break }
		  teams[data.team] <- Msg{CorrGrade, id, self, data.prob}


		case CorrGraded:
		  data, ok := msg.data.(TeamProbMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]string{
				"name": "graded",
				"probid": data.prob,
				"teamid": data.team,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case WriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
			mtype := "recieved"
			if data.admin {
				mtype = "sent"
			}
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"teamid": data.teamid,
				"text": data.msg,
				"type": mtype,
				"solve": false,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SolveProb:
		  data, ok := msg.data.(SolveProbMsg)
		  if !ok { break }
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"teamid": msg.from,
				"text": data.text,
				"type": "recieved",
				"solve": true,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SoldProb:
		  data, ok := msg.data.(string)
		  if !ok { break }
		  err := ws.WriteJSON(map[string]string{
				"name": "sold",
				"teamid": msg.from,
				"probid": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case BoughtProb:
		  data, ok := msg.data.(CorrTicket)
		  if !ok { break }
			err := ws.WriteJSON(map[string]any{
				"name": "bought",
				"teamid": data.TeamId,
				"teamname": data.TeamName,
				"prob": data.Prob,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		}
	}
}
