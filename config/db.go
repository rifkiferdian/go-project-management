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
	if err := ensureUserColumns(DB, os.Getenv("DB_NAME")); err != nil {
		panic(err)
	}
	if err := ensureApprovalFlowDivisionColumn(DB, os.Getenv("DB_NAME")); err != nil {
		panic(err)
	}
	if err := ensureApplicationPermissions(DB); err != nil {
		panic(err)
	}
	if err := ensureTicketAttachmentsTable(DB); err != nil {
		panic(err)
	}
	if err := ensureTicketTodosTable(DB); err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}

func ensureApplicationPermissions(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT 'Copy ticket template', 'web', NOW(), NOW()
		WHERE NOT EXISTS (
			SELECT 1
			FROM permissions
			WHERE name = 'Copy ticket template'
				AND guard_name = 'web'
		)
	`)
	return err
}

func ensureApprovalFlowDivisionColumn(db *sql.DB, dbName string) error {
	var count int
	if err := db.QueryRow(`
		SELECT COUNT(1)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = ?
			AND TABLE_NAME = 'approval_flows'
			AND COLUMN_NAME = 'division_id'
	`, dbName).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err := db.Exec(`
		ALTER TABLE approval_flows
		ADD COLUMN division_id BIGINT UNSIGNED NULL AFTER id,
		ADD KEY approval_flows_division_id_foreign (division_id),
		ADD CONSTRAINT approval_flows_division_id_foreign
			FOREIGN KEY (division_id) REFERENCES divisions (id)
	`)
	return err
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

func ensureUserColumns(db *sql.DB, dbName string) error {
	columns := []struct {
		name       string
		definition string
	}{
		{name: "employee_id", definition: "ADD COLUMN `employee_id` VARCHAR(100) NULL AFTER `email`"},
	}

	for _, column := range columns {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(1)
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = ?
				AND TABLE_NAME = 'users'
				AND COLUMN_NAME = ?
		`, dbName, column.name).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if _, err := db.Exec("ALTER TABLE users " + column.definition); err != nil {
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

func ensureTicketTodosTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ticket_todos (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			ticket_id BIGINT UNSIGNED NOT NULL,
			content VARCHAR(500) NOT NULL,
			is_done TINYINT(1) NOT NULL DEFAULT 0,
			` + "`order`" + ` INT(11) NOT NULL DEFAULT 1,
			done_at DATETIME NULL,
			created_by BIGINT UNSIGNED NULL,
			updated_by BIGINT UNSIGNED NULL,
			deleted_at TIMESTAMP NULL DEFAULT NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY ticket_todos_ticket_id_foreign (ticket_id),
			KEY ticket_todos_created_by_foreign (created_by),
			KEY ticket_todos_updated_by_foreign (updated_by),
			KEY ticket_todos_ticket_done_order_index (ticket_id, is_done, ` + "`order`" + `),
			CONSTRAINT ticket_todos_ticket_id_foreign FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE,
			CONSTRAINT ticket_todos_created_by_foreign FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
			CONSTRAINT ticket_todos_updated_by_foreign FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
		)
	`)
	return err
}
