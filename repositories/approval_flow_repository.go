package repositories

import (
	"database/sql"
	"gobase-app/models"
	"time"
)

type ApprovalFlowRepository struct {
	DB *sql.DB
}

func (r *ApprovalFlowRepository) GetAll() ([]models.ApprovalFlow, error) {
	rows, err := r.DB.Query(`
		SELECT
			f.id,
			COALESCE(f.division_id, 0),
			COALESCE(d.name, 'Belum ditentukan'),
			f.flow_code,
			f.flow_name,
			f.is_active,
			f.created_at,
			f.updated_at
		FROM approval_flows f
		LEFT JOIN divisions d ON d.id = f.division_id AND d.deleted_at IS NULL
		ORDER BY f.is_active DESC, d.name ASC, f.flow_name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flows []models.ApprovalFlow
	for rows.Next() {
		var (
			item      models.ApprovalFlow
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(
			&item.ID,
			&item.DivisionID,
			&item.DivisionName,
			&item.FlowCode,
			&item.FlowName,
			&item.IsActive,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		if updatedAt.Valid {
			item.UpdatedAtDisplay = updatedAt.Time.Format("02 Jan 2006 15:04:05")
		} else {
			item.UpdatedAtDisplay = "-"
		}

		flows = append(flows, item)
	}

	return flows, rows.Err()
}

func (r *ApprovalFlowRepository) ExistsByCode(flowCode string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flows
		WHERE flow_code = ?
	`, flowCode).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowRepository) ExistsByCodeExceptID(flowCode string, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flows
		WHERE flow_code = ?
			AND id <> ?
	`, flowCode, excludeID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowRepository) Create(divisionID int, flowCode, flowName string, isActive bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO approval_flows (division_id, flow_code, flow_name, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, divisionID, flowCode, flowName, isActive)
	return err
}

func (r *ApprovalFlowRepository) Update(id, divisionID int, flowCode, flowName string, isActive bool) error {
	_, err := r.DB.Exec(`
		UPDATE approval_flows
		SET division_id = ?, flow_code = ?, flow_name = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`, divisionID, flowCode, flowName, isActive, id)
	return err
}

func (r *ApprovalFlowRepository) DivisionExists(divisionID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE id = ? AND deleted_at IS NULL
	`, divisionID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM approval_flows
		WHERE id = ?
	`, id)
	return err
}
