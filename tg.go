package main

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
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

		event_args := event.Args()

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

		if len(event_args) >= 1 {
			message := fmt.Sprintf("send:\n /send #%s\n YOUR TEXT`", event_args[0])
			return event.Send(message)
		} else {
			return event.Send(models.GetStartText())
		}
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

		bot_username := event.Bot().Me.Username

		message := fmt.Sprintf("https://t.me/%s?start=%s", bot_username, userkey)

		return event.Send(message)

	})

	bot.Handle("/send", func(event tele.Context) error {

		if len(strings.Split(event.Args()[0], "#")) <= 1 {

			return event.Send("wrong syntax. try again")
		}

		userkey := strings.Split(event.Args()[0], "#")[1]

		command := strings.Split(event.Text(), "\n")

		_, command = command[0], command[1:]

		text := strings.Join(command, "\n")
		text = "from : " + event.Sender().FirstName + "\n" + text

		reciver_userid := models.GetUseridWithKey(db, userkey)

		reciver_chat, _ := event.Bot().ChatByID(reciver_userid)
		event.Bot().Send(reciver_chat, text)

		return event.Send("Done")
	})

	bot.Start()
}
