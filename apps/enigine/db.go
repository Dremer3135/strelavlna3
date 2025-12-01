package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var DIFFS = []string{"A", "B", "C"}

const (
	StateBefore = "before"
	StateRunning = "running"
	StateAfter = "after"
	StateResults = "results"
	StatePaused = "paused"

	DecisionCorrect = "correct"
	DecisionIncorrect = "incorrect"
)

func parseTicketId(tid string) (teamid string, probid string) {
	parts := strings.SplitN(tid, ":", 2)
	teamid = parts[0]
	probid = parts[1]
	return
}

func getState(conn *redis.Client) (string, error) {
	res, err := conn.Get(ctx, "state").Result()
	if err == redis.Nil {
		return StateBefore, nil
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

func setState(conn *redis.Client, state string) error {
	err := conn.Set(ctx, "state", state, time.Duration(0)).Err()
	return err
}

func getStart(conn *redis.Client) (time.Time, error) {
	itime, err := conn.Get(ctx, "start").Int64()
	if err == redis.Nil {
		return time.Now(), nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(itime), nil
}

func setStart(conn *redis.Client, start time.Time) error {
	return conn.Set(ctx, "start", start.UnixMilli(), time.Duration(0)).Err()
}

func getEnd(conn *redis.Client) (time.Time, error) {
	itime, err := conn.Get(ctx, "end").Int64()
	if err == redis.Nil {
		return time.Now(), nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(itime), nil
}

func setEnd(conn *redis.Client, end time.Time) error {
	return conn.Set(ctx, "end", end.UnixMilli(), time.Duration(0)).Err()
}

func getCorrTickets(conn *redis.Client, corrid string) ([]string, error) {
	return conn.SMembers(ctx, "corrtickets:" + corrid).Result()
}

func addCorrTicket(conn *redis.Client, corrid, teamid, probid string) error {
	return conn.SAdd(ctx, "corrtickets:" + corrid, teamid + ":" + probid).Err()
}

func clearCorrTickets(conn *redis.Client, corrid string) error {
	return conn.Del(ctx, "corrtickets:" + corrid).Err()
}

func getMoney(conn *redis.Client, teamid string) (int, error) {
	m, err := conn.Get(ctx, "money:" + teamid).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return m, nil
}

func setMoney(conn *redis.Client, teamid string, money int) error {
	return conn.Set(ctx, "money:" + teamid, money, time.Duration(0)).Err()
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
	Images []string `json:"images" redis:"images"`
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
		"images": strings.Join(p.Images, "\\"),
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
	p.Images = strings.Split(m["images"], "\\")
	if m["images"] == "" { p.Images = []string{} }
} 

func getProb(conn *redis.Client, probid string) (Prob, error) {
	pmap, err := conn.HGetAll(ctx, "prob:" + probid).Result()
	if err == redis.Nil {
		return Prob{}, nil
	}
	if err != nil {
		return Prob{}, err
	}

	res := Prob{}
	res.fromMap(pmap)

	return res, nil
}

func setProb(conn *redis.Client, prob Prob) error {
	return conn.HSet(ctx, "prob:" + prob.Id, prob.toMap()).Err()
}

func getPlayToken(conn *redis.Client, token string) (string, error) {
	id, err := conn.Get(ctx, "playtoken:" + token).Result()
	if err == redis.Nil {
		return "", nil
	}
	return id, err
}

func setPlayToken(conn *redis.Client, token string, teamid string) error {
	return conn.Set(ctx, "playtoken:" + token, teamid, time.Duration(0)).Err()
}

func getCorrToken(conn *redis.Client, token string) (string, error) {
	id, err := conn.Get(ctx, "corrtoken:" + token).Result()
	if err == redis.Nil {
		return "", nil
	}
	return id, err
}

func setCorrToken(conn *redis.Client, token string, teamid string) error {
	return conn.Set(ctx, "corrtoken:" + token, teamid, time.Duration(0)).Err()
}

const (
	PriceBuy = "buy"
	PriceSell = "sell"
	PriceSolve = "solve"
)

func getPrice(conn *redis.Client, ptype string, diff string) (int, error) {
	res, err := conn.Get(ctx, "price:" + ptype + ":" + diff).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return res, nil
}

func setPrice(conn *redis.Client, ptype string, diff string, price int) error {
	return conn.Set(ctx, "price:" + ptype + ":" + diff, price, time.Duration(0)).Err()
}

const (
	OwnedFree = "free"
	OwnedBought = "bought"
	OwnedSold = "sold"
	OwnedSolved = "solved"
)

func getOwnedProbs(conn *redis.Client, teamid, otype, diff string) ([]string, error) {
	return conn.SMembers(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff).Result()
}

func popOwnedProb(conn *redis.Client, teamid, otype, diff string) (string, error) {
	res, err := conn.SPop(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func addOwnedProb(conn *redis.Client, teamid, otype, diff, probid string) error {
	_, err := conn.SAdd(ctx, "oprobs:" + teamid + ":" + otype + ":" + diff, probid).Result()
	return err
}

func moveOwnedProb(conn *redis.Client, teamid, diff, probid, srcotype, dstotype string) error {
	return conn.SMove(ctx, "oprobs:" + teamid + ":" + srcotype + ":" + diff, "oprobs:" + teamid + ":" + dstotype + ":" + diff, probid).Err()
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

func (t *TLineAtom) fromString(s string) error {
	sparts := strings.SplitN(s, ":", 4)
	t.Mside = sparts[0]
	t.Mtype = sparts[1]
	t.Msg = sparts[2]
	mtime, err := strconv.Atoi(sparts[3])
	if err != nil {
		return err
	}
	t.Time = time.UnixMilli(int64(mtime))
	return nil
}

const (
	MSideAdmin = "admin"
	MSidePlayer = "player"
	MTypeText = "message"
	MTypeGif = "gif"
	MTypeSolve = "answer"
	MTypeGrade = "grade"
	MTypeBought = "bought"
	MTypeSold = "sold"
	MTypeSolved = "solved"
	MTypePaste = "paste"
	MTypeCopy = "copy"
	MTypeFocus = "window-focus"
)

func pushTLine(conn *redis.Client, teamid, probid string, tla TLineAtom) error {
	return conn.RPush(ctx, "tline:" + teamid + ":" + probid, tla.String()).Err()
}

func readTLine(conn *redis.Client, teamid, probid string) ([]TLineAtom, error) {
	sres, err := conn.LRange(ctx, "tline:" + teamid + ":" + probid, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	res := make([]TLineAtom, len(sres))

	for i, s := range sres {
		var tla TLineAtom
		if err := tla.fromString(s); err != nil {
			return nil, err
		}
		res[i] = tla
	}

	return res, nil
}

// Owned
func getTState(conn *redis.Client, teamid, probid string) (string, error) {
	res, err := conn.Get(ctx, "tstate:" + teamid + ":" + probid).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

func setTState(conn *redis.Client, teamid, probid, state string) error {
	return conn.Set(ctx, "tstate:" + teamid + ":" + probid, state, time.Duration(0)).Err()
}

func getTCorr(conn *redis.Client, teamid, probid string) (string, error) {
	res, err := conn.Get(ctx, "tcorr:" + teamid + ":" + probid).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

func setTCorr(conn *redis.Client, teamid, probid, corr string) error {
	return conn.Set(ctx, "tcorr:" + teamid + ":" + probid, corr, time.Duration(0)).Err()
}

func getRank(conn *redis.Client, teamid string) (int, error) {
	res, err := conn.Get(ctx, "rank:" + teamid).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return res, nil
}

func setRank(conn *redis.Client, teamid string, rank int) error {
	return conn.Set(ctx, "rank:" + teamid, rank, time.Duration(0)).Err()
}

func getProbDiffValidity(conn *redis.Client, probid string) (string, error) {
	res, err := conn.HGet(ctx, "prob:" + probid, "diff").Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func getProbAnswerValidity(conn *redis.Client, probid string) (string, error) {
	res, err := conn.HGet(ctx, "prob:" + probid, "answer").Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func getTeamName(conn *redis.Client, teamid string) (string, error) {
	res, err := conn.Get(ctx, "teamname:" + teamid).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

func setTeamName(conn *redis.Client, teamid string, name string) error {
	return conn.Set(ctx, "teamname:" + teamid, name, time.Duration(0)).Err()
}

func getNumberRemProbs(conn *redis.Client, teamid string, diff string) (int, error) {
	res, err := conn.SCard(ctx, "oprobs:" + teamid + ":" + OwnedFree + ":" + diff).Result()
	return int(res), err
}

func setTeams(conn *redis.Client, teams []string) error {
	teamsany := []any{}
	for _, team := range teams {
		teamsany = append(teamsany, team)
	}
	return conn.SAdd(ctx, "teams", teamsany...).Err()
}

func getTeams(conn *redis.Client) ([]string, error) {
	return conn.SMembers(ctx, "teams").Result()
}

func setCorrAdmin(conn *redis.Client, adminid string, admin bool) error {
	return conn.Set(ctx, "corradmin:" + adminid, admin, time.Duration(0)).Err()
}

func getCorrAdmin(conn *redis.Client, adminid string) (bool, error) {
	res, err := conn.Get(ctx, "corradmin:" + adminid).Bool()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return res, nil
}

func addConstant(conn *redis.Client, varn string, value float64) error {
	return conn.HSet(ctx, "constants", varn, value).Err()
}

func getConstants(conn *redis.Client) (map[string]float64, error) {
	resr, err := conn.HGetAll(ctx, "constants").Result()
	if err != nil {
		return nil, err
	}
	res := make(map[string]float64)
	for k, v := range resr {
		rv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		res[k] = rv
	}

	return res, nil
}
