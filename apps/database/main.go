package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"

	"github.com/redis/go-redis/v9"
)

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

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {

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
			oimages := rec.GetStringSlice("images")
			nimages := make([]string, len(oimages))
			for i, img := range oimages {
				nimages[i] = "https://strela-vlna.gchd.cz/api/files/probs/" + rec.Id + "/" + img
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
		).Bind(apis.RequireSuperuserAuth())

		e.Router.POST(
			"/api/rdb",
			func(e *core.RequestEvent) error {

				id := e.Request.URL.Query().Get("id")

				contest, err := e.App.FindRecordById("contests", id)
				if err != nil { return err }

				teams, err := e.App.FindAllRecords("teams", dbx.HashExp{"contest": id})
				if err != nil { return err }

				probs, err := e.App.FindAllRecords("probs", dbx.Like("contests", "%" + id + "%"))
				if err != nil { return err }

				rdb := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
				})

				

				return e.String(200, "ok")
			},
		).Bind(apis.RequireSuperuserAuth())

		return e.Next()
	})

	app.OnRecordUpdate("probs").BindFunc(func(e *core.RecordEvent) error {
		if len(e.Record.GetStringSlice("queue")) > 0 {
			return e.Next()
		}
		correctors := []*core.Record{}

		err := e.App.RecordQuery("correctors").All(&correctors)
		if err != nil { return err }

		author := e.Record.GetString("author")

		corrids := make([]string, len(correctors))
		for i, corr := range correctors {
			corrids[i] = corr.GetString("id")
		}

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

		rand.Shuffle(len(corrids), func(i, j int) {
			corrids[i], corrids[j] = corrids[j], corrids[i]
		})

		e.Record.Set("queue", corrids)

		return e.Next()
	})


	if err := app.Start(); err != nil {
			log.Fatal(err)
	}
}
