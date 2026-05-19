package repositories

import (
	"database/sql"
	"gobase-app/models"
	"time"
)

type DivisionRepository struct {
	DB *sql.DB
}

func (r *DivisionRepository) GetAll() ([]models.Division, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, COALESCE(prefix_division, '') AS prefix_division, created_at
		FROM divisions
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var divisions []models.Division
	for rows.Next() {
		var (
			division  models.Division
			createdAt time.Time
		)
		if err := rows.Scan(&division.ID, &division.Name, &division.PrefixDivision, &createdAt); err != nil {
			return nil, err
		}
		division.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		divisions = append(divisions, division)
	}

	return divisions, rows.Err()
}

func (r *DivisionRepository) ExistsByName(name string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE name = ? AND deleted_at IS NULL
	`, name).Scan(&count)
	return count > 0, err
}

func (r *DivisionRepository) ExistsByNameExceptID(name string, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE name = ? AND id <> ? AND deleted_at IS NULL
	`, name, excludeID).Scan(&count)
	return count > 0, err
}

func (r *DivisionRepository) ExistsByPrefix(prefix string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE prefix_division = ? AND deleted_at IS NULL
	`, prefix).Scan(&count)
	return count > 0, err
}

func (r *DivisionRepository) ExistsByPrefixExceptID(prefix string, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE prefix_division = ? AND id <> ? AND deleted_at IS NULL
	`, prefix, excludeID).Scan(&count)
	return count > 0, err
}

func (r *DivisionRepository) Create(name, prefix string) error {
	_, err := r.DB.Exec(`
		INSERT INTO divisions (name, prefix_division, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`, name, nullableDivisionPrefix(prefix))
	return err
}

func (r *DivisionRepository) Update(id int, name, prefix string) error {
	_, err := r.DB.Exec(`
		UPDATE divisions
		SET name = ?, prefix_division = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, nullableDivisionPrefix(prefix), id)
	return err
}

func (r *DivisionRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		UPDATE divisions
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`, id)
	return err
}

func nullableDivisionPrefix(prefix string) interface{} {
	if prefix == "" {
		return nil
	}
	return prefix
}
