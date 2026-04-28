package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {

	var err error

	// enable parseTime so DATETIME/TIMESTAMP scan into time.Time instead of []byte
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	if err := ensureTicketScheduleColumns(DB, os.Getenv("DB_NAME")); err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}

func ensureTicketScheduleColumns(db *sql.DB, dbName string) error {
	columns := []struct {
		name       string
		definition string
	}{
		{name: "starts_at", definition: "ADD COLUMN `starts_at` DATE NULL AFTER `estimation`"},
		{name: "ends_at", definition: "ADD COLUMN `ends_at` DATE NULL AFTER `starts_at`"},
	}

	for _, column := range columns {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(1)
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = ?
				AND TABLE_NAME = 'tickets'
				AND COLUMN_NAME = ?
		`, dbName, column.name).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if _, err := db.Exec("ALTER TABLE tickets " + column.definition); err != nil {
			return err
		}
	}

	return nil
}
