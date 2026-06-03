package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	return sql.Open("sqlite", "rakkiz.db")
}