package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func LoadDB(DbName string) *sql.DB {

	db, err := sql.Open("sqlite3", DbName)
	CheckErr(err)
	return db

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}