package repositories

import (
	"database/sql"
	"gobase-app/models"
	"strconv"
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
	DivisionIDs    []int64
}

type UserUpdateParams struct {
	ID             int
	HashedPassword string
	Name           string
	Email          string
	DivisionIDs    []int64
}

// GetAll mengambil seluruh data user aktif beserta role yang terkait.
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT 
			u.id,
			u.name, 
			u.email, 
			COALESCE(GROUP_CONCAT(DISTINCT d.id ORDER BY d.id SEPARATOR ','), '') AS division_ids_csv,
			COALESCE(GROUP_CONCAT(DISTINCT d.name ORDER BY d.name SEPARATOR ', '), '') AS division_display,
			u.created_at,
			COALESCE(GROUP_CONCAT(DISTINCT r2.name ORDER BY r2.name SEPARATOR ', '), '') AS role_display
		FROM users u
		LEFT JOIN user_divisions ud ON ud.user_id = u.id
		LEFT JOIN divisions d ON d.id = ud.division_id AND d.deleted_at IS NULL
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
			u              models.User
			divisionIDsCSV string
			createdAt      time.Time
		)

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&divisionIDsCSV,
			&u.DivisionDisplay,
			&createdAt,
			&u.RoleDisplay,
		); err != nil {
			return nil, err
		}

		u.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		u.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		u.DivisionIDs = splitCSVToIntSlice(divisionIDsCSV)
		u.DivisionNames = splitAndTrimCSV(u.DivisionDisplay)
		if u.DivisionDisplay == "" {
			u.DivisionDisplay = "-"
		}
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

	if len(params.DivisionIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO user_divisions (user_id, division_id, created_at, updated_at)
			VALUES (?, ?, NOW(), NOW())
		`)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		defer stmt.Close()

		for _, divisionID := range params.DivisionIDs {
			if _, err := stmt.Exec(userID, divisionID); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
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

	if _, err := tx.Exec(`DELETE FROM user_divisions WHERE user_id = ?`, params.ID); err != nil {
		tx.Rollback()
		return err
	}

	if len(params.DivisionIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO user_divisions (user_id, division_id, created_at, updated_at)
			VALUES (?, ?, NOW(), NOW())
		`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		for _, divisionID := range params.DivisionIDs {
			if _, err := stmt.Exec(params.ID, divisionID); err != nil {
				tx.Rollback()
				return err
			}
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

// GetDivisions mengambil daftar divisi aktif untuk kebutuhan dropdown user.
func (r *UserRepository) GetDivisions() ([]models.DivisionOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name
		FROM divisions
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var divisions []models.DivisionOption
	for rows.Next() {
		var division models.DivisionOption
		if err := rows.Scan(&division.ID, &division.Name); err != nil {
			return nil, err
		}
		divisions = append(divisions, division)
	}

	return divisions, rows.Err()
}

// FindExistingDivisionIDs mengembalikan map id divisi yang ditemukan di database.
func (r *UserRepository) FindExistingDivisionIDs(ids []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)
	if len(ids) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
		SELECT id
		FROM divisions
		WHERE id IN (` + strings.Join(placeholders, ",") + `) AND deleted_at IS NULL
	`
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result[id] = true
	}

	return result, rows.Err()
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

func splitCSVToIntSlice(val string) []int {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}

	parts := strings.Split(val, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil || id <= 0 {
			continue
		}
		result = append(result, id)
	}
	return result
}
