package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
	ArenaUrl = "wss://sv.skrat.org"
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

	// tt, ct := genMockRdb(10000, 2000, 2000)
	// fmt.Printf("%#v\n:==-==:\n%#v\n", tt, ct)

	runFinalTest(playtoks[:500], corrtoks[:100], 5)

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
	playerid string
	state    string
	focused  string
	money    int
	bought   map[string]struct{}
	sold     map[string]struct{}
	solved   map[string]struct{}
	ungraded map[string]struct{}
	buycost  map[string]int
}

type SimStats struct {
	activePlayers   atomic.Int64
	activeCorrs     atomic.Int64
	playerFocuses   atomic.Int64
	playerWrites    atomic.Int64
	playerSolves    atomic.Int64
	playerSells     atomic.Int64
	playerBuys      atomic.Int64
	probsBought     atomic.Int64
	probsSold       atomic.Int64
	probsSolved     atomic.Int64
	corrFocuses     atomic.Int64
	corrWrites      atomic.Int64
	corrGrades      atomic.Int64
	ungradedTickets atomic.Int64
	serverErrors    atomic.Int64
	wsDisconnects   atomic.Int64
}

var stats SimStats

func redecode(v any, dst any) bool {
	b, err := json.Marshal(v)
	if err != nil {
		return false
	}
	return json.Unmarshal(b, dst) == nil
}

func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case string:
		i, err := strconv.Atoi(x)
		return i, err == nil
	default:
		return 0, false
	}
}

func isTLineUngraded(tline []TLineAtom) bool {
	lastAns := -1
	lastGrade := -1
	for i, a := range tline {
		if a.Mtype == "answer" || a.Mtype == "solve" {
			lastAns = i
		} else if a.Mtype == "grade" {
			lastGrade = i
		}
	}
	return lastAns > lastGrade
}

func playerSimLoop(token string) error {
	conn, _, err := websocket.DefaultDialer.Dial(ArenaUrl+"/ws?token="+token, nil)
	if err != nil {
		return err
	}
	stats.activePlayers.Add(1)
	defer stats.activePlayers.Add(-1)

	state := PlayerState{
		bought:   make(map[string]struct{}),
		sold:     make(map[string]struct{}),
		solved:   make(map[string]struct{}),
		ungraded: make(map[string]struct{}),
		buycost:  make(map[string]int),
	}

	pchan := make(chan map[string]any, 10)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		for {
			var msg map[string]any
			err := conn.ReadJSON(&msg)
			if err != nil {
				stats.wsDisconnects.Add(1)
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
			if !ok || msg == nil {
				conn.Close()
				return nil
			}

			name, ok := msg["name"].(string)
			if !ok {
				continue
			}

			switch name {
			case "error":
				stats.serverErrors.Add(1)
				fmt.Printf("[player %s] server error: %v\n", token, msg["error"])

			case "initload":
				var data InitLoad
				if !redecode(msg["data"], &data) {
					continue
				}
				state.playerid = data.PlayerId
				state.state = data.State
				state.money = data.Money
				state.focused = ""
				for _, prob := range data.Bought {
					state.bought[prob.Id] = struct{}{}
					if isTLineUngraded(data.TLines[prob.Id]) {
						state.ungraded[prob.Id] = struct{}{}
					}
				}
				for _, prob := range data.Sold {
					state.sold[prob.Id] = struct{}{}
				}
				for _, prob := range data.Solved {
					state.solved[prob.Id] = struct{}{}
				}
				state.buycost = data.BuyCost
				fmt.Printf("[player %s/%s] initload: state=%s money=%d bought=%d solved=%d sold=%d\n",
					token, data.PlayerId, data.State, data.Money, len(data.Bought), len(data.Solved), len(data.Sold))

			case "focus":
				playerid, ok := msg["playerid"].(string)
				if !ok || playerid != state.playerid {
					continue
				}
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				state.focused = probid

			case "written":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				mtype, ok := msg["type"].(string)
				if !ok {
					continue
				}
				if mtype == "answer" {
					state.ungraded[probid] = struct{}{}
				} else if mtype == "grade" {
					delete(state.ungraded, probid)
				}

			case "bought":
				var prob Prob
				if !redecode(msg["prob"], &prob) {
					continue
				}
				money, ok := toInt(msg["money"])
				if !ok {
					continue
				}
				state.bought[prob.Id] = struct{}{}
				state.money = money
				stats.probsBought.Add(1)

			case "sold":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				money, ok := toInt(msg["money"])
				if !ok {
					continue
				}
				delete(state.bought, probid)
				delete(state.ungraded, probid)
				state.sold[probid] = struct{}{}
				state.money = money
				stats.probsSold.Add(1)

			case "solved":
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				money, ok := toInt(msg["money"])
				if !ok {
					continue
				}
				delete(state.bought, probid)
				delete(state.ungraded, probid)
				state.solved[probid] = struct{}{}
				state.money = money
				stats.probsSolved.Add(1)

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
				money, ok := toInt(msg["money"])
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
				if len(cids) == 0 {
					continue
				}
				bids := make([]string, 0, len(cids))
				for id := range cids {
					bids = append(bids, id)
				}
				fprob := bids[rng.Intn(len(bids))]
				stats.playerFocuses.Add(1)
				conn.WriteJSON(map[string]string{
					"name":   "focus",
					"probid": fprob,
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

				mtype := "message"
				if rng.Intn(20) == 0 {
					mtype = "copy"
				}
				if rng.Intn(20) == 0 {
					mtype = "paste"
				}

				stats.playerWrites.Add(1)
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
				if _, ungraded := state.ungraded[state.focused]; ungraded {
					continue
				}
				state.ungraded[state.focused] = struct{}{}
				msg := fmt.Sprint(rng.Float64()*10000 - 5000)
				stats.playerSolves.Add(1)
				conn.WriteJSON(map[string]string{
					"name":   "solve",
					"probid": state.focused,
					"answer": msg,
				})
			} else if rng.Intn(100) == 0 {
				if _, ok := state.bought[state.focused]; !ok {
					continue
				}
				stats.playerSells.Add(1)
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
				stats.playerBuys.Add(1)
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
	state    string
	focused  string
	bought   map[string]struct{}
	sold     map[string]struct{}
	solved   map[string]struct{}
	ungraded map[string]struct{}
}

func (s *CorrState) addUngraded(tickid string) {
	if _, exists := s.ungraded[tickid]; !exists {
		s.ungraded[tickid] = struct{}{}
		stats.ungradedTickets.Add(1)
	}
}

func (s *CorrState) removeUngraded(tickid string) {
	if _, exists := s.ungraded[tickid]; exists {
		delete(s.ungraded, tickid)
		stats.ungradedTickets.Add(-1)
	}
}

func corrSimLoop(token string) error {
	conn, _, err := websocket.DefaultDialer.Dial(ArenaUrl+"/ws/corr?token="+token, nil)
	if err != nil {
		return err
	}
	stats.activeCorrs.Add(1)
	defer stats.activeCorrs.Add(-1)

	state := CorrState{
		bought:   make(map[string]struct{}),
		sold:     make(map[string]struct{}),
		solved:   make(map[string]struct{}),
		ungraded: make(map[string]struct{}),
	}
	defer func() {
		stats.ungradedTickets.Add(-int64(len(state.ungraded)))
	}()

	cchan := make(chan map[string]any, 10)

	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	go func() {
		for {
			var msg map[string]any
			err := conn.ReadJSON(&msg)
			if err != nil {
				stats.wsDisconnects.Add(1)
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
			case "error":
				stats.serverErrors.Add(1)
				fmt.Printf("[corr %s] server error: %v\n", token, msg["error"])

			case "initload":
				var data CorrInitLoad
				if !redecode(msg["data"], &data) {
					continue
				}
				state.state = data.State
				state.focused = ""
				for tickid := range data.BoughtTickets {
					state.bought[tickid] = struct{}{}
					if isTLineUngraded(data.TLines[tickid]) {
						state.addUngraded(tickid)
					}
				}
				for tickid := range data.SoldTickets {
					state.sold[tickid] = struct{}{}
				}
				for tickid := range data.SolvedTickets {
					state.solved[tickid] = struct{}{}
				}
				fmt.Printf("[corr %s] initload: state=%s bought=%d solved=%d sold=%d\n",
					token, data.State, len(data.BoughtTickets), len(data.SolvedTickets), len(data.SoldTickets))

			case "focus":
				id, ok := msg["id"].(string)
				if !ok {
					continue
				}
				state.focused = id

			case "written":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				mtype, ok := msg["type"].(string)
				if !ok {
					continue
				}
				tickid := teamid + ":" + probid
				if mtype == "answer" {
					if _, isBought := state.bought[tickid]; isBought {
						state.addUngraded(tickid)
					}
				} else if mtype == "grade" {
					state.removeUngraded(tickid)
				}

			case "graded":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				probid, ok := msg["probid"].(string)
				if !ok {
					continue
				}
				state.removeUngraded(teamid + ":" + probid)

			case "bought":
				teamid, ok := msg["teamid"].(string)
				if !ok {
					continue
				}
				var prob Prob
				if !redecode(msg["prob"], &prob) {
					continue
				}
				tickid := teamid + ":" + prob.Id
				state.bought[tickid] = struct{}{}
				var tlines map[string][]TLineAtom
				if redecode(msg["tlines"], &tlines) && isTLineUngraded(tlines[prob.Id]) {
					state.addUngraded(tickid)
				}

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
				state.removeUngraded(tickid)
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
				state.removeUngraded(tickid)
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

			if rng.Intn(3) == 0 {
				var cids map[string]struct{}
				if len(state.ungraded) > 0 {
					cids = state.ungraded
				} else {
					fn := rng.Intn(50)
					switch fn {
					case 0:
						cids = state.sold
					case 1:
						cids = state.solved
					default:
						cids = state.bought
					}
				}
				if len(cids) == 0 {
					continue
				}
				bids := make([]string, 0, len(cids))
				for id := range cids {
					bids = append(bids, id)
				}
				ftick := bids[rng.Intn(len(bids))]
				stats.corrFocuses.Add(1)
				conn.WriteJSON(map[string]string{
					"name": "focus",
					"id":   ftick,
				})
			} else if rng.Intn(5) == 0 {
				if _, ok := state.bought[state.focused]; !ok {
					continue
				}
				if _, ok := state.ungraded[state.focused]; !ok {
					continue
				}
				state.removeUngraded(state.focused)
				decision := "incorrect"
				if rng.Intn(2) == 0 {
					decision = "correct"
				}
				stats.corrGrades.Add(1)
				conn.WriteJSON(map[string]string{
					"name":     "grade",
					"id":       state.focused,
					"decision": decision,
				})
			} else if rng.Intn(20) == 0 {
				if state.focused == "" {
					continue
				}
				msg := ""
				for range rng.Intn(20) + 3 {
					msg += words[rng.Intn(len(words))]
					msg += " "
				}
				msg = strings.TrimSuffix(msg, " ")

				stats.corrWrites.Add(1)
				conn.WriteJSON(map[string]string{
					"name":    "write",
					"id":      state.focused,
					"message": msg,
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

	fmt.Printf("Starting stress test: %d teams (%d players/team = %d players), %d correctors\n",
		len(teamTokens), playersPerTeam, len(teamTokens)*playersPerTeam, len(corrTokens))

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				elapsed := time.Since(start).Round(time.Second)
				fmt.Printf("[%v] active: %d players, %d corrs | actions sent -> buy:%d sell:%d solve:%d write:%d focus:%d | corr -> grade:%d ungraded:%d write:%d focus:%d | events recv -> bought:%d sold:%d solved:%d | errs:%d disconn:%d\n",
					elapsed,
					stats.activePlayers.Load(),
					stats.activeCorrs.Load(),
					stats.playerBuys.Load(),
					stats.playerSells.Load(),
					stats.playerSolves.Load(),
					stats.playerWrites.Load(),
					stats.playerFocuses.Load(),
					stats.corrGrades.Load(),
					stats.ungradedTickets.Load(),
					stats.corrWrites.Load(),
					stats.corrFocuses.Load(),
					stats.probsBought.Load(),
					stats.probsSold.Load(),
					stats.probsSolved.Load(),
					stats.serverErrors.Load(),
					stats.wsDisconnects.Load(),
				)
			}
		}
	}()

	var wg sync.WaitGroup

	for _, ctoken := range corrTokens {
		wg.Add(1)
		go func(tok string) {
			defer wg.Done()
			if err := corrSimLoop(tok); err != nil {
				fmt.Printf("corr %s error: %v\n", tok, err)
			}
		}(ctoken)
		time.Sleep(20 * time.Millisecond)
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
			time.Sleep(200 * time.Millisecond)
		}
	}

	wg.Wait()
	close(done)
}
