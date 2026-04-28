package repositories

import (
	"database/sql"
	"gobase-app/models"
	"time"
)

type ReferentialRepository struct {
	DB *sql.DB
}

func (r *ReferentialRepository) GetActivities() ([]models.Activity, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, description, created_at
		FROM activities
		WHERE deleted_at IS NULL
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Activity
	for rows.Next() {
		var (
			item      models.Activity
			createdAt time.Time
		)
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ReferentialRepository) CreateActivity(name, description string) error {
	_, err := r.DB.Exec(`
		INSERT INTO activities (name, description, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`, name, description)
	return err
}

func (r *ReferentialRepository) UpdateActivity(id int, name, description string) error {
	_, err := r.DB.Exec(`
		UPDATE activities
		SET name = ?, description = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, description, id)
	return err
}

func (r *ReferentialRepository) DeleteActivity(id int) error {
	_, err := r.DB.Exec(`UPDATE activities SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
}

func (r *ReferentialRepository) GetProjectStatuses() ([]models.StatusReference, error) {
	return r.getStatusReferences("project_statuses")
}

func (r *ReferentialRepository) CreateProjectStatus(name, color string, isDefault bool) error {
	return r.createStatusReference("project_statuses", name, color, isDefault)
}

func (r *ReferentialRepository) UpdateProjectStatus(id int, name, color string, isDefault bool) error {
	return r.updateStatusReference("project_statuses", id, name, color, isDefault)
}

func (r *ReferentialRepository) DeleteProjectStatus(id int) error {
	return r.deleteStatusReference("project_statuses", id)
}

func (r *ReferentialRepository) GetTicketPriorities() ([]models.StatusReference, error) {
	return r.getStatusReferences("ticket_priorities")
}

func (r *ReferentialRepository) CreateTicketPriority(name, color string, isDefault bool) error {
	return r.createStatusReference("ticket_priorities", name, color, isDefault)
}

func (r *ReferentialRepository) UpdateTicketPriority(id int, name, color string, isDefault bool) error {
	return r.updateStatusReference("ticket_priorities", id, name, color, isDefault)
}

func (r *ReferentialRepository) DeleteTicketPriority(id int) error {
	return r.deleteStatusReference("ticket_priorities", id)
}

func (r *ReferentialRepository) GetTicketStatuses() ([]models.TicketStatusReference, error) {
	rows, err := r.DB.Query(`
		SELECT
			ts.id,
			ts.name,
			ts.color,
			ts.is_default,
			ts.order,
			ts.project_id,
			COALESCE(p.name, '') AS project_name,
			ts.created_at
		FROM ticket_statuses ts
		LEFT JOIN projects p ON p.id = ts.project_id
		WHERE ts.deleted_at IS NULL
		ORDER BY ts.project_id IS NOT NULL, ts.order, ts.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketStatusReference
	for rows.Next() {
		var (
			item      models.TicketStatusReference
			projectID sql.NullInt64
			project   sql.NullString
			createdAt time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Color,
			&item.IsDefault,
			&item.Order,
			&projectID,
			&project,
			&createdAt,
		); err != nil {
			return nil, err
		}
		if projectID.Valid {
			item.ProjectID = int(projectID.Int64)
		}
		if project.Valid && project.String != "" {
			item.ProjectName = project.String
		} else {
			item.ProjectName = "Global"
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ReferentialRepository) CreateTicketStatus(name, color string, isDefault bool, order int, projectID *int) error {
	_, err := r.DB.Exec(`
		INSERT INTO ticket_statuses (name, color, is_default, `+"`order`"+`, project_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`, name, color, isDefault, order, nullableProjectID(projectID))
	return err
}

func (r *ReferentialRepository) UpdateTicketStatus(id int, name, color string, isDefault bool, order int, projectID *int) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_statuses
		SET name = ?, color = ?, is_default = ?, `+"`order`"+` = ?, project_id = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, color, isDefault, order, nullableProjectID(projectID), id)
	return err
}

func (r *ReferentialRepository) DeleteTicketStatus(id int) error {
	_, err := r.DB.Exec(`UPDATE ticket_statuses SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
}

func (r *ReferentialRepository) GetTicketTypes() ([]models.TicketTypeReference, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, icon, color, is_default, created_at
		FROM ticket_types
		WHERE deleted_at IS NULL
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketTypeReference
	for rows.Next() {
		var (
			item      models.TicketTypeReference
			createdAt time.Time
		)
		if err := rows.Scan(&item.ID, &item.Name, &item.Icon, &item.Color, &item.IsDefault, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ReferentialRepository) CreateTicketType(name, icon, color string, isDefault bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO ticket_types (name, icon, color, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, name, icon, color, isDefault)
	return err
}

func (r *ReferentialRepository) UpdateTicketType(id int, name, icon, color string, isDefault bool) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_types
		SET name = ?, icon = ?, color = ?, is_default = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, icon, color, isDefault, id)
	return err
}

func (r *ReferentialRepository) DeleteTicketType(id int) error {
	_, err := r.DB.Exec(`UPDATE ticket_types SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
}

func (r *ReferentialRepository) GetProjectOptions() ([]models.ProjectOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.ProjectOption
	for rows.Next() {
		var item models.ProjectOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		projects = append(projects, item)
	}

	return projects, rows.Err()
}

func (r *ReferentialRepository) getStatusReferences(table string) ([]models.StatusReference, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, color, is_default, created_at
		FROM ` + table + `
		WHERE deleted_at IS NULL
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.StatusReference
	for rows.Next() {
		var (
			item      models.StatusReference
			createdAt time.Time
		)
		if err := rows.Scan(&item.ID, &item.Name, &item.Color, &item.IsDefault, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ReferentialRepository) createStatusReference(table, name, color string, isDefault bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO `+table+` (name, color, is_default, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, name, color, isDefault)
	return err
}

func (r *ReferentialRepository) updateStatusReference(table string, id int, name, color string, isDefault bool) error {
	_, err := r.DB.Exec(`
		UPDATE `+table+`
		SET name = ?, color = ?, is_default = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, color, isDefault, id)
	return err
}

func (r *ReferentialRepository) deleteStatusReference(table string, id int) error {
	_, err := r.DB.Exec(`UPDATE `+table+` SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
}

func nullableProjectID(projectID *int) interface{} {
	if projectID == nil || *projectID <= 0 {
		return nil
	}
	return *projectID
}
