package main

import (
	"fmt"
	"log"
	"tg-contact-bot/models"
	"time"

	tele "gopkg.in/telebot.v3"
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

		currentDate := models.GetCurrentDate()
		userid := event.Sender().ID

		if models.CheckExistsUserStartedOnDB(db, userid) == false {
			_, err := useridCourser.Exec(userid, currentDate)
			models.CheckErr(err)

			createUserReult := models.CreateUserData(event.Sender().Username, userid)

			log.Println("New Start Added", " And database data response", createUserReult)
		} else {
			log.Println("old user send start command again", userid)
		}

		return event.Send(models.GetStartText())
	})

	bot.Handle("/getkey", func(event tele.Context) error {
		var (
			userkey string
		)

		userid := event.Sender().ID

		sqlStmt := fmt.Sprintf("SELECT userkey FROM user WHERE userid = %d", userid)

		rows, err := db.Query(sqlStmt, 1)

		models.CheckErr(err)

		for rows.Next() {
			err := rows.Scan(&userkey)
			if err != nil {
				log.Fatal(err)
			}
		}

		return event.Send(userkey)

	})

	bot.Start()
}
