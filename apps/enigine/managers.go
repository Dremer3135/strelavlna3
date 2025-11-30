package main

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
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

	Start
	End
	Results
)

const (
	writeWait = 10 * time.Second
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
	time time.Time
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

type GradeProbMsg struct {
	team string
	prob string
	decision string
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

type ResultsMsg struct {
	money int
	rank int
}

var whiteSpaceRegex = regexp.MustCompile(`\\S+`)

func checkAnswer(cans, tans string) bool {
  cans = whiteSpaceRegex.ReplaceAllLiteralString(cans, "")
  tans = whiteSpaceRegex.ReplaceAllLiteralString(tans, "")
	cans = strings.ReplaceAll(cans, ",", ".")
	tans = strings.ReplaceAll(tans, ",", ".")
	cf, errc := strconv.ParseFloat(cans, 64)
	tf, errt := strconv.ParseFloat(tans, 64)
	if errc == nil && errt == nil {
	  if tf > cf * 0.99 && tf < cf * 1.01 { return true }
	}
	if cans == tans { return true }
  return false
}

func teamManager(self chan Msg, admins chan Msg, id string, tname string) {
	players := map[string]chan Msg{}
	conn := NewRdbConn()
	correctors := map[string]chan Msg{}
	for {
		msg, ok := <- self
		fmt.Printf("team: %#v\n", msg)
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
		  plchan := make(chan Msg, 1000)
		  go playerManager(data, plchan, self, plid, id)
		  players[plid] = plchan

			iload, err := initLoad(conn, id, plid)
			if err != nil {
				plchan <- Msg{ServerError, id, self, "could not load initial data"}
				break
			}
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

			state, err := getState(conn)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if state != StateRunning { msg.callback <- Msg{UserError, id, self, "not running"}; break }

			prob, money, remprobs, err := buyProb(conn, id, diff)
			if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() }; break }

			adminid := ""
			for _, a := range prob.Queue {
				_, ok := correctors[a]
				if ok {
					adminid = a
					break
				}
			} 
			if adminid == "" {
				msg.callback <- Msg{UserError, id, self, "no corrector"}
				break
			}
			corr := correctors[adminid]
			corr <- Msg{BoughtProb, "", self, CorrTicket{id, tname, prob}}

			if err := setTCorr(conn, id, prob.Id, adminid); err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if err := addCorrTicket(conn, adminid, id, prob.Id); err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}

			prob.Answer = "<dobrej pokus>"

			for _, pl := range players {
				select {
				case pl <- Msg{BoughtProb, id, self, BoughtProbMsg{
					prob, money, remprobs,
				}}:
				default:
				}
			}

		case SellProb:
			probid, ok := msg.data.(string)
		  if !ok { break }

			state, err := getState(conn)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if state != StateRunning { msg.callback <- Msg{UserError, id, self, "not running"}; break }

			money, err := sellProb(conn, id, probid)
			if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() }; break }

			corr, err := getTCorr(conn, id, probid)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if corr != "" {
				correctors[corr] <- Msg{SoldProb, id, self, probid}
			}

			for _, pl := range players {
				select {
				case pl <- Msg{SoldProb, id, self, IdMoneyMsg{probid, money}}:
				default:
				}
			}

		case SolveProb:
			data, ok := msg.data.(SolveProbMsg)
		  if !ok { break }

			state, err := getState(conn)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if state != StateRunning { msg.callback <- Msg{UserError, id, self, "not running"}; break }
			// for _, pl := range players {
			// 	pl <- Msg{SolveProbReq, id, self, data}
			// }

			go func(){
				self <- Msg{WriteMsg, id, self, WriteMsgMsg{data.probid, id, MTypeSolve, data.text, false, time.Now()}}
			}()

			prob, err := getProb(conn, data.probid)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if checkAnswer(data.text, prob.Answer) {
				money, err := solveProb(conn, id, data.probid)
				if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() }; break }

				for _, pl := range players {
					select {
					case pl <- Msg{SolvedProb, id, self, IdMoneyMsg{data.probid, money}}:
					default:
					}
				}

				go func(){
					self <- Msg{WriteMsg, id, self, WriteMsgMsg{data.probid, id, MTypeGrade, DecisionCorrect, true, time.Now()}}
				}()

				corr, err := getTCorr(conn, id, data.probid)
				if err != nil {
					msg.callback <- Msg{ServerError, id, self, "db error"}
					break
				}
				if corr != "" {
					correctors[corr] <- Msg{CorrGraded, id, self, TeamProbMsg{id, data.probid}}
				}
			} // else {
			// 	corr := getTCorr(conn, id, data.probid)
			// 	correctors[corr] <- Msg{SolveProb, id, self, data}
			// }

		case WriteMsg:
			data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }

			side := MSidePlayer
			if data.admin {
				side = MSideAdmin
			}
			if err := pushTLine(conn, id, data.probid, TLineAtom{side, data.mtype, data.msg, time.Now()}); err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}

			for _, pl := range players {
				select {
				case pl <- Msg{WriteMsg, id, self, data}:
				default:
				}
			}

			corr, err := getTCorr(conn, id, data.probid)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}
			if corr != "" {
				correctors[corr] <- Msg{WriteMsg, id, self, data}
			}

		case CorrGrade:
			data, ok := msg.data.(GradeProbMsg)
		  if !ok { break }

			corr, err := getTCorr(conn, id, data.prob)
			if err != nil {
				msg.callback <- Msg{ServerError, id, self, "db error"}
				break
			}

			if data.decision == DecisionCorrect {
				money, err := solveProb(conn, id, data.prob)
				if err != nil { msg.callback <- Msg{UserError, id, self, err.Error() }; break }

				for _, pl := range players {
					select {
					case pl <- Msg{SolvedProb, id, self, IdMoneyMsg{data.prob, money}}:
					default:
					}
				}
				if corr != "" {
					correctors[corr] <- Msg{SolvedProb, id, self, data.prob}
				}
			}

			// pushTLine(conn, data.team, data.prob, TLineAtom{MSideAdmin, MTypeGrade, data.decision, time.Now()})
			// msgs := Msg{WriteMsg, id, self, WriteMsgMsg{data.prob, data.team, MTypeGrade, data.decision, true, time.Now()}}
			// for _, pl := range players {
			// 	pl <- msgs
			// }

			go func(){
				self <- Msg{WriteMsg, id, self, WriteMsgMsg{data.prob, data.team, MTypeGrade, data.decision, true, time.Now()}}
			}()

			if corr != "" {
				correctors[corr] <- Msg{CorrGraded, id, self, TeamProbMsg{id, data.prob}}
			}

		case Start:
			for _, pl := range players {
				select {
				case pl <- Msg{Start, id, self, nil}:
				default:
				}
			}

		case End:
			for _, pl := range players {
				select {
				case pl <- Msg{End, id, self, nil}:
				default:
				}
			}

		case Results:
			data, ok := msg.data.(TeamMoneyResult)
		  if !ok { break }
			for _, pl := range players {
				select {
				case pl <- Msg{Results, id, self, data}:
				default:
				}
			}

		case Focus:
			data, ok := msg.data.(PlayerFocusMsg)
		  if !ok { break }
		  for _, pl := range players {
				select {
				case pl <- Msg{Focus, id, self, data}:
				default:
				}
			}

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
				text, ok := msg["message"]
				if !ok { self <- Msg{UserError, id, self, "no text"}; break }
				mtype, ok := msg["mtype"]
				if !ok { self <- Msg{UserError, id, self, "no mtype"}; break }
				if mtype != MTypeText && mtype != MTypeGif && mtype != MTypeCopy && mtype != MTypePaste && mtype != MTypeFocus { self <- Msg{UserError, id, self, "invalid mtype"}; break }
				// if len(text) > 50 { self <- Msg{UserError, id, self, "too long"}; break }
				team <- Msg{WriteMsg, id, self, WriteMsgMsg{probid, teamid, mtype, text, false, time.Now()}}

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
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "error",
				"error": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case ServerError:
			data, ok := msg.data.(string)
			if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "error",
				"error": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case PlayerInitLoaded:
		  data, ok := msg.data.(InitLoad)
			fmt.Printf("%v %v\n", ok, data)
			if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "initload",
				"data": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }
			fmt.Printf("err: %v\n", err)

		case Focus:
		  data, ok := msg.data.(PlayerFocusMsg)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "focus",
				"playerid": data.playerid,
				"probid": data.probid,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case WriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
			rtype := "sent"
			if data.admin {
				rtype = "recieved"
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"text": data.msg,
				"origin": rtype,
				"type": data.mtype,
				"time": data.time.UnixMilli(),
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case BoughtProb:
		  data, ok := msg.data.(BoughtProbMsg)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
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
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "sold",
				"probid": data.id,
				"money": strconv.Itoa(data.money),
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SolvedProb:
		  data, ok := msg.data.(IdMoneyMsg)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "solved",
				"probid": data.id,
				"money": strconv.Itoa(data.money),
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		// case SolveProbReq:
		//   data, ok := msg.data.(SolveProbMsg)
		//   if !ok { break }
		// 	err := ws.WriteJSON(map[string]any{
		// 		"name": "written",
		// 		"probid": data.probid,
		// 		"text": teamid,
		// 		"type": "sent",
		// 		"solve": false,
		// 		"time": data.time,
		// 	})
		// 	if err != nil { self <- Msg{WsError, id, self, err} }

		case Start:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "start",
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case End:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "end",
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Results:
		  data, ok := msg.data.(TeamMoneyResult)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "results",
				"money": data.Money,
				"rank": data.Rank,
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
		fmt.Printf("admin: %#v\n", msg)
		switch msg.tp {

		case TeamRegister:
			id, ok := msg.data.(string)
		  if !ok { break }
			tchan := make(chan Msg, 1000)
			tname, err := getTeamName(conn, id)
			if err != nil {
				fmt.Printf("admin error: %v\n", err)
				continue
			}
		  go teamManager(tchan, self, id, tname)
		  teams[id] = tchan
			for _, ch := range correctors {
				select {
				case ch <- Msg{TeamRegister, "", self, CorrectorJoinedMsg{id, tchan}}:
				default:
				}
			}

		case TeamChanError:
			tchan := make(chan Msg, 1000)
			tname, err := getTeamName(conn, msg.from)
			if err != nil {
				fmt.Printf("admin error: %v\n", err)
				continue
			}
		  go teamManager(tchan, self, msg.from, tname)
		  teams[msg.from] = tchan

		case CorrectorJoined:
			data, ok := msg.data.(CorrectorJoinedReqMsg)
		  if !ok { break }
			cchan := make(chan Msg, 1000)
		  go correctorManager(data.Conn, cchan, self, data.Id, teams)
		  correctors[data.Id] = cchan
			iload, err := corrInitLoad(conn, data.Id)
			if err != nil {
				cchan <- Msg{ServerError, "", self, "could not load initial data"}
				break
			}
			cchan <- Msg{CorrectorInitLoaded, "", self, iload}
		  for _, ch := range teams {
				select {
				case ch <- Msg{CorrectorJoined, "", self, CorrectorJoinedMsg{data.Id, cchan}}:
				default:
				}
			}

		case CorrectorLeft:
			_, ok := msg.data.(error)
		  if !ok { break }
		  for _, ch := range teams {
				select {
				case ch <- Msg{CorrectorLeft, "", self, msg.from}:
				default:
				}
			}
		  delete(correctors, msg.from)
			tickets, err := getCorrTickets(conn, msg.from)
			if err != nil {
				fmt.Printf("admin error: %v\n", err)
				continue
			}
		  for _, t := range tickets {
				teamid, probid := parseTicketId(t)
				tstate, err := getTState(conn, teamid, probid)
				if err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
				if tstate != OwnedBought { continue }
				prob, err := getProb(conn, probid)
				if err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
				adminid := ""
				for _, a := range prob.Queue {
					_, ok := correctors[a]
					if ok {
						adminid = a
						break
					}
				} 
				if adminid == "" { continue }
				corr := correctors[adminid]
				tname, err := getTeamName(conn, teamid)
				if err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
				corr <- Msg{BoughtProb, "", self, CorrTicket{teamid, tname, prob}}
				if err := setTCorr(conn, teamid, prob.Id, adminid); err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
				if err := addCorrTicket(conn, adminid, teamid, probid); err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
			}
			if err := clearCorrTickets(conn, msg.from); err != nil {
				fmt.Printf("admin error: %v\n", err)
			}

		case PlayerJoinRequest:
			data, ok := msg.data.(PlayerJoinReqMsg)
			if !ok { break }
			team, ok := teams[data.Teamid]
			if !ok { break }
			team <- Msg{PlayerJoinRequest, "", self, data.Conn}

		case Start:
		  if err := setState(conn, StateRunning); err != nil {
				fmt.Printf("admin error: %v\n", err)
				continue
			}
		  for _, team := range teams {
				select {
				case team <- Msg{Start, "", self, nil}:
				default:
				}
			}
		  for _, corr := range correctors {
				select {
				case corr <- Msg{Start, "", self, nil}:
				default:
				}
			}
		
		case End:
		  if err := setState(conn, StateAfter); err != nil {
				fmt.Printf("admin error: %v\n", err)
				continue
			}
		  for _, team := range teams {
				select {
				case team <- Msg{End, "", self, nil}:
				default:
				}
			}
		  for _, corr := range correctors {
				select {
				case corr <- Msg{End, "", self, nil}:
				default:
				}
			}

		case Results:
			moneys := make(map[int][]string)
			for id := range teams {
				money, err := getMoney(conn, id)
				if err != nil {
					fmt.Printf("admin error: %v\n", err)
					continue
				}
				moneys[money] = append(moneys[money], id) 
			}
			tmrl := make([]ManyTeamMoneyResult, 0, len(moneys))
			for k, v := range moneys {
				tmrl = append(tmrl, ManyTeamMoneyResult{v, k})
			}
			slices.SortFunc(tmrl, func(a, b ManyTeamMoneyResult) int {
				if a.money < b.money { return 1 }
				if a.money > b.money { return -1 }
				return 0
			})
			resmap := make(map[string]TeamMoneyResult)
			i := 0
			for _, tmr := range tmrl {
				rank := i + 1
				for _, id := range tmr.ids {
					if err := setRank(conn, id, rank); err != nil {
						fmt.Printf("admin error: %v\n", err)
						continue
					}
					tname, err := getTeamName(conn, id)
					if err != nil {
						fmt.Printf("admin error: %v\n", err)
						continue
					}
					resmap[id] = TeamMoneyResult{tmr.money, rank, tname}
					select {
						case teams[id] <- Msg{Results, "", self, TeamMoneyResult{tmr.money, rank, ""}}:
						default:
					}
				}
				i += len(tmr.ids)
			}
			for _, corr := range correctors {
				select {
					case corr <- Msg{Results, "", self, resmap}:
					default:
				}
			}

		}
	}
}

type TeamMoneyResult struct {
	Money int `json:"money"`
	Rank int `json:"rank"`
	Name string `json:"name"`
}

type ManyTeamMoneyResult struct {
	ids []string
	money int
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
			  tid, ok := msg["id"]
				if !ok { self <- Msg{UserError, id, self, "no teamid"}}
			  text, ok := msg["message"]
				if !ok { self <- Msg{UserError, id, self, "no text"}}
				teamid, probid := parseTicketId(tid)
			  self <- Msg{CorrWriteMsg, id, self, WriteMsgMsg{probid, teamid, MTypeText, text, true, time.Now()}}
			  // self <- Msg{WriteMsg, id, self, WriteMsgMsg{probid, teamid, MTypeText, text, true, time.Now()}}

			case "grade":
			  tid, ok := msg["id"]
				if !ok { self <- Msg{UserError, id, self, "no teamid"}}
			  decision, ok := msg["decision"]
				if !ok { self <- Msg{UserError, id, self, "no decision"}}
				teamid, probid := parseTicketId(tid)
			  self <- Msg{CorrGrade, id, self, GradeProbMsg{teamid, probid, decision}}

			case "start":
			  admins <- Msg{Start, id, self, nil}

			case "end":
			  admins <- Msg{End, id, self, nil}

			case "results":
			  admins <- Msg{Results, id, self, nil}

			}
		}
	}()
	chanloop: for {
		msg, ok := <- self
		fmt.Printf("corr: %#v\n", msg)
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
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{"name": "error", "error": data})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case ServerError:
			data, ok := msg.data.(string)
			if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{"name": "error", "error": data})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case TeamRegister:
		  data, ok := msg.data.(CorrectorJoinedMsg)
			if !ok { break }
			teams[data.Id] = data.Chan

		case CorrectorInitLoaded:
		  data, ok := msg.data.(CorrInitLoad)
			if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
		  err := ws.WriteJSON(map[string]any{
				"name": "initload",
				"data": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Focus:
			data, ok := msg.data.(string)
			if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{"name": "focus", "id": data})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case CorrWriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
		  teams[data.teamid] <- Msg{WriteMsg, id, self, data}

		case CorrGrade:
		  data, ok := msg.data.(GradeProbMsg)
		  if !ok { break }
		  teams[data.team] <- Msg{CorrGrade, id, self, data}

		case CorrGraded:
		  data, ok := msg.data.(TeamProbMsg)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]string{
				"name": "graded",
				"probid": data.prob,
				"teamid": data.team,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case WriteMsg:
		  data, ok := msg.data.(WriteMsgMsg)
		  if !ok { break }
			rtype := "recieved"
			if data.admin {
				rtype = "sent"
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "written",
				"probid": data.probid,
				"teamid": data.teamid,
				"text": data.msg,
				"type": data.mtype,
				"time": data.time.UnixMilli(),
				"origin": rtype,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		// case SolveProb:
		//   data, ok := msg.data.(SolveProbMsg)
		//   if !ok { break }
		// 	err := ws.WriteJSON(map[string]any{
		// 		"name": "written",
		// 		"probid": data.probid,
		// 		"teamid": msg.from,
		// 		"text": data.text,
		// 		"origin": "recieved",
		// 		"type": MTypeSolve,
		// 		"time": time.Now().UnixMilli(),
		// 	})
		// 	if err != nil { self <- Msg{WsError, id, self, err} }

		case SoldProb:
		  data, ok := msg.data.(string)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
		  err := ws.WriteJSON(map[string]string{
				"name": "sold",
				"teamid": msg.from,
				"probid": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case SolvedProb:
		  data, ok := msg.data.(string)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
		  err := ws.WriteJSON(map[string]string{
				"name": "solved",
				"teamid": msg.from,
				"probid": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case BoughtProb:
		  data, ok := msg.data.(CorrTicket)
		  if !ok { break }
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(map[string]any{
				"name": "bought",
				"teamid": data.TeamId,
				"teamname": data.TeamName,
				"prob": data.Prob,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Start:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
		  err := ws.WriteJSON(map[string]string{
				"name": "start",
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		case Results:
		  data, ok := msg.data.(map[string]TeamMoneyResult)
		  if !ok { break }
		  err := ws.WriteJSON(map[string]any{
				"name": "results",
				"data": data,
			})
			if err != nil { self <- Msg{WsError, id, self, err} }

		}
	}
}
