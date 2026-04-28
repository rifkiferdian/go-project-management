package repositories

import (
	"database/sql"
	"gobase-app/models"
	"strings"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

const userModelType = "App\\Models\\User"

type UserCreateParams struct {
	HashedPassword string
	Name           string
	Email          string
}

type UserUpdateParams struct {
	ID             int
	HashedPassword string
	Name           string
	Email          string
}

// GetAll mengambil seluruh data user aktif beserta role yang terkait.
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT 
			u.id,
			u.name, 
			u.email, 
			u.created_at,
			COALESCE(GROUP_CONCAT(DISTINCT r2.name ORDER BY r2.name SEPARATOR ', '), '') AS role_display
		FROM users u
		LEFT JOIN model_has_roles mhr ON mhr.model_id = u.id AND mhr.model_type = ?
		LEFT JOIN roles r2 ON r2.id = mhr.role_id
		WHERE u.deleted_at IS NULL
		GROUP BY u.id, u.name, u.email, u.created_at
		ORDER BY u.created_at DESC
	`, userModelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var (
			u         models.User
			createdAt time.Time
		)

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&createdAt,
			&u.RoleDisplay,
		); err != nil {
			return nil, err
		}

		u.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		u.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		if u.RoleDisplay == "" {
			u.RoleDisplay = "-"
		}
		u.RoleNames = splitAndTrimCSV(u.RoleDisplay)
		users = append(users, u)
	}

	return users, rows.Err()
}

// CreateUserWithRoles menyimpan data user baru beserta assignment rolenya dalam satu transaksi.
func (r *UserRepository) CreateUserWithRoles(params UserCreateParams, roleIDs []int64) (int64, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.Exec(`
		INSERT INTO users (name, email, password, type, created_at, updated_at)
		VALUES (?, ?, ?, 'db', NOW(), NOW())
	`, params.Name, params.Email, params.HashedPassword)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(roleIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO model_has_roles (role_id, model_type, model_id)
			VALUES (?, ?, ?)
		`)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		defer stmt.Close()

		for _, roleID := range roleIDs {
			if _, err := stmt.Exec(roleID, userModelType, userID); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return userID, nil
}

// UpdateUserWithRoles memperbarui data user beserta role assignments dalam satu transaksi.
func (r *UserRepository) UpdateUserWithRoles(params UserUpdateParams, roleIDs []int64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if params.HashedPassword != "" {
		if _, err := tx.Exec(`
			UPDATE users
			SET password = ?, name = ?, email = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
		`, params.HashedPassword, params.Name, params.Email, params.ID); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if _, err := tx.Exec(`
			UPDATE users
			SET name = ?, email = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
		`, params.Name, params.Email, params.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.Exec(`DELETE FROM model_has_roles WHERE model_id = ? AND model_type = ?`, params.ID, userModelType); err != nil {
		tx.Rollback()
		return err
	}

	if len(roleIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO model_has_roles (role_id, model_type, model_id)
			VALUES (?, ?, ?)
		`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		for _, roleID := range roleIDs {
			if _, err := stmt.Exec(roleID, userModelType, params.ID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// ExistsByEmail mengecek apakah email sudah digunakan.
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, nil
	}

	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE email = ? AND deleted_at IS NULL`, email).Scan(&count)
	return count > 0, err
}

// ExistsByEmailExceptID mengecek apakah email sudah digunakan user lain.
func (r *UserRepository) ExistsByEmailExceptID(email string, id int) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, nil
	}

	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE email = ? AND id <> ? AND deleted_at IS NULL`, email, id).Scan(&count)
	return count > 0, err
}

// GetRoleIDsByNames mengambil role_id berdasarkan nama role yang diberikan.
func (r *UserRepository) GetRoleIDsByNames(names []string) (map[string]int64, error) {
	result := make(map[string]int64)
	if len(names) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	for i, name := range names {
		placeholders[i] = "?"
		args[i] = name
	}

	query := `
		SELECT id, name
		FROM roles
		WHERE name IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		result[name] = id
	}

	return result, rows.Err()
}

// DeleteUser menandai user sebagai deleted.
func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.DB.Exec(`UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
}

func splitAndTrimCSV(val string) []string {
	val = strings.TrimSpace(val)
	if val == "" || val == "-" {
		return nil
	}

	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
