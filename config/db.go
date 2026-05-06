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
	if err := ensureTicketAttachmentsTable(DB); err != nil {
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

func ensureTicketAttachmentsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ticket_attachments (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			ticket_id BIGINT UNSIGNED NOT NULL,
			user_id BIGINT UNSIGNED NULL,
			original_name VARCHAR(255) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_path VARCHAR(500) NOT NULL,
			file_size BIGINT UNSIGNED NOT NULL DEFAULT 0,
			mime_type VARCHAR(255) NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL DEFAULT NULL,
			PRIMARY KEY (id),
			KEY ticket_attachments_ticket_id_index (ticket_id),
			KEY ticket_attachments_user_id_index (user_id)
		)
	`)
	return err
}
