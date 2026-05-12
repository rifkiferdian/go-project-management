package repositories

import (
	"database/sql"
	"gobase-app/models"
	"strings"
	"time"
)

type ProjectRepository struct {
	DB *sql.DB
}

// GetAll mengambil seluruh project aktif beserta owner, status, dan ringkasan relasi.
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	rows, err := r.DB.Query(`
			SELECT
				p.id,
				p.name,
				COALESCE(p.description, '') AS description,
				p.owner_id,
				u.name AS owner_name,
				COALESCE(p.developer_id, 0) AS developer_id,
				COALESCE(dev.name, '-') AS developer_name,
				COALESCE(GROUP_CONCAT(DISTINCT d.id ORDER BY d.id SEPARATOR ','), '') AS request_division_ids_csv,
				COALESCE(GROUP_CONCAT(DISTINCT d.name ORDER BY d.name SEPARATOR ', '), '-') AS request_division,
				p.status_id,
			ps.name AS status_name,
			ps.color AS status_color,
			COALESCE(p.priority_id, 0) AS priority_id,
			COALESCE(pp.name, '-') AS priority_name,
			COALESCE(pp.color, '#cecece') AS priority_color,
			p.ticket_prefix,
			p.status_type,
			p.type,
			COUNT(DISTINCT pu.user_id) AS member_count,
			COUNT(DISTINCT t.id) AS ticket_count,
			p.created_at
			FROM projects p
			JOIN users u ON u.id = p.owner_id
			LEFT JOIN users dev ON dev.id = p.developer_id AND dev.deleted_at IS NULL
			JOIN project_statuses ps ON ps.id = p.status_id
			LEFT JOIN project_priorities pp ON pp.id = p.priority_id AND pp.deleted_at IS NULL
		LEFT JOIN project_divisions pd ON pd.project_id = p.id
		LEFT JOIN divisions d ON d.id = pd.division_id AND d.deleted_at IS NULL
		LEFT JOIN project_users pu ON pu.project_id = p.id
		LEFT JOIN tickets t ON t.project_id = p.id AND t.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
			GROUP BY
				p.id, p.name, p.description, p.owner_id, u.name, p.developer_id, dev.name, p.status_id, ps.name, ps.color, p.priority_id, pp.name, pp.color,
				p.ticket_prefix, p.status_type, p.type, p.created_at
			ORDER BY p.created_at DESC
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var (
			project               models.Project
			createdAt             time.Time
			requestDivisionIDsCSV string
		)

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.OwnerName,
			&project.DeveloperID,
			&project.DeveloperName,
			&requestDivisionIDsCSV,
			&project.RequestDivision,
			&project.StatusID,
			&project.StatusName,
			&project.StatusColor,
			&project.PriorityID,
			&project.PriorityName,
			&project.PriorityColor,
			&project.TicketPrefix,
			&project.StatusType,
			&project.Type,
			&project.MemberCount,
			&project.TicketCount,
			&createdAt,
		); err != nil {
			return nil, err
		}

		project.RequestDivisionIDs = splitCSVToIntSlice(requestDivisionIDsCSV)
		project.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		project.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

// GetByID mengambil detail project tunggal.
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	var project models.Project
	err := r.DB.QueryRow(`
			SELECT
				id,
				name,
				COALESCE(description, '') AS description,
				owner_id,
				COALESCE(developer_id, 0) AS developer_id,
				status_id,
				COALESCE(priority_id, 0) AS priority_id,
				ticket_prefix,
			status_type,
			type
		FROM projects
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.DeveloperID,
		&project.StatusID,
		&project.PriorityID,
		&project.TicketPrefix,
		&project.StatusType,
		&project.Type,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) ExistsByTicketPrefix(prefix string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM projects WHERE ticket_prefix = ? AND deleted_at IS NULL`, prefix).Scan(&count)
	return count > 0, err
}

func (r *ProjectRepository) ExistsByTicketPrefixExceptID(prefix string, id int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM projects WHERE ticket_prefix = ? AND id <> ? AND deleted_at IS NULL`, prefix, id).Scan(&count)
	return count > 0, err
}

func (r *ProjectRepository) Create(params models.ProjectCreateInput) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
			INSERT INTO projects (name, description, owner_id, developer_id, status_id, priority_id, ticket_prefix, status_type, type, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		`, params.Name, params.Description, params.OwnerID, params.DeveloperID, params.StatusID, params.PriorityID, params.TicketPrefix, params.StatusType, params.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	projectID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := r.replaceProjectDivisions(tx, int(projectID), params.DivisionIDs); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ProjectRepository) Update(params models.ProjectUpdateInput) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
			UPDATE projects
			SET name = ?, description = ?, owner_id = ?, developer_id = ?, status_id = ?, priority_id = ?, ticket_prefix = ?, status_type = ?, type = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
		`, params.Name, params.Description, params.OwnerID, params.DeveloperID, params.StatusID, params.PriorityID, params.TicketPrefix, params.StatusType, params.Type, params.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.replaceProjectDivisions(tx, params.ID, params.DivisionIDs); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ProjectRepository) Delete(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if err := ensureExists(tx, `SELECT COUNT(1) FROM projects WHERE id = ? AND deleted_at IS NULL`, id); err != nil {
		tx.Rollback()
		return err
	}

	steps := []struct {
		query string
		args  []interface{}
	}{
		{query: `DELETE FROM ticket_subscribers WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)`, args: []interface{}{id}},
		{query: `DELETE FROM ticket_hours WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)`, args: []interface{}{id}},
		{query: `DELETE FROM ticket_comments WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)`, args: []interface{}{id}},
		{query: `DELETE FROM ticket_attachments WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)`, args: []interface{}{id}},
		{query: `DELETE FROM ticket_activities WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)`, args: []interface{}{id}},
		{
			query: `DELETE FROM ticket_relations
				WHERE ticket_id IN (SELECT id FROM tickets WHERE project_id = ?)
					OR relation_id IN (SELECT id FROM tickets WHERE project_id = ?)`,
			args: []interface{}{id, id},
		},
		{query: `DELETE FROM tickets WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM sprints WHERE project_id = ?`, args: []interface{}{id}},
		{query: `UPDATE epics SET parent_id = NULL WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM epics WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM ticket_statuses WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM project_users WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM project_favorites WHERE project_id = ?`, args: []interface{}{id}},
		{query: `DELETE FROM project_divisions WHERE project_id = ?`, args: []interface{}{id}},
	}

	for _, step := range steps {
		if _, err := tx.Exec(step.query, step.args...); err != nil {
			tx.Rollback()
			return err
		}
	}

	result, err := tx.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if affectedRows == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (r *ProjectRepository) GetStatusOptions() ([]models.ProjectStatusOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, color
		FROM project_statuses
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.ProjectStatusOption
	for rows.Next() {
		var status models.ProjectStatusOption
		if err := rows.Scan(&status.ID, &status.Name, &status.Color); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

func (r *ProjectRepository) GetDivisionOptions() ([]models.DivisionOption, error) {
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

func (r *ProjectRepository) GetPriorityOptions() ([]models.ProjectPriorityOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, color
		FROM project_priorities
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priorities []models.ProjectPriorityOption
	for rows.Next() {
		var priority models.ProjectPriorityOption
		if err := rows.Scan(&priority.ID, &priority.Name, &priority.Color); err != nil {
			return nil, err
		}
		priorities = append(priorities, priority)
	}

	return priorities, rows.Err()
}

func (r *ProjectRepository) FindExistingDivisionIDs(ids []int64) (map[int64]bool, error) {
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

func (r *ProjectRepository) IsUserInDivision(userID int, divisionName string) (bool, error) {
	if userID <= 0 || strings.TrimSpace(divisionName) == "" {
		return false, nil
	}

	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM users u
		JOIN user_divisions ud ON ud.user_id = u.id
		JOIN divisions d ON d.id = ud.division_id
		WHERE u.id = ?
			AND u.deleted_at IS NULL
			AND d.deleted_at IS NULL
			AND LOWER(TRIM(d.name)) = LOWER(TRIM(?))
	`, userID, divisionName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProjectRepository) replaceProjectDivisions(tx *sql.Tx, projectID int, divisionIDs []int64) error {
	if _, err := tx.Exec(`DELETE FROM project_divisions WHERE project_id = ?`, projectID); err != nil {
		return err
	}

	if len(divisionIDs) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(`
		INSERT INTO project_divisions (project_id, division_id, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, divisionID := range divisionIDs {
		if _, err := stmt.Exec(projectID, divisionID); err != nil {
			return err
		}
	}

	return nil
}
