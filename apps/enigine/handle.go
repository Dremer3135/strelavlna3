package main

import (
	"errors"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

type InitLoad struct {
	TeamId string `json:"teamid"`
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
}

func initLoad(conn *redis.Client, teamid string, playerid string) InitLoad {
	state := getState(conn)
	res := InitLoad{}
	res.State = state

	res.TeamId = teamid
	res.PlayerId = playerid
	res.Money = getMoney(conn, teamid)
	res.Start = int(getStart(conn).Sub(time.Now()).Milliseconds())
	res.End = int(getEnd(conn).Sub(time.Now()).Milliseconds())

  if state == StateBefore {
		return res
	}

	tlines := map[string][]TLineAtom{}

	bought := make([]Prob, 0)
	for _, d := range diffs {
		l := getOwnedProbs(conn, teamid, OwnedBought, d)
		for _, i := range l {
			p := getProb(conn, i)
			p.Answer = ""
			p.Code = ""
			bought = append(bought, p)
			tlines[i] = readTLine(conn, teamid, i)
		}
	}
	res.Bought = bought

	solved := make([]Prob, 0)
	for _, d := range diffs {
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
	for _, d := range diffs {
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

func buyProb(conn *redis.Client, teamid string, diff string) (Prob, error) {
	money := getMoney(conn, teamid)
	price := getPrice(conn, PriceBuy, diff)
	if money < price {
		return Prob{}, errors.New("not enough money")
	}
	probid := popOwnedProb(conn, teamid, OwnedFree, diff)
	if probid == "" {
		return Prob{}, errors.New("not available")
	}
	prob := getProb(conn, probid)
	if prob.Infinite {
		addOwnedProb(conn, teamid, OwnedFree, diff, probid)
		// TODO: generate prob
	}
	addOwnedProb(conn, teamid, OwnedBought, diff, probid)
	setMoney(conn, teamid, money - price)
	setTState(conn, teamid, probid, OwnedBought)
	pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	return prob, nil
}

func sellProb(conn *redis.Client, teamid, probid string) error {
	diff := getProbDiffValidity(conn, probid)
	if diff == "" {
		return errors.New("not valid prob")
	}
	bought := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if !slices.Contains(bought, probid) {
		return errors.New("not bought")
	}
	moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSold)
	reward := getPrice(conn, PriceSell, diff)
	money := getMoney(conn, teamid)
	setMoney(conn, teamid, money + reward)
	setTState(conn, teamid, probid, OwnedSold)
	pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeBought, "", time.Now()})
	return nil
}

func solveProb(conn *redis.Client, teamid, probid string) error {
	diff := getProbDiffValidity(conn, probid)
	if diff == "" {
		return errors.New("not valid prob")
	}
	bought := getOwnedProbs(conn, teamid, OwnedBought, diff)
	if !slices.Contains(bought, probid) {
		return errors.New("not bought")
	}
	moveOwnedProb(conn, teamid, diff, probid, OwnedBought, OwnedSolved)
	reward := getPrice(conn, PriceSolve, diff)
	money := getMoney(conn, teamid)
	setMoney(conn, teamid, money + reward)
	setTState(conn, teamid, probid, OwnedSolved)
	pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeSolved, "", time.Now()})
	return nil
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
		pushTLine(conn, teamid, probid, TLineAtom{MSidePlayer, MTypeSolved, "", time.Now()})
		return true, nil
	}
	return false, nil
}
