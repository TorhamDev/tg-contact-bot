package main

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"tg-contact-bot/models"
	"time"
)

func main() {

	cfg := models.LoadConfigs()

	db := models.LoadDB(cfg.Bot.DbName)

	useridCourser, err := db.Prepare("INSERT INTO started(userid, created) values(?,?)")
	models.CheckErr(err)

	pref := tele.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tele.LongPoller{Timeout: time.Duration(cfg.Bot.PollerTime) * time.Second},
	}

	bot, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(event tele.Context) error {

		query := fmt.Sprintf("SELECT userid=%d FROM started", event.Sender().ID)
		rows, err := db.Query(query)
		models.CheckErr(err)

		if models.CheckExistsOnDB(rows) == false {
			_, err := useridCourser.Exec(event.Sender().ID, "2020-12-09")
			models.CheckErr(err)
			fmt.Println("New Start Added")
		} else {
			fmt.Println("old user send start command again")
		}

		return event.Send("Hello!")
	})

	bot.Start()
}
