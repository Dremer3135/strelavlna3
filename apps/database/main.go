package main

import (
	"encoding/json"
	"io"
	"log"
	"net/mail"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
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
				Limit(1).
				All(&schools)

			if len(schools) < 1 {
				return
			}

			school := schools[0]

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

			email := school.GetString("email_1")
			if email == "" {
				email = school.GetString("email_2")
			}
			if email == "" {
				return
			}

			msg := &mailer.Message{
				From: mail.Address{
					Name: app.Settings().Meta.SenderName,
					Address: app.Settings().Meta.SenderAddress,
				},
				To: []mail.Address{{Address: email}},
				Subject: text.GetString("data"),
				HTML: text.GetString("text"),
			}

			err = app.NewMailClient().Send(msg)

			if err != nil {
				app.Logger().Error("mail send failed", "err", err)
			}
		},
	)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.POST("/loadprobs", func(e *core.RequestEvent) error {

			body, err := io.ReadAll(e.Request.Body)
			if err != nil { return err }
			e.Request.Body.Close()

			data := []map[string]string{}
			err = json.Unmarshal(body, &data)
			if err != nil { return err }

			coll, _ := e.App.FindCollectionByNameOrId("probs")

			for _, prob := range data {
				rec := core.NewRecord(coll)
				rec.Set("id", prob["id"])
				rec.Set("name", prob["name"])
				rec.Set("diff", prob["diff"])
				rec.Set("type", prob["type"])
				rec.Set("text", prob["text"])
				rec.Set("answer", prob["solution"])
				rec.Set("author", prob["author"])
				rec.Set("infinite", prob["infinite"] == "1")
				err = e.App.Save(rec)
				if err != nil { return err }
			}

			return e.String(200, "ok")
		})

		e.Router.POST(
			"/sql",
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
		)

		return e.Next()
	})


	if err := app.Start(); err != nil {
			log.Fatal(err)
	}
}
