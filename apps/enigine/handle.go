package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

type InitLoad struct {
	TeamId string `json:"teamid"`
	TeamName string `json:"teamname"`
	PlayerId string `json:"playerid"`
	State string `json:"state"`
	Bought []Prob `json:"bought"`
	Sold []Prob `json:"sold"`
	Solved []Prob `json:"solved"`
	TLines map[string][]TLineAtom `json:"tlines"`
	Money int `json:"money"`
	Start int `json:"start"`
	End int `json:"end"`
	Rank int `json:"rank"`
	RemProbs map[string]int `json:"remprobs"`
	BuyCost map[string]int `json:"buycost"`
	SellCost map[string]int `json:"sellcost"`
	SolveCost map[string]int `json:"solvecost"`
}

func initLoad(conn *redis.Client, teamid string, playerid string) InitLoad {
	state := getState(conn)
	res := InitLoad{}
	res.State = state

	res.TeamId = teamid
	res.TeamName = getTeamName(conn, teamid)
	res.PlayerId = playerid
	res.Money = getMoney(conn, teamid)
	res.Start = int(getStart(conn).Sub(time.Now()).Milliseconds())
	res.End = int(getEnd(conn).Sub(time.Now()).Milliseconds())

	res.BuyCost = make(map[string]int)
	res.SellCost = make(map[string]int)
	res.SolveCost = make(map[string]int)

	res.RemProbs = map[string]int{}
	for _, diff := range DIFFS {
		res.RemProbs[diff] = getNumberRemProbs(conn, teamid, diff)
	}

	for _, diff := range DIFFS {
		res.BuyCost[diff] = getPrice(conn, PriceBuy, diff)
		res.SellCost[diff] = getPrice(conn, PriceSell, diff)
		res.SolveCost[diff] = getPrice(conn, PriceSolve, diff)
	}

	res.Bought = make([]Prob, 0)
	res.Sold = make([]Prob, 0)
	res.Solved = make([]Prob, 0)

	res.TLines = make(map[string][]TLineAtom)

  if state == StateBefore {
		return res
	}

	tlines := map[string][]TLineAtom{}

	bought := make([]Prob, 0)
	for _, d := range DIFFS {
		l := getOwnedProbs(conn, teamid, OwnedBought, d)
		for _, i := range l {
			p := getProb(conn, i)
			p.Answer = "<dobrej pokus>"
			p.Code = "<dobrej pokus>"
			bought = append(bought, p)
			tlines[i] = readTLine(conn, teamid, i)
		}
	}
	res.Bought = bought

	solved := make([]Prob, 0)
	for _, d := range DIFFS {
		l := getOwnedProbs(conn, teamid, OwnedSolved, d)
		for _, i := range l {
			p := getProb(conn, i)
			p.Answer = ""
			p.Code = ""
			solved = append(solved, p)
			tlines[i] = readTLine(conn, teamid, i)
		}
	}
	res.Solved = solved

	sold := make([]Prob, 0)
	for _, d := range DIFFS {
		l := getOwnedProbs(conn, teamid, OwnedSold, d)
		for _, i := range l {
			p := getProb(conn, i)
			p.Answer = ""
			p.Code = ""
			sold = append(sold, p)
			tlines[i] = readTLine(conn, teamid, i)
		}
	}
	res.Sold = sold

	res.TLines = tlines

	if state == StateAfter {
		res.Rank = getRank(conn, teamid)
	}

	return res
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func buyProb(conn *redis.Client, teamid string, diff string) (prob Prob, money int, remprobs map[string]int, err error) {
	money = getMoney(conn, teamid)
	price := getPrice(conn, PriceBuy, diff)
	if money < price {
		err = errors.New("not enough money")
		return
	}
	probid := popOwnedProb(conn, teamid, OwnedFree, diff)
	if probid == "" {
		err = errors.New("not available")
		return
	}

	prob = getProb(conn, probid)
	if prob.Auto {
		genProb(conn, &prob)
		if prob.Infinite {
			addOwnedProb(conn, teamid, OwnedFree, diff, probid)
		}
	}
	prob.Id = prob.Id + "/" + randSeq(7)
	setProb(conn, prob)
	addOwnedProb(conn, teamid, OwnedBought, diff, probid)
	money -= price
	setMoney(conn, teamid, money)
	setTState(conn, teamid, probid, OwnedBought)
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	remprobs = make(map[string]int)
	for _, diff := range DIFFS {
		remprobs[diff] = getNumberRemProbs(conn, teamid, diff)
	}
	return
}

func sellProb(conn *redis.Client, teamid, probid string) (money int, err error) {
	diff := getProbDiffValidity(conn, probid)
	if diff == "" {
		return 0, errors.New("not valid prob")
	}
	bought := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if !slices.Contains(bought, probid) {
		return 0, errors.New("not bought")
	}
	moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSold)
	reward := getPrice(conn, PriceSell, diff)
	money = getMoney(conn, teamid)
	money += reward
	setMoney(conn, teamid, money)
	setTState(conn, teamid, probid, OwnedSold)
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	return
}

func solveProb(conn *redis.Client, teamid, probid string) (money int, err error) {
	diff := getProbDiffValidity(conn, probid)
	if diff == "" {
		return 0, errors.New("not valid prob")
	}
	bought := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if !slices.Contains(bought, probid) {
		return 0, errors.New("not bought")
	}
	moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSolved)
	reward := getPrice(conn, PriceSolve, diff)
	money = getMoney(conn, teamid)
	money += reward
	setMoney(conn, teamid, money)
	setTState(conn, teamid, probid, OwnedSolved)
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeSolved, "", time.Now()})
	return
}

func autoGrade(conn *redis.Client, teamid, probid, answer string) (bool, error) {
	ans := getProbAnswerValidity(conn, probid)
	if ans == "" {
		return false, errors.New("prob not valid")
	}
	diff := getProbDiffValidity(conn, probid)
	bought := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if !slices.Contains(bought, probid) {
		return false, errors.New("prob not bought")
	}
	if ans == answer {
		moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSolved)
		reward := getPrice(conn, PriceSolve, diff)
		money := getMoney(conn, teamid)
		setMoney(conn, teamid, money + reward)
		setTState(conn, teamid, probid, OwnedSolved)
		// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeSolved, "", time.Now()})
		return true, nil
	}
	return false, nil
}

type CorrInitLoad struct {
	BoughtTickets map[string]CorrTicket `json:"bought_tickets"`
	SolvedTickets map[string]CorrTicket `json:"solved_tickets"`
	SoldTickets map[string]CorrTicket `json:"sold_tickets"`
	TLines map[string][]TLineAtom `json:"tlines"`
	Start int `json:"start"`
	End int `json:"end"`
	State string `json:"state"`
	Admin bool `json:"admin"`
	Id string `json:"id"`
}

type CorrTicket struct {
	TeamId string `json:"team_id"`
	TeamName string `json:"team_name"`
	Prob Prob `json:"prob"`
}

func corrInitLoad(conn *redis.Client, id string) CorrInitLoad {
	tickets := getCorrTickets(conn, id)
	res := CorrInitLoad{
		BoughtTickets: make(map[string]CorrTicket),
		SolvedTickets: make(map[string]CorrTicket),
		SoldTickets: make(map[string]CorrTicket),
	}

	res.State = getState(conn)
	res.Start = int(getStart(conn).Sub(time.Now()).Milliseconds())
	res.End = int(getEnd(conn).Sub(time.Now()).Milliseconds())
	res.Id = id

	res.Admin = getCorrAdmin(conn, id)

	res.TLines = make(map[string][]TLineAtom)
	
	for _, tickid := range tickets {
		teamid, probid := parseTicketId(tickid)
		tstate := getTState(conn, teamid, probid)
		res.TLines[tickid] = readTLine(conn, teamid, probid)
		if tstate == OwnedBought {
			res.BoughtTickets[tickid] = CorrTicket{
				TeamId: teamid,
				TeamName: getTeamName(conn, teamid),
				Prob: getProb(conn, probid),
			}
		}
		if tstate == OwnedSolved {
			res.SolvedTickets[tickid] = CorrTicket{
				TeamId: teamid,
				TeamName: getTeamName(conn, teamid),
				Prob: getProb(conn, probid),
			}
		}
		if tstate == OwnedSold {
			res.SoldTickets[tickid] = CorrTicket{
				TeamId: teamid,
				TeamName: getTeamName(conn, teamid),
				Prob: getProb(conn, probid),
			}
		}
	}

	return res
}

func genProb(conn *redis.Client, prob *Prob) error {
	constsMap := getConstants(conn)
	jbody, err := json.Marshal(struct{
		Code string `json:"code"`
		Text string `json:"text"`
		Answer string `json:"answer"`
		Consts map[string]float64 `json:"consts"`
		Timeout float32 `json:"timeout"`
		MemMB int `json:"mem_mb"`
	}{
		prob.Code,
		prob.Text,
		prob.Answer,
		constsMap,
		1,
		32,
	})
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Post("http://localhost:8000/run", "application/json", bytes.NewBuffer(jbody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	prespb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	presp := map[string]any{}
	err = json.Unmarshal(prespb, &presp)
	if err != nil {
		return err
	}
	succ, ok := presp["success"].(bool)
	if !ok || !succ {
		return errors.New(fmt.Sprint(presp["error"]))
	}
	text, ok := presp["text"].(string)
	if !ok {
		return errors.New("idk")
	}
	answer, ok := presp["answer"].(string)
	if !ok {
		return errors.New("idk")
	}
	prob.Text = text
	prob.Answer = answer
	return nil
}
