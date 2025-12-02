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

func initLoad(conn *redis.Client, teamid string, playerid string) (InitLoad, error) {
	res := InitLoad{}
	var err error
	res.State, err = getState(conn)
	if err != nil {
		return res, err
	}

	res.TeamId = teamid
	res.TeamName, err = getTeamName(conn, teamid)
	if err != nil {
		return res, err
	}
	res.PlayerId = playerid
	res.Money, err = getMoney(conn, teamid)
	if err != nil {
		return res, err
	}
	start, err := getStart(conn)
	if err != nil {
		return res, err
	}
	res.Start = int(start.Sub(time.Now()).Milliseconds())
	end, err := getEnd(conn)
	if err != nil {
		return res, err
	}
	res.End = int(end.Sub(time.Now()).Milliseconds())

	res.BuyCost = make(map[string]int)
	res.SellCost = make(map[string]int)
	res.SolveCost = make(map[string]int)

	res.RemProbs = map[string]int{}
	for _, diff := range DIFFS {
		rem, err := getNumberRemProbs(conn, teamid, diff)
		if err != nil {
			return res, err
		}
		res.RemProbs[diff] = rem
	}

	for _, diff := range DIFFS {
		buyCost, err := getPrice(conn, PriceBuy, diff)
		if err != nil {
			return res, err
		}
		res.BuyCost[diff] = buyCost
		sellCost, err := getPrice(conn, PriceSell, diff)
		if err != nil {
			return res, err
		}
		res.SellCost[diff] = sellCost
		solveCost, err := getPrice(conn, PriceSolve, diff)
		if err != nil {
			return res, err
		}
		res.SolveCost[diff] = solveCost
	}

	res.Bought = make([]Prob, 0)
	res.Sold = make([]Prob, 0)
	res.Solved = make([]Prob, 0)

	res.TLines = make(map[string][]TLineAtom)

  if res.State == StateBefore {
		return res, nil
	}

	tlines := map[string][]TLineAtom{}

	bought := make([]Prob, 0)
	for _, d := range DIFFS {
		l, err := getOwnedProbs(conn, teamid, OwnedBought, d)
		if err != nil {
			return res, err
		}
		for _, i := range l {
			p, err := getProb(conn, i)
			if err != nil {
				return res, err
			}
p.Answer = "<dobrej pokus>"
			p.Code = "<dobrej pokus>"
			bought = append(bought, p)
			tl, err := readTLine(conn, teamid, i)
			if err != nil {
				return res, err
			}
			tlines[i] = tl
		}
	}
	res.Bought = bought

	solved := make([]Prob, 0)
	for _, d := range DIFFS {
		l, err := getOwnedProbs(conn, teamid, OwnedSolved, d)
		if err != nil {
			return res, err
		}
		for _, i := range l {
			p, err := getProb(conn, i)
			if err != nil {
				return res, err
			}
			p.Answer = ""
			p.Code = ""
			solved = append(solved, p)
			tl, err := readTLine(conn, teamid, i)
			if err != nil {
				return res, err
			}
			tlines[i] = tl
		}
	}
	res.Solved = solved

	sold := make([]Prob, 0)
	for _, d := range DIFFS {
		l, err := getOwnedProbs(conn, teamid, OwnedSold, d)
		if err != nil {
			return res, err
		}
		for _, i := range l {
			p, err := getProb(conn, i)
			if err != nil {
				return res, err
			}
			p.Answer = ""
			p.Code = ""
			sold = append(sold, p)
			tl, err := readTLine(conn, teamid, i)
			if err != nil {
				return res, err
			}
			tlines[i] = tl
		}
	}
	res.Sold = sold

	res.TLines = tlines

	if res.State == StateAfter {
		rank, err := getRank(conn, teamid)
		if err != nil {
			return res, err
		}
		res.Rank = rank
	}

	return res, nil
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
	money, err = getMoney(conn, teamid)
	if err != nil {
		return
	}
	price, err := getPrice(conn, PriceBuy, diff)
	if err != nil {
		return
	}
	if money < price {
		err = errors.New("not enough money")
		return
	}
	probid, err := popOwnedProb(conn, teamid, OwnedFree, diff)
	if err != nil {
		return
	}
	if probid == "" {
		err = errors.New("not available")
		return
	}

	prob, err = getProb(conn, probid)
	if err != nil {
		return
	}
	if prob.Auto {
		err = genProb(conn, &prob)
		if err != nil {
			return
		}
		if prob.Infinite {
			err = addOwnedProb(conn, teamid, OwnedFree, diff, probid)
			if err != nil {
				return
			}
		}
		prob.Id = prob.Id + "/" + randSeq(7)
		err = setProb(conn, prob)
		if err != nil {
			return
		}
	}
	err = addOwnedProb(conn, teamid, OwnedBought, diff, prob.Id)
	if err != nil {
		return
	}
	money -= price
	err = setMoney(conn, teamid, money)
	if err != nil {
		return
	}
	err = setTState(conn, teamid, prob.Id, OwnedBought)
	if err != nil {
		return
	}
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	remprobs = make(map[string]int)
	for _, d := range DIFFS {
		var rem int
		rem, err = getNumberRemProbs(conn, teamid, d)
		if err != nil {
			return
		}
		remprobs[d] = rem
	}
	return
}

func sellProb(conn *redis.Client, teamid, probid string) (money int, err error) {
	diff, err := getProbDiffValidity(conn, probid)
	if err != nil {
		return 0, err
	}
	if diff == "" {
		return 0, errors.New("not valid prob")
	}
	bought, err := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if err != nil {
		return 0, err
	}
	if !slices.Contains(bought, probid) {
		return 0, errors.New("not bought")
	}
	err = moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSold)
	if err != nil {
		return 0, err
	}
	reward, err := getPrice(conn, PriceSell, diff)
	if err != nil {
		return 0, err
	}
	money, err = getMoney(conn, teamid)
	if err != nil {
		return 0, err
	}
	money += reward
	err = setMoney(conn, teamid, money)
	if err != nil {
		return 0, err
	}
	err = setTState(conn, teamid, probid, OwnedSold)
	if err != nil {
		return 0, err
	}
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	return
}

func solveProb(conn *redis.Client, teamid, probid string) (money int, err error) {
	diff, err := getProbDiffValidity(conn, probid)
	if err != nil {
		return 0, err
	}
	if diff == "" {
		return 0, errors.New("not valid prob")
	}
	bought, err := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if err != nil {
		return 0, err
	}
	if !slices.Contains(bought, probid) {
		return 0, errors.New("not bought")
	}
	err = moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSold)
	if err != nil {
		return 0, err
	}
	reward, err := getPrice(conn, PriceSolve, diff)
	if err != nil {
		return 0, err
	}
	money, err = getMoney(conn, teamid)
	if err != nil {
		return 0, err
	}
	money += reward
	err = setMoney(conn, teamid, money)
	if err != nil {
		return 0, err
	}
	err = setTState(conn, teamid, probid, OwnedSolved)
	if err != nil {
		return 0, err
	}
	// pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeSolved, "", time.Now()})
	return
}

func autoGrade(conn *redis.Client, teamid, probid, answer string) (bool, error) {
	ans, err := getProbAnswerValidity(conn, probid)
	if err != nil {
		return false, err
	}
	if ans == "" {
		return false, errors.New("prob not valid")
	}
	diff, err := getProbDiffValidity(conn, probid)
	if err != nil {
		return false, err
	}
	bought, err := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if err != nil {
		return false, err
	}
	if !slices.Contains(bought, probid) {
		return false, errors.New("prob not bought")
	}
	if ans == answer {
		err = moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSolved)
		if err != nil {
			return false, err
		}
		reward, err := getPrice(conn, PriceSolve, diff)
		if err != nil {
			return false, err
		}
		money, err := getMoney(conn, teamid)
		if err != nil {
			return false, err
		}
		err = setMoney(conn, teamid, money+reward)
		if err != nil {
			return false, err
		}
		err = setTState(conn, teamid, probid, OwnedSolved)
		if err != nil {
			return false, err
		}
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
	TLines map[string][]TLineAtom `json:"tlines"`
}

func corrInitLoad(conn *redis.Client, id string) (CorrInitLoad, error) {
	res := CorrInitLoad{
		BoughtTickets: make(map[string]CorrTicket),
		SolvedTickets: make(map[string]CorrTicket),
		SoldTickets: make(map[string]CorrTicket),
	}
	var err error

	tickets, err := getCorrTickets(conn, id)
	if err != nil {
		return res, err
	}

	res.State, err = getState(conn)
	if err != nil {
		return res, err
	}
	start, err := getStart(conn)
	if err != nil {
		return res, err
	}
	res.Start = int(start.Sub(time.Now()).Milliseconds())
	end, err := getEnd(conn)
	if err != nil {
		return res, err
	}
	res.End = int(end.Sub(time.Now()).Milliseconds())
	res.Id = id

	res.Admin, err = getCorrAdmin(conn, id)
	if err != nil {
		return res, err
	}

	res.TLines = make(map[string][]TLineAtom)
	
	for _, tickid := range tickets {
		teamid, probid := parseTicketId(tickid)
		tstate, err := getTState(conn, teamid, probid)
		if err != nil {
			return res, err
		}
		res.TLines[tickid], err = readTLine(conn, teamid, probid)
		if err != nil {
			return res, err
		}
		if tstate == OwnedBought {
			teamName, err := getTeamName(conn, teamid)
			if err != nil {
				return res, err
			}
			prob, err := getProb(conn, probid)
			if err != nil {
				return res, err
			}
			res.BoughtTickets[tickid] = CorrTicket{
				TeamId:   teamid,
				TeamName: teamName,
				Prob:     prob,
			}
		}
		if tstate == OwnedSolved {
			teamName, err := getTeamName(conn, teamid)
			if err != nil {
				return res, err
			}
			prob, err := getProb(conn, probid)
			if err != nil {
				return res, err
			}
			res.SolvedTickets[tickid] = CorrTicket{
				TeamId:   teamid,
				TeamName: teamName,
				Prob:     prob,
			}
		}
		if tstate == OwnedSold {
			teamName, err := getTeamName(conn, teamid)
			if err != nil {
				return res, err
			}
			prob, err := getProb(conn, probid)
			if err != nil {
				return res, err
			}
			res.SoldTickets[tickid] = CorrTicket{
				TeamId:   teamid,
				TeamName: teamName,
				Prob:     prob,
			}
		}
	}

	return res, nil
}

func genProb(conn *redis.Client, prob *Prob) error {
	constsMap, err := getConstants(conn)
	if err != nil { return err }
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
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post("http://localhost:8000/run", "application/json", bytes.NewBuffer(jbody))
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
