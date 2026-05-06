package repositories

import (
	"database/sql"
	"gobase-app/models"
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
			p.status_id,
			ps.name AS status_name,
			ps.color AS status_color,
			p.ticket_prefix,
			p.status_type,
			p.type,
			COUNT(DISTINCT pu.user_id) AS member_count,
			COUNT(DISTINCT t.id) AS ticket_count,
			p.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		JOIN project_statuses ps ON ps.id = p.status_id
		LEFT JOIN project_users pu ON pu.project_id = p.id
		LEFT JOIN tickets t ON t.project_id = p.id AND t.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
		GROUP BY
			p.id, p.name, p.description, p.owner_id, u.name, p.status_id, ps.name, ps.color,
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
			project   models.Project
			createdAt time.Time
		)

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.OwnerName,
			&project.StatusID,
			&project.StatusName,
			&project.StatusColor,
			&project.TicketPrefix,
			&project.StatusType,
			&project.Type,
			&project.MemberCount,
			&project.TicketCount,
			&createdAt,
		); err != nil {
			return nil, err
		}

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
			status_id,
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
		&project.StatusID,
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
	_, err := r.DB.Exec(`
		INSERT INTO projects (name, description, owner_id, status_id, ticket_prefix, status_type, type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, params.Name, params.Description, params.OwnerID, params.StatusID, params.TicketPrefix, params.StatusType, params.Type)
	return err
}

func (r *ProjectRepository) Update(params models.ProjectUpdateInput) error {
	_, err := r.DB.Exec(`
		UPDATE projects
		SET name = ?, description = ?, owner_id = ?, status_id = ?, ticket_prefix = ?, status_type = ?, type = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, params.Name, params.Description, params.OwnerID, params.StatusID, params.TicketPrefix, params.StatusType, params.Type, params.ID)
	return err
}

func (r *ProjectRepository) Delete(id int) error {
	_, err := r.DB.Exec(`UPDATE projects SET deleted_at = NOW(), updated_at = NOW() WHERE id = ?`, id)
	return err
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
