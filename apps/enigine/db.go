package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var diffs = []string{"A", "B", "C"}

const (
	StateBefore = "before"
	StateRunning = "running"
	StateAfter = "after"
)

func parseTicketId(tid string) (teamid string, probid string) {
	parts := strings.SplitN(tid, ":", 1)
	teamid = parts[0]
	probid = parts[1]
	return
}

func getState(conn *redis.Client) string {
	res, err := conn.Get(ctx, "state").Result()
	if err != nil { panic(err) }

	return res
}

func setState(conn *redis.Client, state string) {
	err := conn.Set(ctx, "state", state, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getStart(conn *redis.Client) time.Time {
	itime, err := conn.Get(ctx, "start").Int64()
	if err != nil { panic(err) }

	return time.UnixMilli(itime)
}

func setStart(conn *redis.Client, start time.Time) {
	err := conn.Set(ctx, "start", start.UnixMilli(), time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getEnd(conn *redis.Client) time.Time {
	itime, err := conn.Get(ctx, "end").Int64()
	if err != nil { panic(err) }

	return time.UnixMilli(itime)
}

func setEnd(conn *redis.Client, end time.Time) {
	err := conn.Set(ctx, "end", end.UnixMilli(), time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getCorrTickets(conn *redis.Client, corrid string) []string {
	res, err := conn.SMembers(ctx, "corrtickets:" + corrid).Result()
	if err != nil { panic(err) }

	return res
}

func addCorrTicket(conn *redis.Client, corrid, teamid, probid string) {
	err := conn.SAdd(ctx, "corrtickets:" + corrid, teamid + ":" + probid).Err()
	if err != nil { panic(err) }
}

func clearCorrTickets(conn *redis.Client, corrid string) {
	err := conn.Del(ctx, "corrtickets:" + corrid).Err()
	if err != nil { panic(err) }
}

func getMoney(conn *redis.Client, teamid string) int {
	m, err := conn.Get(ctx, "money:" + teamid).Int()
	if err != nil { panic(err) }

	return m
}

func setMoney(conn *redis.Client, teamid string, money int) {
	err := conn.Set(ctx, "money:" + teamid, money, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

type Prob struct {
	Id string `json:"id" redis:"id"`
	Name string `json:"name" redis:"name"`
	Diff string `json:"diff" redis:"diff"`
	Text string `json:"text" redis:"text"`
	Answer string `json:"answer" redis:"answer"`
	Code string `json:"code" redis:"code"`
	Auto bool `json:"auto" redis:"auto"`
	Infinite bool `json:"infinite" redis:"infinite"`
	Queue []string `json:"queue" redis:"queue"`
}

func (p Prob) toMap() map[string]string {
	res := map[string]string{
		"id": p.Id,
		"name": p.Name,
		"diff": p.Diff,
		"text": p.Text,
		"answer": p.Answer,
		"code": p.Code,
		"auto": "false",
		"infinite": "false",
		"queue": strings.Join(p.Queue, ":"),
	}
	if p.Auto { res["auto"] = "true" }
	if p.Infinite { res["infinite"] = "true" }
	return res
}

func (p *Prob) fromMap(m map[string]string) {
	p.Id = m["id"]
	p.Name = m["name"]
	p.Diff = m["diff"]
	p.Text = m["text"]
	p.Answer = m["answer"]
	p.Code = m["code"]
	p.Auto = m["auto"] == "true"
	p.Infinite = m["infinite"] == "true"
	p.Queue = strings.Split(m["queue"], ":")
} 

func getProb(conn *redis.Client, probid string) Prob {
	pmap, err := conn.HGetAll(ctx, "prob:" + probid).Result()
	if err != nil { panic(err) }

	res := Prob{}
	res.fromMap(pmap)

	return res
}

func setProb(conn *redis.Client, prob Prob) {
	err := conn.HSet(ctx, "prob:" + prob.Id, prob.toMap()).Err()
	if err != nil { panic(err) }
}

func getPlayToken(conn *redis.Client, token string) string {
	id, err := conn.Get(ctx, "playtoken:" + token).Result()
	if err == nil { return id }
	if err == redis.Nil { return "" }
	panic(err)
}

func setPlayToken(conn *redis.Client, token string, teamid string) {
	err := conn.Set(ctx, "playtoken:" + token, teamid, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getCorrToken(conn *redis.Client, token string) string {
	id, err := conn.Get(ctx, "corrtoken:" + token).Result()
	if err == nil { return id }
	if err == redis.Nil { return "" }
	panic(err)
}

func setCorrToken(conn *redis.Client, token string, teamid string) {
	err := conn.Set(ctx, "corrtoken:" + token, teamid, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

const (
	PriceBuy = "buy"
	PriceSell = "sell"
	PriceSolve = "solve"
)

func getPrice(conn *redis.Client, ptype string, diff string) int {
	res, err := conn.Get(ctx, "price:" + ptype + ":" + diff).Int()
	if err != nil { panic(err) }

	return res
}

func setPrice(conn *redis.Client, ptype string, diff string, price int) {
	err := conn.Set(ctx, "price:" + ptype + ":" + diff, price, time.Duration(0))
	if err != nil { panic(err) }
}

const (
	OwnedFree = "free"
	OwnedBought = "bought"
	OwnedSold = "sold"
	OwnedSolved = "solved"
)

func getOwnedProbs(conn *redis.Client, teamid, otype, diff string) []string {
	res, err := conn.SMembers(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff).Result()
	if err != nil { panic(err) }

	return res
}

func popOwnedProb(conn *redis.Client, teamid, otype, diff string) string {
	res, err := conn.SPop(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff).Result()
	if err == nil { return res }
	if err == redis.Nil { return "" }
	panic(err)
}

func addOwnedProb(conn *redis.Client, teamid, otype, diff, probid string) {
	_, err := conn.SAdd(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff, probid).Result()
	if err != nil { panic(err) }
}

func moveOwnedProb(conn *redis.Client, teamid, diff, probid, srcotype, dstotype string) {
	err := conn.SMove(ctx, "oprobs:" + teamid + ":" + srcotype + ":" + diff, "oprobs:" + teamid + ":" + dstotype + ":" + diff, probid).Err()
	if err != nil { panic(err) }
}

type TLineAtom struct {
	Mside string `json:"mside"`
	Mtype string `json:"mtype"`
	Msg string `json:"msg"`
	Time time.Time `json:"time"`
}

func (t TLineAtom) String() string {
	return t.Mside + ":" + t.Mtype + ":" + t.Msg + ":" + strconv.Itoa(int(t.Time.UnixMilli()))
}

func (t *TLineAtom) fromString(s string) {
	sparts := strings.SplitN(s, ":", 4)
	t.Mside = sparts[0]
	t.Mtype = sparts[1]
	t.Msg = sparts[2]
	mtime, err := strconv.Atoi(sparts[3])
	if err != nil { panic(err) }
	t.Time = time.UnixMilli(int64(mtime))
}

const (
	MSideAdmin = "admin"
	MSidePlayer = "player"
	MTypeText = "text"
	MTypeGif = "gif"
	MTypeSolve = "solve"
	MTypeGrade = "grade"
	MTypeBought = "bought"
	MTypeSold = "sold"
	MTypeSolved = "solved"
)

func pushTLine(conn *redis.Client, teamid, probid string, tla TLineAtom) {
	err := conn.RPush(ctx, "tline:" + teamid + ":" + probid, tla.String()).Err()
	if err != nil { panic(err) }
}

func readTLine(conn *redis.Client, teamid, probid string) []TLineAtom {
	sres, err := conn.LRange(ctx, "tline:" + teamid + ":" + probid, 0, -1).Result()
	if err != nil { panic(err) }

	res := make([]TLineAtom, len(sres))

	for i, s := range sres {
		var tla TLineAtom
		tla.fromString(s)
		res[i] = tla
	}
	
	return res
}

// Owned
func getTState(conn *redis.Client, teamid, probid string) string {
	res, err := conn.Get(ctx, "tstate:" + teamid + ":" + probid).Result()
	if err != nil { panic(err) }

	return res
}

func setTState(conn *redis.Client, teamid, probid, state string) {
	err := conn.Set(ctx, "tstate:" + teamid + ":" + probid, state, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getTCorr(conn *redis.Client, teamid, probid string) string {
	res, err := conn.Get(ctx, "tcorr:" + teamid + ":" + probid).Result()
	if err != nil { panic(err) }

	return res
}

func setTCorr(conn *redis.Client, teamid, probid, corr string) {
	err := conn.Set(ctx, "tcorr:" + teamid + ":" + probid, corr, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getRank(conn *redis.Client, teamid string) int {
	res, err := conn.Get(ctx, "rank:" + teamid).Int()
	if err != nil { panic(err) }

	return res
}

func setRank(conn *redis.Client, teamid string, rank int) {
	err := conn.Set(ctx, "rank:" + teamid, rank, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getProbDiffValidity(conn *redis.Client, probid string) string {
	res, err := conn.HGet(ctx, "prob:" + probid, "diff").Result()
	if err == nil { return res }
	if err == redis.Nil { return "" }
	panic(err)
}

func getProbAnswerValidity(conn *redis.Client, probid string) string {
	res, err := conn.HGet(ctx, "prob:" + probid, "answer").Result()
	if err == nil { return res }
	if err == redis.Nil { return "" }
	panic(err)
}

func getTeamName(conn *redis.Client, teamid string) string {
	res, err := conn.Get(ctx, "teamname:" + teamid).Result()
	if err != nil { panic(err) }

	return res
}

func setTeamName(conn *redis.Client, teamid string, name string) {
	err := conn.Set(ctx, "teamname:" + teamid, name, time.Duration(0)).Err()
	if err != nil { panic(err) }
}

func getNumberRemProbs(conn *redis.Client, teamid string, diff string) int {
	res, err := conn.SCard(ctx, "oprobs:" + teamid + ":" + OwnedFree + ":" + diff).Result()
	if err != nil { panic(err) }

	return int(res)
}
