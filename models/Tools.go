package models

import (
	"database/sql"

)

func CheckExistsOnDB(rows *sql.Rows) bool{
	var exists bool = false

	for rows.Next() {
		exists = true
	}

	return exists
}