package models

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"tg-contact-bot/consts"
	"time"
)

func CheckExistsUserStartedOnDB(db *sql.DB, userid int64) bool {

	sqlStmt := `SELECT userid FROM started WHERE userid = ?`
	err := db.QueryRow(sqlStmt, userid).Scan(&userid)

	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return false
	}

	return true
}

func GetCurrentDate() string {

	currentTime := time.Now()
	y, m, d := currentTime.Date()
	dateOnly := fmt.Sprintf("%d-%d-%d", y, m, d)

	return dateOnly
}

func GetStartText() string {

	return consts.StartText
}

func CreateStringKeyWithLength(length int) string {

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := range b {
		b[i] = consts.Charset[seededRand.Intn(len(consts.Charset))]
	}
	return string(b)
}

func CreateUserData(username string, userid int64) bool {
	cfg := LoadConfigs()

	db := LoadDB(cfg.Bot.DbName)

	UserCourser, err := db.Prepare("INSERT INTO user(username, userid, userkey, created) values(?,?,?,?)")
	CheckErr(err)

	userRandomKey := CreateStringKeyWithLength(10)

	for {
		if !CheckExistsUserStartedOnDB(db, userid) {
			userRandomKey = CreateStringKeyWithLength(10)
			continue
		} else {
			break
		}
	}

	_, Eerr := UserCourser.Exec(username, userid, userRandomKey, GetCurrentDate())

	CheckErr(Eerr)

	return true
}
