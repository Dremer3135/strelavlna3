package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"

	// "github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var printSocketChanMap = map[string]chan PrinterProb{
	"1": make(chan PrinterProb, 100),
	"2": make(chan PrinterProb, 100),
	"3": make(chan PrinterProb, 100),
}

func printingSocketLoop(index string) func(*websocket.Conn) {
	ch := printSocketChanMap[index]
	return func(conn *websocket.Conn) {
		clchan := make(chan bool)
		go func(){
			conn.ReadMessage()
			clchan <- true
		}()
		loop: for {
			select {
			case <-clchan:
				break loop
			case prob, ok := <- ch:
				if !ok {
					conn.WriteJSON(map[string]any{
						"error": "chan closed",
					})
					break loop
				}
				conn.WriteJSON(prob)
			}
		}
	}
}

type PrinterProb struct {
	Diff string
	Name string
	Id string
	TeamName string
	Code string
}

type PaperProb struct {
  Diff string
  Name string
  Index int
  Text string
  Img []string
  Buy int
  Sell int
  Solve int
  TeamName string
  Id string
  AuthorName string
  AuthorSocials string
}

type PaperSol struct {
  Name string
  Index int
  Solution string
}

type PaperConstant struct {
  Id string
  Value float64
  Name string
  Symbol string
  Unit string
  Desc string
  Group string
}

func RequireAuth() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:   "adminauth",
		Func: requireAdminAuth("correctors"),
	}
}

func requireAdminAuth(optCollectionNames ...string) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return e.UnauthorizedError("The request requires valid record authorization token.", nil)
		}

		// check record collection name
		if len(optCollectionNames) > 0 && !slices.Contains(optCollectionNames, e.Auth.Collection().Name) {
			return e.ForbiddenError("The authorized record is not allowed to perform this action.", nil)
		}

		if !e.Auth.GetBool("admin") {
			return e.ForbiddenError("not admin", nil)
		}

		return e.Next()
	}
}

func genProb(app core.App, id string) (string, string, error) {
	rec := core.Record{}
	err := app.RecordQuery("probs").AndWhere(dbx.HashExp{"id": id}).Limit(1).One(&rec)
	if err != nil {
		return "", "", err
	}
	consts := []*core.Record{}
	err = app.RecordQuery("constants").All(&consts)
	if err != nil {
		return "", "", err
	}
	constsMap := map[string]float64{}
	for _, cnst := range consts {
		constsMap[cnst.GetString("variable_name")] = cnst.GetFloat("value")
	}
	jbody, err := json.Marshal(struct{
		Code string `json:"code"`
		Text string `json:"text"`
		Answer string `json:"answer"`
		Consts map[string]float64 `json:"consts"`
		Timeout float32 `json:"timeout"`
		MemMB int `json:"mem_mb"`
	}{
		rec.GetString("code"),
		rec.GetString("text"),
		rec.GetString("answer"),
		constsMap,
		1,
		32,
	})
	if err != nil {
		return "", "", err
	}
	resp, err := http.DefaultClient.Post("http://localhost:8000/run", "application/json", bytes.NewBuffer(jbody))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	prespb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	presp := map[string]any{}
	err = json.Unmarshal(prespb, &presp)
	if err != nil {
		return "", "", err
	}
	succ, ok := presp["success"].(bool)
	if !ok || !succ {
		return "", "", errors.New(fmt.Sprint(presp["error"]))
	}
	text, ok := presp["text"].(string)
	if !ok {
		return "", "", errors.New("idk")
	}
	answer, ok := presp["answer"].(string)
	if !ok {
		return "", "", errors.New("idk")
	}
	return text, answer, nil
}

func latexEscapeComment(s string) string {
  res := s
  res = strings.ReplaceAll(res, `%`, `\%`)
  res = strings.ReplaceAll(res, `<br>`, `\n`)
  res = strings.ReplaceAll(res, `°`, `\degree`)
  return res
}

func latexEscape(s string) string {
  res := s
  res = strings.ReplaceAll(res, `%`, `\%`)
  res = strings.ReplaceAll(res, `{`, `\{`)
  res = strings.ReplaceAll(res, `}`, `\}`)
  res = strings.ReplaceAll(res, `&`, `\&`)
  res = strings.ReplaceAll(res, `^`, `\^`)
  res = strings.ReplaceAll(res, `#`, `\#`)
  res = strings.ReplaceAll(res, `_`, `\_`)
  res = strings.ReplaceAll(res, `$`, `\$`)
  res = strings.ReplaceAll(res, `<br>`, `\n`)
  return res
}

type TeamGameData struct {
	Money int
	Free []string
	Bought []string
	Solved []string
	Sold []string
}

func main() {
	app := pocketbase.New()

	app.Cron().MustAdd(
		"spam",
		"* * * * *",
		func() {
			schools := []*core.Record{}

			err := app.RecordQuery("skoly").
				AndWhere(dbx.HashExp{"spam": false}).
				OrderBy("created").
				Limit(1).
				All(&schools)

			if len(schools) < 1 {
				return
			}

			if err != nil {
				app.Logger().Error("school query failed", "err", err)
				return
			}

			text := core.Record{}

			err = app.RecordQuery("texts").
				AndWhere(dbx.HashExp{"name": "spam_mail"}).
				Limit(1).
				One(&text)

			if err != nil {
				app.Logger().Error("text query failed", "err", err)
				return
			}

			schoolmails := []mail.Address{}

			for _, school := range schools {
				email := school.GetString("email_1")
				if email == "" {
					email = school.GetString("email_2")
				}
				if email == "" {
					return
				}
				addr, err := mail.ParseAddress(email)
				if err != nil { continue }
				schoolmails = append(schoolmails, *addr)
			}

			msg := &mailer.Message{
				From: mail.Address{
					Name: app.Settings().Meta.SenderName,
					Address: app.Settings().Meta.SenderAddress,
				},
				To: schoolmails[0:1],
				Bcc: schoolmails[1:],
				Subject: text.GetString("data"),
				HTML: text.GetString("text"),
			}

			err = app.NewMailClient().Send(msg)

			if err != nil {
				app.Logger().Error("mail send failed", "err", err)
			}

			for _, school := range schools {
				// school.Set("spam", true)
				// err := app.Save(school)
				err := app.DB().Update("skoly", dbx.Params{"spam": true}, dbx.HashExp{"id": school.Id})
				if err != nil {
					app.Logger().Error("app.Save failed", "err", err)
				}
			}

		},
	)

	app.Cron().MustAdd(
		"final_mail",
		"* * * * *",
		func() {
			teams := []*core.Record{}

			err := app.RecordQuery("teams").
				AndWhere(dbx.HashExp{"token": ""}).
				OrderBy("created").
				Limit(1).
				All(&teams)

			if len(teams) < 1 {
				return
			}

			team := teams[0]

			if err != nil {
				app.Logger().Error("team query failed", "err", err)
				return
			}

			text := core.Record{}

			err = app.RecordQuery("texts").
				AndWhere(dbx.HashExp{"name": "reg_confirm"}).
				Limit(1).
				One(&text)

			if err != nil {
				app.Logger().Error("text query failed", "err", err)
				return
			}

			contest, err := app.FindRecordById("contests", team.GetString("contest"))
			if err != nil {
				app.Logger().Error("texttempl query failed", "err", err)
				return
			}

			token := security.RandomString(5)

			var renbuf bytes.Buffer
			tmpl, err := template.New("reg_confirm").Parse(text.GetString("text"))
			if err != nil {
				app.Logger().Error("texttempl query failed", "err", err)
				return
			}
			err = tmpl.Execute(&renbuf, struct{
				Code,
				CompSubject,
				CompName,
				Whatsapp,
				TeamName,
				Email,
				Player1,
				Player2,
				Player3,
				Player4,
				Player5,
				OnlineRound,
				FinalRound,
				RegistrationStart,
				RegistrationEnd string
			}{
				token,
				contest.GetString("subject"),
				contest.GetString("name"),
				contest.GetString("whatsapp"),
				team.GetString("name"),
				team.GetString("email"),
				team.GetString("player1"),
				team.GetString("player2"),
				team.GetString("player3"),
				team.GetString("player4"),
				team.GetString("player5"),
				contest.GetDateTime("onlineStart").Time().Format("2. 1. 2006 15:04:05"),
				contest.GetDateTime("onSiteStart").Time().Format("2. 1. 2006 15:04:05"),
				contest.GetDateTime("registration_start").Time().Format("2. 1. 2006 15:04:05"),
				contest.GetDateTime("registration_end").Time().Format("2. 1. 2006 15:04:05"),
			})
			if err != nil {
				app.Logger().Error("templ failed", "err", err)
				return
			}

			teacher, err := app.FindRecordById("teachers", team.GetString("teacher"))
			if err != nil {
				app.Logger().Error("teacher failed", "err", err)
				return
			}

			msg := renbuf.String()

			err = app.NewMailClient().Send(&mailer.Message{
				From: mail.Address{
					Address: "strela-vlna@gchd.cz",
					Name: "Střela Vlna",
				},
				To: []mail.Address{ {Address: teacher.GetString("email")}, },
				Cc: []mail.Address{
					{Address: team.GetString("player1email")},
					{Address: team.GetString("player2email")},
					{Address: team.GetString("player3email")},
					{Address: team.GetString("player4email")},
					{Address: team.GetString("player5email")},
				},
				Subject: "Potvrzení registrace do soutěže " + contest.GetString("name"),
				HTML: msg,
			})
			if err != nil {
				app.Logger().Error("mail failed", "err", err)
				return
			}

			team.Set("finalEmail", false)
			team.Set("token", token)
			err = app.Save(team)
			if err != nil {
				app.Logger().Error("team save failed", "err", err)
			}
		},
	)

	app.Cron().MustAdd(
		"between_mail",
		"* * * * *",
		func() {
			teams := []*core.Record{}

			err := app.RecordQuery("teams").
				AndWhere(dbx.Not(dbx.HashExp{"rank": 0})).
				AndWhere(dbx.HashExp{"betweenEmail": false}).
				OrderBy("created").
				Limit(1).
				All(&teams)

			if len(teams) < 1 {
				return
			}

			team := teams[0]

			if err != nil {
				app.Logger().Error("team query failed", "err", err)
				return
			}

			text := core.Record{}

			tmpname := "mail_between"
			if team.GetInt("rank") <= 15 {
				tmpname = "mail_advancing"
			}

			err = app.RecordQuery("texts").
				AndWhere(dbx.HashExp{"name": tmpname}).
				Limit(1).
				One(&text)

			if err != nil {
				app.Logger().Error("text query failed", "err", err)
				return
			}

			contest, err := app.FindRecordById("contests", team.GetString("contest"))
			if err != nil {
				app.Logger().Error("texttempl query failed", "err", err)
				return
			}

			var renbuf bytes.Buffer
			tmpl, err := template.New("reg_confirm").Parse(text.GetString("text"))
			if err != nil {
				app.Logger().Error("texttempl query failed", "err", err)
				return
			}
			err = tmpl.Execute(&renbuf, struct{
				CompName,
				TeamName,
				Rank string
			}{
				contest.GetString("name"),
				team.GetString("name"),
				strconv.Itoa(team.GetInt("rank")),
			})
			if err != nil {
				app.Logger().Error("templ failed", "err", err)
				return
			}

			teacher, err := app.FindRecordById("teachers", team.GetString("teacher"))
			if err != nil {
				app.Logger().Error("teacher failed", "err", err)
				return
			}

			msg := renbuf.String()

			err = app.NewMailClient().Send(&mailer.Message{
				From: mail.Address{
					Address: "strela-vlna@gchd.cz",
					Name: "Střela Vlna",
				},
				To: []mail.Address{ {Address: teacher.GetString("email")}, },
				Cc: []mail.Address{
					{Address: team.GetString("player1email")},
					{Address: team.GetString("player2email")},
					{Address: team.GetString("player3email")},
					{Address: team.GetString("player4email")},
					{Address: team.GetString("player5email")},
				},
				Subject: "Děkujeme za účast v " + contest.GetString("name"),
				HTML: msg,
			})
			if err != nil {
				app.Logger().Error("mail failed", "err", err)
				return
			}

			team.Set("betweenEmail", true)
			err = app.Save(team)
			if err != nil {
				app.Logger().Error("team save failed", "err", err)
			}
		},
	)

	app.Cron().MustAdd(
		"corr_token_mail",
		"*/10 * * * *",
		func() {
			corrs := []*core.Record{}

			err := app.RecordQuery("correctors").
				AndWhere(dbx.HashExp{"token": ""}).
				OrderBy("created").
				Limit(1).
				All(&corrs)

			if len(corrs) < 1 {
				return
			}

			corr := corrs[0]

			if err != nil {
				app.Logger().Error("corr query failed", "err", err)
				return
			}

			text := core.Record{}

			err = app.RecordQuery("texts").
				AndWhere(dbx.HashExp{"name": "corr_token_mail"}).
				Limit(1).
				One(&text)

			if err != nil {
				app.Logger().Error("text query failed", "err", err)
				return
			}

			token := security.RandomString(5)

			var renbuf bytes.Buffer
			tmpl, err := template.New("corr_token_mail").Parse(text.GetString("text"))
			if err != nil {
				app.Logger().Error("texttempl query failed", "err", err)
				return
			}
			err = tmpl.Execute(&renbuf, struct{
				Token string 
				Name string
			}{ token, corr.GetString("username") })
			if err != nil {
				app.Logger().Error("templ failed", "err", err)
				return
			}

			msg := renbuf.String()

			err = app.NewMailClient().Send(&mailer.Message{
				From: mail.Address{
					Address: "strela-vlna@gchd.cz",
					Name: "Střela Vlna",
				},
				To: []mail.Address{ {Address: corr.GetString("email")}, },
				Subject: text.GetString("data"),
				HTML: msg,
			})
			if err != nil {
				app.Logger().Error("mail failed", "err", err)
				return
			}

			corr.Set("token", token)
			err = app.Save(corr)
			if err != nil {
				app.Logger().Error("corr save failed", "err", err)
			}
		},
	)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {

		e.Router.GET("/api/setgamedata", func(e *core.RequestEvent) error {
			teamf := e.Request.URL.Query().Get("teamf")
			probf := e.Request.URL.Query().Get("probf")

			probs, err := e.App.FindRecordsByFilter("probs", probf, "created", -1, 0)
			if err != nil { return err }

			teams, err := e.App.FindRecordsByFilter("teams", teamf, "created", -1, 0)
			if err != nil { return err }

			gamedata := TeamGameData{100, []string{}, []string{}, []string{}, []string{}}

			for _, prob := range probs {
				gamedata.Free = append(gamedata.Free, prob.Id)
			}

			for _, team := range teams {
				team.Set("inPersonGameData", gamedata)
				err := e.App.Save(team)
				if err != nil { return err }
			}

			return nil

		}).BindFunc(requireAdminAuth())

		e.Router.GET("/api/buyprob", func(e *core.RequestEvent) error {
			teamid := e.Request.URL.Query().Get("team")
			diff := e.Request.URL.Query().Get("diff")
			printid := e.Request.URL.Query().Get("printer")

			err := e.App.RunInTransaction(func(txApp core.App) error {

				team, err := txApp.FindRecordById("teams", teamid)
				if err != nil { return err }

				gamedata := TeamGameData{}
				err = json.Unmarshal([]byte(team.GetString("inPersonGameData")), &gamedata)
				if err != nil { return err }

				contest, err := txApp.FindRecordById("contests", team.GetString("contest"))
				if err != nil { return err }

				sconfig := contest.GetString("config")
				config := struct{
					Buy map[string]int
					Sell map[string]int
					Solve map[string]int
				}{}
				err = json.Unmarshal([]byte(sconfig), &config)
				if err != nil { return err }

				if gamedata.Money < config.Buy[diff] {
					return errors.New("not enough money")
				}

				var prob *core.Record

				for range len(gamedata.Free) {
					probid := gamedata.Free[rand.Intn(len(gamedata.Free))]
					prob, err = txApp.FindRecordById("probs", probid)
					if err != nil { return err }
					if prob.GetString("diff") == diff { break }
				}

				if prob == nil {
					return errors.New("no probs left")
				}

				text, ans := prob.GetString("text"), prob.GetString("answer")

				if prob.GetBool("auto") {
					text, ans, err = genProb(txApp, prob.Id)
					if err != nil { return err }
				}

				coll, err := txApp.FindCollectionByNameOrId("tickets")
				if err != nil { return err }

				tick := core.NewRecord(coll)

				tick.Set("team", teamid)
				tick.Set("prob", prob.Id)
				tick.Set("code", security.RandomString(7))
				tick.Set("text", text)
				tick.Set("answer", ans)

				err = txApp.Save(tick)
				if err != nil { return err }

				if !prob.GetBool("infinite") {
					gamedata.Free = slices.DeleteFunc(gamedata.Free, func(a string) bool { return a == prob.Id })
				}
				gamedata.Bought = append(gamedata.Bought, prob.Id)

				gamedata.Money -= config.Buy[diff]

				datab, err := json.Marshal(gamedata)
				if err != nil { return err }

				team.Set("inPersonGameData", string(datab))

				err = txApp.Save(team)
				if err != nil { return err }

				printSocketChanMap[printid] <- PrinterProb{
					Diff: prob.GetString("diff"),
					Name: prob.GetString("name"),
					Id: prob.Id,
					TeamName: team.GetString("name"),
					Code: tick.GetString("code"),
				}
				return nil
			})

			if err != nil { return e.String(400, err.Error()) }

			return e.String(200, "ok")
		})

		e.Router.GET("/api/getprob", func(e *core.RequestEvent) error {
			tickid := e.Request.URL.Query().Get("id")

			tick, err := e.App.FindFirstRecordByData("tickets", "code", tickid)
			if err != nil { return err }

			team, err := e.App.FindRecordById("teams", tick.GetString("team"))
			if err != nil { return err }

			gamedata := TeamGameData{}
			err = json.Unmarshal([]byte(team.GetString("inPersonGameData")), &gamedata)
			if err != nil { return err }

			state := "invalid"
			if slices.Contains(gamedata.Bought, tick.GetString("prob")) {
				state = "bought"
			}
			if slices.Contains(gamedata.Solved, tick.GetString("prob")) {
				state = "solved"
			}
			if slices.Contains(gamedata.Sold, tick.GetString("prob")) {
				state = "sold"
			}

			return e.JSON(200, map[string]any{
				"answer": tick.GetString("answer"),
				"money": gamedata.Money,
				"state": state,
			})
		})

		e.Router.GET("/api/gradeprob", func(e *core.RequestEvent) error {
			tickid := e.Request.URL.Query().Get("id")

			err := e.App.RunInTransaction(func(txApp core.App) error {
				tick, err := e.App.FindFirstRecordByData("tickets", "code", tickid)
				if err != nil { return err }

				team, err := e.App.FindRecordById("teams", tick.GetString("team"))
				if err != nil { return err }

				gamedata := TeamGameData{}
				err = json.Unmarshal([]byte(team.GetString("inPersonGameData")), &gamedata)
				if err != nil { return err }

				contest, err := txApp.FindRecordById("contests", team.GetString("contest"))
				if err != nil { return err }

				sconfig := contest.GetString("config")
				config := struct{
					Buy map[string]int
					Sell map[string]int
					Solve map[string]int
				}{}
				err = json.Unmarshal([]byte(sconfig), &config)
				if err != nil { return err }

				gamedata.Bought = slices.DeleteFunc(gamedata.Bought, func(a string) bool { return a == tick.GetString("prob") })
				gamedata.Solved = append(gamedata.Solved, tick.GetString("prob"))

				prob, err := txApp.FindRecordById("probs", tick.GetString("prob"))
				if err != nil { return err }
				gamedata.Money += config.Solve[prob.GetString("diff")]

				datab, err := json.Marshal(gamedata)
				if err != nil { return err }

				team.Set("inPersonGameData", string(datab))

				err = e.App.Save(team)
				if err != nil { return err }

				return nil
			})

			if err != nil { return e.String(400, err.Error()) }

			return e.String(200, "ok")
		})

		e.Router.GET("/api/sellprob", func(e *core.RequestEvent) error {
			tickid := e.Request.URL.Query().Get("id")

			err := e.App.RunInTransaction(func(txApp core.App) error {
				tick, err := e.App.FindFirstRecordByData("tickets", "code", tickid)
				if err != nil { return err }

				team, err := e.App.FindRecordById("teams", tick.GetString("team"))
				if err != nil { return err }

				gamedata := TeamGameData{}
				err = json.Unmarshal([]byte(team.GetString("inPersonGameData")), &gamedata)
				if err != nil { return err }

				contest, err := txApp.FindRecordById("contests", team.GetString("contest"))
				if err != nil { return err }

				sconfig := contest.GetString("config")
				config := struct{
					Buy map[string]int
					Sell map[string]int
					Solve map[string]int
				}{}
				err = json.Unmarshal([]byte(sconfig), &config)
				if err != nil { return err }

				gamedata.Bought = slices.DeleteFunc(gamedata.Bought, func(a string) bool { return a == tick.GetString("prob") })
				gamedata.Sold = append(gamedata.Sold, tick.GetString("prob"))

				prob, err := txApp.FindRecordById("probs", tick.GetString("prob"))
				if err != nil { return err }
				gamedata.Money += config.Sell[prob.GetString("diff")]

				datab, err := json.Marshal(gamedata)
				if err != nil { return err }

				team.Set("inPersonGameData", string(datab))

				err = e.App.Save(team)
				if err != nil { return err }

				return nil
			})

			if err != nil { return e.String(400, err.Error()) }

			return e.String(200, "ok")
		})

		e.Router.GET("/api/paperprob", func(e *core.RequestEvent) error {
			id := e.Request.URL.Query().Get("id")
			prob, err := e.App.FindRecordById("probs", id)
			if err != nil { return err }
			ntext, nans, err := genProb(e.App, id)
			if err != nil { return err }
			prob.Set("text", ntext)
			prob.Set("answer", nans)

			bts, err := os.ReadFile("/home/strelavlna/strelavlna3/apps/database/prob_templ_box.tex")
			if err != nil { return err }

			tmpl, err := template.New("box_probs_papers").Parse(string(bts))
			if err != nil { return err }

			renbuf := bytes.Buffer{}
			err = tmpl.Execute(&renbuf, struct{
				Text string
				Imgs []string
				Diff string
				Name string
				Index int
			}{prob.GetString("text"), prob.GetStringSlice("images"), prob.GetString("diff"), prob.GetString("name"), -1})
			if err != nil { return err }

			papers := renbuf.String()
			papers = html.UnescapeString(papers)
			papers = strings.ReplaceAll(papers, "%", "\\%")
			papers = strings.ReplaceAll(papers, "°", "\\degree")

			for _, img := range prob.GetStringSlice("images") {
				papers = strings.ReplaceAll(papers, " " + img + " ", img)
			}

			args := []string{"/home/strelavlna/strelavlna3/apps/database/texprob.sh", "/home/strelavlna/strelavlna3/apps/database/pb_data/storage/" + prob.BaseFilesPath(), papers}

			fmt.Printf("%#v\n", args)

			cmd := exec.Command(args[0], args[1:]...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			bres, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("cmd.Output: %w: %s", err, stderr.String())
			}

			fname := strings.TrimSpace(string(bres)) // Remove potential newline

			// Extract directory and base filename
			dir := filepath.Dir(fname)
			base := filepath.Base(fname)

			// Serve the file from the specific temporary directory
			return e.FileFS(os.DirFS(dir), base)
		})

		e.Router.POST("/api/code", func(e *core.RequestEvent) error {
			data := struct{
				Id string `json:"id"`
			}{}
			err := e.BindBody(&data)
			if err != nil {
				return e.Error(400, "invalid data", err)
			}
			rec := core.Record{}
			err = e.App.RecordQuery("probs").AndWhere(dbx.HashExp{"id": data.Id}).Limit(1).One(&rec)
			if err != nil {
				return e.Error(400, "invalid prob id", err)
			}
			oimages := rec.GetStringSlice("images")
			nimages := make([]string, len(oimages))
			for i, img := range oimages {
				nimages[i] = "https://strela-vlna.gchd.cz/api/files/probs/" + rec.Id + "/" + img
			}
			if !rec.GetBool("auto") {
				return e.JSON(200, struct{
					Text string `json:"text"`
					Answer string `json:"answer"`
					Images []string `json:"images"`
					Diff string `json:"diff"`
					Name string `json:"name"`
					Id string `json:"id"`
				}{rec.GetString("text"), rec.GetString("answer"), nimages, rec.GetString("diff"), rec.GetString("name"), rec.Id})
			}
			consts := []*core.Record{}
			err = e.App.RecordQuery("constants").All(&consts)
			if err != nil {
				return e.Error(500, "idk", err)
			}
			constsMap := map[string]float64{}
			for _, cnst := range consts {
				constsMap[cnst.GetString("variable_name")] = cnst.GetFloat("value")
			}
			jbody, err := json.Marshal(struct{
				Code string `json:"code"`
				Text string `json:"text"`
				Answer string `json:"answer"`
				Consts map[string]float64 `json:"consts"`
				Timeout float32 `json:"timeout"`
				MemMB int `json:"mem_mb"`
			}{
				rec.GetString("code"),
				rec.GetString("text"),
				rec.GetString("answer"),
				constsMap,
				1,
				32,
			})
			if err != nil {
				return e.Error(500, "atp idk bro", err)
			}
			resp, err := http.DefaultClient.Post("http://localhost:8000/run", "application/json", bytes.NewBuffer(jbody))
			if err != nil {
				return e.Error(500, "atp idk bro", err)
			}
			defer resp.Body.Close()
			prespb, err := io.ReadAll(resp.Body)
			if err != nil {
				return e.Error(500, "atp idk bro", err)
			}
			presp := map[string]any{}
			err = json.Unmarshal(prespb, &presp)
			if err != nil {
				return e.Error(500, "atp idk bro", err)
			}
			succ, ok := presp["success"].(bool)
			if !ok || !succ {
				return e.JSON(400, struct{
					Data string `json:"data"`

				}{Data: fmt.Sprint(presp["error"])})
			}
			text, ok := presp["text"].(string)
			if !ok {
				return e.Error(500, "atp idk bro", err)
			}
			answer, ok := presp["answer"].(string)
			if !ok {
				return e.Error(500, "atp idk bro", err)
			}
			return e.JSON(200, struct{
				Text string `json:"text"`
				Answer string `json:"answer"`
				Images []string `json:"images"`
				Diff string `json:"diff"`
				Name string `json:"name"`
				Id string `json:"id"`
			}{text, answer, nimages, rec.GetString("diff"), rec.GetString("name"), rec.Id})
		}).Bind(apis.RequireAuth("correctors"))

		e.Router.GET("/api/printsocket", func(e *core.RequestEvent) error {
			printid := e.Request.URL.Query().Get("id")
			conn, err := upgrader.Upgrade(e.Response, e.Request, nil)
			if err != nil { return err }
			go printingSocketLoop(printid)(conn)
			return nil
		})

		e.Router.GET("/api/printrprob", func(e *core.RequestEvent) error {
			printid := e.Request.URL.Query().Get("id")
			probs, err := e.App.FindAllRecords("probs")
			if err != nil { return err }
			prob := probs[rand.Intn(len(probs))]
			printSocketChanMap[printid] <- PrinterProb{
				Diff: prob.GetString("diff"),
				Name: prob.GetString("name"),
				Id: prob.Id,
				TeamName: "Bambuláci",
				Code: "676767",
			}
			return e.String(200, "ok")
		})

		// e.Router.POST("/loadprobs", func(e *core.RequestEvent) error {
		//
		// 	body, err := io.ReadAll(e.Request.Body)
		// 	if err != nil { return err }
		// 	e.Request.Body.Close()
		//
		// 	data := []map[string]string{}
		// 	err = json.Unmarshal(body, &data)
		// 	if err != nil { return err }
		//
		// 	coll, _ := e.App.FindCollectionByNameOrId("probs")
		//
		// 	for _, prob := range data {
		// 		rec := core.NewRecord(coll)
		// 		rec.Set("id", prob["id"])
		// 		rec.Set("name", prob["name"])
		// 		rec.Set("diff", prob["diff"])
		// 		rec.Set("type", prob["type"])
		// 		rec.Set("text", prob["text"])
		// 		rec.Set("answer", prob["solution"])
		// 		rec.Set("author", prob["author"])
		// 		rec.Set("infinite", prob["infinite"] == "1")
		// 		err = e.App.Save(rec)
		// 		if err != nil { return err }
		// 	}
		//
		// 	return e.String(200, "ok")
		// })

		e.Router.POST(
			"/api/sql",
			func(e *core.RequestEvent) error {
				data, err := io.ReadAll(e.Request.Body)
				if err != nil { return err }
				body := string(data)

				e.Request.Body.Close()

				rows, err := app.DB().NewQuery(body).Rows()
				if err != nil { return err }

				res := []map[string]string{}

				for rows.Next() {
					row := dbx.NullStringMap{}
					err := rows.ScanMap(row)
					if err != nil { return err }
					rrow := map[string]string{}
					for k, v := range row {
						if !v.Valid { continue }
						rrow[k] = v.String
					}
					res = append(res, rrow)
				}

				return e.JSON(200, res)
			},
		).Bind(RequireAuth())

		e.Router.GET(
			"/api/papers",
			func(e *core.RequestEvent) error {
				id := e.Request.URL.Query().Get("id")
				filter := e.Request.URL.Query().Get("filter")

				contest, err := e.App.FindRecordById("contests", id)
				if err != nil { return err }

				sconfig := contest.GetString("config")
				config := struct{
					Buy map[string]int
					Sell map[string]int
					Solve map[string]int
				}{}
				err = json.Unmarshal([]byte(sconfig), &config)
				if err != nil { return err }

				probs, err := e.App.FindRecordsByFilter("probs", filter, "created", -1, 0)
				if err != nil { return err }

				fmt.Println(len(probs))

				// teams, err := e.App.FindAllRecords("teams", dbx.HashExp{"contest": id})
				// if err != nil { return err }

				consts, err := e.App.FindAllRecords("constants")
				if err != nil { return err }

				gprobs := make([]PaperProb, 0)
				gsols := make([]PaperSol, 0)
				gconsts := make([]PaperConstant, 0)

				imgsurls := make([]string, 0)

				i := 1
				for _, prob := range probs {
					for _, img := range prob.GetStringSlice("images") {
						imgsurls = append(imgsurls, "https://strela-vlna.gchd.cz/api/files/probs/" + prob.Id + "/" + img)
					}
					res := PaperProb{
						Diff: prob.GetString("diff"),
						Name: latexEscape(prob.GetString("name")),
						Index: i,
						Text: latexEscapeComment(prob.GetString("text")),
						Img: prob.GetStringSlice("images"),
						Buy: config.Buy[prob.GetString("diff")],
						Sell: config.Sell[prob.GetString("diff")],
						Solve: config.Solve[prob.GetString("diff")],
						TeamName: "",
						Id: prob.Id,
						AuthorName: "",
						AuthorSocials: "",
					}
					sres := PaperSol{
						Name: prob.GetString("name"),
						Index: i,
						Solution: prob.GetString("answer"),
					}

					if prob.GetBool("auto") {
						text, ans, err := genProb(e.App, prob.Id)
						if err != nil { return err }
						res.Text = text
						sres.Solution = ans
					}

					gprobs = append(gprobs, res)
					gsols = append(gsols, sres)
					i++

				}

				npprobs := make([]PaperProb, 0)
				// for _, tm := range teams {
					for _, pr := range gprobs {
						// pr.TeamName = tm.GetString("name")
						npprobs = append(npprobs, pr)
					}
				// }

				for _, cnst := range consts {
					gconsts = append(gconsts, PaperConstant{
						Id: cnst.Id,
						Value: cnst.GetFloat("value"),
						Name: cnst.GetString("name"),
						Symbol: cnst.GetString("symbol"),
						Unit: cnst.GetString("unit"),
						Desc: cnst.GetString("desc"),
						Group: cnst.GetString("group"),
					})
				}

				funcsmap := template.FuncMap{ "iseven": func(i int) bool { return i % 2 == 0 } }


				bts, err := os.ReadFile("/home/strelavlna/strelavlna3/apps/database/prob_templ.tex")
				if err != nil { return err }

				tmpl, err := template.New("probs_papers").Funcs(funcsmap).Parse(string(bts))
				if err != nil { return err }

				renbuf := bytes.Buffer{}
				err = tmpl.Execute(&renbuf, npprobs)
				if err != nil { return err }

				papers := renbuf.String()
				papers = html.UnescapeString(papers)


				sol_bts, err := os.ReadFile("/home/strelavlna/strelavlna3/apps/database/prob_sol_templ.tex")
				if err != nil { return err }

				sol_tmpl, err := template.New("probs_sol_papers").Funcs(funcsmap).Parse(string(sol_bts))
				if err != nil { return err }

				sol_renbuf := bytes.Buffer{}
				err = sol_tmpl.Execute(&sol_renbuf, gsols)
				if err != nil { return err }

				sol_papers := sol_renbuf.String()
				sol_papers = html.UnescapeString(sol_papers)


				const_bts, err := os.ReadFile("/home/strelavlna/strelavlna3/apps/database/consts_templ.tex")
				if err != nil { return err }

				const_tmpl, err := template.New("const_papers").Funcs(funcsmap).Parse(string(const_bts))
				if err != nil { return err }

				const_renbuf := bytes.Buffer{}
				err = const_tmpl.Execute(&const_renbuf, gconsts)
				if err != nil { return err }

				const_papers := const_renbuf.String()
				const_papers = html.UnescapeString(const_papers)

				return e.String(200, papers + "\n\n\n" + sol_papers + "\n\n\n" + const_papers + "\n\n\n" + strings.Join(imgsurls, " "))
			},
		).Bind(RequireAuth())

		e.Router.GET(
			"/api/rdb",
			func(e *core.RequestEvent) error {

				id := e.Request.URL.Query().Get("id")

				rdb := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
				})

				contest, err := e.App.FindRecordById("contests", id)
				if err != nil { return err }

				sconfig := contest.GetString("config")
				config := struct{
					Buy map[string]int
					Sell map[string]int
					Solve map[string]int
				}{}
				err = json.Unmarshal([]byte(sconfig), &config)
				if err != nil { return err }

				for d, c := range config.Buy {
					setPrice(rdb, PriceBuy, d, c)
				}
				for d, c := range config.Sell {
					setPrice(rdb, PriceSell, d, c)
				}
				for d, c := range config.Solve {
					setPrice(rdb, PriceSolve, d, c)
				}

				teams, err := e.App.FindAllRecords("teams", dbx.HashExp{"contest": id})
				if err != nil { return err }

				teamids := []string{}
				for _, team := range teams {
					teamids = append(teamids, team.Id)
				}

				setTeams(rdb, teamids)

				for _, team := range teams {
					setMoney(rdb, team.Id, 100)
					setPlayToken(rdb, team.GetString("token"), team.Id)
					setTeamName(rdb, team.Id, team.GetString("name"))
				}

				probs, err := e.App.FindAllRecords("probs", dbx.Like("contests", id), dbx.HashExp{"online": true})
				if err != nil { return err }

				for _, prob := range probs {
					imgs := []string{}
					for _, img := range prob.GetStringSlice("images") {
						imgs = append(imgs, "https://strela-vlna.gchd.cz/api/files/probs/" + prob.Id + "/" + img)
					}
					setProb(rdb, Prob{
						Id: prob.Id,
						Name: prob.GetString("name"),
						Diff: prob.GetString("diff"),
						Text: prob.GetString("text"),
						Answer: prob.GetString("answer"),
						Code: prob.GetString("code"),
						Auto: prob.GetBool("auto"),
						Infinite: prob.GetBool("infinite"),
						Queue: prob.GetStringSlice("queue"),
						Images: imgs,
					})
					for _, team := range teams {
						addOwnedProb(rdb, team.Id, OwnedFree, prob.GetString("diff"), prob.Id)
					}
				}

				corectors, err := e.App.FindAllRecords("correctors")
				if err != nil { return err }

				for _, corr := range corectors {
					setCorrToken(rdb, corr.GetString("token"), corr.Id)
					setCorrAdmin(rdb, corr.Id, corr.GetBool("admin"))
				}

				setState(rdb, StateBefore)

				setStart(rdb, contest.GetDateTime("onlineStart").Time())
				setEnd(rdb, contest.GetDateTime("onlineEnd").Time())

				constants, err := e.App.FindAllRecords("constants")
				if err != nil { return err }

				for _, con := range constants {
					addConstant(rdb, con.GetString("variable_name"), con.GetFloat("value"))
				}

				rdb.Save(ctx).Result()

				return e.String(200, "ok")
			},
		).Bind(RequireAuth())

		return e.Next()
	})

	app.OnRecordUpdate("probs").BindFunc(func(e *core.RecordEvent) error {
		// if len(e.Record.GetStringSlice("queue")) > 0 {
		// 	return e.Next()
		// }
		correctors := []*core.Record{}

		err := e.App.RecordQuery("correctors").All(&correctors)
		if err != nil { return err }

		author := e.Record.GetString("author")

		corrids := make([]string, len(correctors))
		for i, corr := range correctors {
			corrids[i] = corr.GetString("id")
		}

		rand.Shuffle(len(corrids), func(i, j int) {
			corrids[i], corrids[j] = corrids[j], corrids[i]
		})

		if slices.Contains(corrids, author) {
			ncorrids := make([]string, 1, len(corrids))
			ncorrids[0] = author
			for _, corr := range corrids {
				if corr != author {
					ncorrids = append(ncorrids, corr)
				}
			}
			corrids = ncorrids
		}

		e.Record.Set("queue", corrids)

		return e.Next()
	})


	if err := app.Start(); err != nil {
			log.Fatal(err)
	}
}
