package repositories

import (
	"database/sql"
	"gobase-app/models"
	"time"
)

type ApprovalFlowStepRepository struct {
	DB *sql.DB
}

func (r *ApprovalFlowStepRepository) GetFlowOptions() ([]models.ApprovalFlowOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, flow_code, flow_name, is_active
		FROM approval_flows
		ORDER BY is_active DESC, flow_name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.ApprovalFlowOption
	for rows.Next() {
		var item models.ApprovalFlowOption
		if err := rows.Scan(&item.ID, &item.FlowCode, &item.FlowName, &item.IsActive); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func (r *ApprovalFlowStepRepository) ExistsFlowByID(flowID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flows
		WHERE id = ?
	`, flowID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepRepository) GetAll() ([]models.ApprovalFlowStep, error) {
	return r.GetAllByFlowID(0)
}

func (r *ApprovalFlowStepRepository) GetAllByFlowID(flowID int) ([]models.ApprovalFlowStep, error) {
	query := `
		SELECT
			s.id,
			s.approval_flow_id,
			f.flow_code,
			f.flow_name,
			s.step_order,
			s.step_name,
			s.approval_rule,
			s.is_active,
			s.created_at,
			s.updated_at
		FROM approval_flow_steps s
		JOIN approval_flows f ON f.id = s.approval_flow_id
	`
	args := make([]interface{}, 0, 1)
	if flowID > 0 {
		query += " WHERE s.approval_flow_id = ?"
		args = append(args, flowID)
	}
	query += " ORDER BY f.flow_name ASC, s.step_order ASC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []models.ApprovalFlowStep
	for rows.Next() {
		var (
			item      models.ApprovalFlowStep
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(
			&item.ID,
			&item.ApprovalFlowID,
			&item.FlowCode,
			&item.FlowName,
			&item.StepOrder,
			&item.StepName,
			&item.ApprovalRule,
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
		steps = append(steps, item)
	}

	return steps, rows.Err()
}

func (r *ApprovalFlowStepRepository) ExistsStepOrder(flowID, stepOrder int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flow_steps
		WHERE approval_flow_id = ?
			AND step_order = ?
	`, flowID, stepOrder).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepRepository) ExistsStepOrderExceptID(flowID, stepOrder, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flow_steps
		WHERE approval_flow_id = ?
			AND step_order = ?
			AND id <> ?
	`, flowID, stepOrder, excludeID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepRepository) Create(flowID, stepOrder int, stepName, approvalRule string, isActive bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO approval_flow_steps (approval_flow_id, step_order, step_name, approval_rule, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`, flowID, stepOrder, stepName, approvalRule, isActive)
	return err
}

func (r *ApprovalFlowStepRepository) Update(id, flowID, stepOrder int, stepName, approvalRule string, isActive bool) error {
	_, err := r.DB.Exec(`
		UPDATE approval_flow_steps
		SET approval_flow_id = ?, step_order = ?, step_name = ?, approval_rule = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`, flowID, stepOrder, stepName, approvalRule, isActive, id)
	return err
}

func (r *ApprovalFlowStepRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM approval_flow_steps
		WHERE id = ?
	`, id)
	return err
}
