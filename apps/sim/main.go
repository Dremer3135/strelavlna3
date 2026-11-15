package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func NewRdbConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

const (
	ArenaUrl = "https://sv.skrat.org"
)

var words []string

func main() {
	data, err := os.ReadFile("words.txt")
	if err != nil {
		panic(err)
	}
	for bts := range bytes.SplitSeq(data, []byte("\n")) {
		words = append(words, string(bts))
	}

	tt, ct := genMockRdb(100, 20, 200)
	fmt.Println(tt)
	fmt.Println(ct)
}

type InitLoad struct {
	TeamId    string                 `json:"teamid"`
	TeamName  string                 `json:"teamname"`
	PlayerId  string                 `json:"playerid"`
	State     string                 `json:"state"`
	Bought    []Prob                 `json:"bought"`
	Sold      []Prob                 `json:"sold"`
	Solved    []Prob                 `json:"solved"`
	TLines    map[string][]TLineAtom `json:"tlines"`
	Money     int                    `json:"money"`
	Start     int                    `json:"start"`
	End       int                    `json:"end"`
	Rank      int                    `json:"rank"`
	RemProbs  map[string]int         `json:"remprobs"`
	BuyCost   map[string]int         `json:"buycost"`
	SellCost  map[string]int         `json:"sellcost"`
	SolveCost map[string]int         `json:"solvecost"`
}

type PlayerState struct {
	state   string
	focused string
	money   int
	bought  map[string]struct{}
	sold    map[string]struct{}
	solved  map[string]struct{}
	buycost map[string]int
}

func playerSimLoop(token string) error {
	conn, _, err := websocket.DefaultDialer.Dial(ArenaUrl+"/ws?token="+token, nil)
	if err != nil {
		return err
	}

	state := PlayerState{}

	pchan := make(chan map[string]any, 10)

	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	go func() {
		for {
			var msg map[string]any
			err := conn.ReadJSON(&msg)
			if err != nil {
				pchan <- nil
				return
			}
			pchan <- msg
		}
	}()

	tick := time.NewTicker(time.Second)

	for {
		select {
		case msg, ok := <-pchan:
			if !ok {
				conn.Close()
				return nil
			}

			name, ok := msg["name"].(string)
			if !ok {
				continue
			}

			switch name {
			case "initload":
				data, ok := msg["data"].(InitLoad)
				if !ok {
					continue
				}
				state.state = data.State
				state.money = data.Money
				state.focused = ""
				for _, prob := range data.Bought {
					state.bought[prob.Id] = struct{}{}
				}
				for _, prob := range data.Sold {
					state.sold[prob.Id] = struct{}{}
				}
				for _, prob := range data.Solved {
					state.solved[prob.Id] = struct{}{}
				}
				state.buycost = data.BuyCost

			case "focus":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				state.focused = probid

			case "bought":
				prob, ok := msg["prob"].(Prob)
				if !ok {
					continue
				}
				money, ok := msg["money"].(int)
				if !ok {
					continue
				}
				state.bought[prob.Id] = struct{}{}
				state.money = money

			case "sold":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				money, ok := msg["money"].(int)
				if !ok {
					continue
				}
				delete(state.bought, probid)
				state.sold[probid] = struct{}{}
				state.money = money

			case "solved":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				money, ok := msg["money"].(int)
				if !ok {
					continue
				}
				delete(state.bought, probid)
				state.solved[probid] = struct{}{}
				state.money = money

			case "start":
				state.state = "running"

			case "end":
				state.state = "after"

			case "results":
				state.state = "results"

			case "paused":
				state.state = "paused"

			case "resumed":
				state.state = "running"

			case "chmoney":
				money, ok := msg["money"].(int)
				if !ok {
					continue
				}
				state.money = money

			}

		case <-tick.C:

			if state.state == "results" {
				if rng.Intn(300) == 0 {
					conn.Close()
					return nil
				}
				continue
			}

			if state.state != "running" {
				continue
			}

			if rng.Intn(5) == 0 {
				fn := rng.Intn(10)
				var cids map[string]struct{}
				switch fn {
				case 0:
					cids = state.sold
				case 1:
					cids = state.solved
				default:
					cids = state.bought
				}
				bids := make([]string, 0, len(cids))
				for id := range cids {
					bids = append(bids, id)
				}
				fprob := bids[rng.Intn(len(bids))]
				conn.WriteJSON(map[string]string{
					"name":   "focus",
					"probid": fprob,
				})
			} else if rng.Intn(30) == 0 {
				msg := ""
				for range rng.Intn(20) + 3 {
					msg += words[rng.Intn(len(words))]
					msg += " "
				}
				msg = strings.TrimSuffix(msg, " ")

				mtype := "message"
				if rng.Intn(20) == 0 {
					mtype = "copy"
				}
				if rng.Intn(20) == 0 {
					mtype = "paste"
				}

				conn.WriteJSON(map[string]string{
					"name":    "write",
					"probid":  state.focused,
					"message": msg,
					"mtype":   mtype,
				})

			} else if rng.Intn(50) == 0 {
				if _, ok := state.bought[state.focused]; !ok {
					continue
				}
				msg := fmt.Sprint(rng.Float64()*10000 - 5000)
				conn.WriteJSON(map[string]string{
					"name":   "solve",
					"probid": state.focused,
					"answer": msg,
				})
			} else if rng.Intn(100) == 0 {
				if _, ok := state.bought[state.focused]; !ok {
					continue
				}
				conn.WriteJSON(map[string]string{
					"name":   "sell",
					"probid": state.focused,
				})
			} else if rng.Intn(30) == 0 {
				diff := ""
				if state.money > state.buycost["C"] {
					diff = "C"
				} else if state.money > state.buycost["B"] {
					diff = "B"
				} else if state.money > state.buycost["A"] {
					diff = "A"
				}
				if diff == "" {
					continue
				}
				conn.WriteJSON(map[string]string{
					"name": "buy",
					"diff": diff,
				})
			}

		}
	}

}

type CorrInitLoad struct {
	BoughtTickets map[string]CorrTicket  `json:"bought_tickets"`
	SolvedTickets map[string]CorrTicket  `json:"solved_tickets"`
	SoldTickets   map[string]CorrTicket  `json:"sold_tickets"`
	TLines        map[string][]TLineAtom `json:"tlines"`
	Start         int                    `json:"start"`
	End           int                    `json:"end"`
	State         string                 `json:"state"`
	Admin         bool                   `json:"admin"`
	Id            string                 `json:"id"`
}

type CorrTicket struct {
	TeamId   string                 `json:"team_id"`
	TeamName string                 `json:"team_name"`
	Prob     Prob                   `json:"prob"`
	TLines   map[string][]TLineAtom `json:"tlines"`
}

type CorrState struct {
	state   string
	focused string
	bought  map[string]struct{}
	sold    map[string]struct{}
	solved  map[string]struct{}
}

func corrSimLoop(token string) error {
	conn, _, err := websocket.DefaultDialer.Dial(ArenaUrl+"/ws/corr?token="+token, nil)
	if err != nil {
		return err
	}

	state := CorrState{
		bought: make(map[string]struct{}),
		sold:   make(map[string]struct{}),
		solved: make(map[string]struct{}),
	}

	cchan := make(chan map[string]any, 10)

	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	go func() {
		for {
			var msg map[string]any
			err := conn.ReadJSON(&msg)
			if err != nil {
				cchan <- nil
				return
			}
			cchan <- msg
		}
	}()

	tick := time.NewTicker(time.Second)

	for {
		select {
		case msg, ok := <-cchan:
			if !ok || msg == nil {
				conn.Close()
				return nil
			}

			name, ok := msg["name"].(string)
			if !ok {
				continue
			}

			switch name {
			case "initload":
				data, ok := msg["data"].(CorrInitLoad)
				if !ok {
					continue
				}
				state.state = data.State
				state.focused = ""
				for tickid := range data.BoughtTickets {
					state.bought[tickid] = struct{}{}
				}
				for tickid := range data.SoldTickets {
					state.sold[tickid] = struct{}{}
				}
				for tickid := range data.SolvedTickets {
					state.solved[tickid] = struct{}{}
				}

			case "focus":
				id, ok := msg["id"].(string)
				if !ok {
					continue
				}
				state.focused = id

			case "bought":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				prob, ok := msg["prob"].(Prob)
				if !ok {
					continue
				}
				state.bought[teamid+":"+prob.Id] = struct{}{}

			case "sold":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				tickid := teamid + ":" + probid
				delete(state.bought, tickid)
				state.sold[tickid] = struct{}{}

			case "solved":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				tickid := teamid + ":" + probid
				delete(state.bought, tickid)
				state.solved[tickid] = struct{}{}

			case "start":
				state.state = "running"

			case "end":
				state.state = "after"

			case "results":
				state.state = "results"

			case "paused":
				state.state = "paused"

			case "resumed":
				state.state = "running"
			}

		case <-tick.C:

			if state.state == "results" {
				if rng.Intn(300) == 0 {
					conn.Close()
					return nil
				}
				continue
			}

			if state.state != "running" {
				continue
			}

			if rng.Intn(5) == 0 {
				fn := rng.Intn(10)
				var cids map[string]struct{}
				switch fn {
				case 0:
					cids = state.sold
				case 1:
					cids = state.solved
				default:
					cids = state.bought
				}
				if len(cids) == 0 {
					continue
				}
				bids := make([]string, 0, len(cids))
				for id := range cids {
					bids = append(bids, id)
				}
				ftick := bids[rng.Intn(len(bids))]
				conn.WriteJSON(map[string]string{
					"name": "focus",
					"id":   ftick,
				})
			} else if rng.Intn(30) == 0 {
				if state.focused == "" {
					continue
				}
				msg := ""
				for range rng.Intn(20) + 3 {
					msg += words[rng.Intn(len(words))]
					msg += " "
				}
				msg = strings.TrimSuffix(msg, " ")

				conn.WriteJSON(map[string]string{
					"name":    "write",
					"id":      state.focused,
					"message": msg,
				})
			} else if rng.Intn(50) == 0 {
				if _, ok := state.bought[state.focused]; !ok {
					continue
				}
				decision := "incorrect"
				if rng.Intn(2) == 0 {
					decision = "correct"
				}
				conn.WriteJSON(map[string]string{
					"name":     "grade",
					"id":       state.focused,
					"decision": decision,
				})
			}

		}
	}
}

const idChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func randId(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = idChars[rng.Intn(len(idChars))]
	}
	return string(b)
}

func randWords(rng *rand.Rand, minW, maxW int) string {
	if len(words) == 0 {
		return randId(rng, 8)
	}
	count := minW
	if maxW > minW {
		count += rng.Intn(maxW - minW + 1)
	}
	parts := make([]string, 0, count)
	for range count {
		w := strings.TrimSpace(words[rng.Intn(len(words))])
		if w != "" {
			parts = append(parts, w)
		}
	}
	return strings.Join(parts, " ")
}

func genMockRdb(numTeams, numCorrs, numProbs int) ([]string, []string) {
	rdb := NewRdbConn()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if err := rdb.FlushDB(ctx).Err(); err != nil {
		panic(err)
	}

	buyPrices := map[string]int{"A": 10, "B": 30, "C": 80}
	sellPrices := map[string]int{"A": 5, "B": 15, "C": 40}
	solvePrices := map[string]int{"A": 20, "B": 60, "C": 160}

	for _, d := range diffs {
		setPrice(rdb, PriceBuy, d, buyPrices[d])
		setPrice(rdb, PriceSell, d, sellPrices[d])
		setPrice(rdb, PriceSolve, d, solvePrices[d])
	}

	teamIds := make([]string, numTeams)
	teamTokens := make([]string, numTeams)
	for i := range numTeams {
		tid := randId(rng, 15)
		token := randId(rng, 5)
		teamIds[i] = tid
		teamTokens[i] = token

		setMoney(rdb, tid, 140)
		setPlayToken(rdb, token, tid)
		setTeamName(rdb, tid, randWords(rng, 2, 3))
	}
	if numTeams > 0 {
		setTeams(rdb, teamIds)
	}

	corrIds := make([]string, numCorrs)
	corrTokens := make([]string, numCorrs)
	for i := range numCorrs {
		cid := randId(rng, 15)
		token := randId(rng, 5)
		corrIds[i] = cid
		corrTokens[i] = token

		setCorrToken(rdb, token, cid)
		setCorrAdmin(rdb, cid, true)
	}

	for i := range numProbs {
		pid := randId(rng, 15)
		diff := diffs[i%len(diffs)]

		queue := make([]string, len(corrIds))
		copy(queue, corrIds)
		rng.Shuffle(len(queue), func(a, b int) {
			queue[a], queue[b] = queue[b], queue[a]
		})

		prob := Prob{
			Id:       pid,
			Name:     randWords(rng, 2, 4),
			Diff:     diff,
			Text:     randWords(rng, 15, 40),
			Answer:   fmt.Sprintf("%.2f", rng.Float64()*1000-500),
			Code:     "",
			Auto:     false,
			Infinite: diff == "A",
			Queue:    queue,
			Images:   []string{},
		}
		setProb(rdb, prob)

		for _, tid := range teamIds {
			addOwnedProb(rdb, tid, OwnedFree, diff, pid)
		}
	}

	setState(rdb, StateRunning)
	setStart(rdb, time.Now())
	setEnd(rdb, time.Now().Add(2*time.Hour))

	if _, err := rdb.Save(ctx).Result(); err != nil {
		panic(err)
	}

	return teamTokens, corrTokens
}

func getTokensFromRdb(rdb *redis.Client) (teamTokens []string, corrTokens []string) {
	pkeys, err := rdb.Keys(ctx, "playtoken:*").Result()
	if err != nil {
		panic(err)
	}
	for _, k := range pkeys {
		teamTokens = append(teamTokens, strings.TrimPrefix(k, "playtoken:"))
	}

	ckeys, err := rdb.Keys(ctx, "corrtoken:*").Result()
	if err != nil {
		panic(err)
	}
	for _, k := range ckeys {
		corrTokens = append(corrTokens, strings.TrimPrefix(k, "corrtoken:"))
	}
	return
}

func runFinalTest(teamTokens []string, corrTokens []string, playersPerTeam int) {
	if len(teamTokens) == 0 && len(corrTokens) == 0 {
		rdb := NewRdbConn()
		teamTokens, corrTokens = getTokensFromRdb(rdb)
	}

	var wg sync.WaitGroup

	for _, ctoken := range corrTokens {
		wg.Add(1)
		go func(tok string) {
			defer wg.Done()
			if err := corrSimLoop(tok); err != nil {
				fmt.Printf("corr %s error: %v\n", tok, err)
			}
		}(ctoken)
	}

	for _, ttoken := range teamTokens {
		for range playersPerTeam {
			wg.Add(1)
			go func(tok string) {
				defer wg.Done()
				if err := playerSimLoop(tok); err != nil {
					fmt.Printf("player %s error: %v\n", tok, err)
				}
			}(ttoken)
		}
	}

	wg.Wait()
}
