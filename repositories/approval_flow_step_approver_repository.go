package repositories

import (
	"database/sql"
	"fmt"
	"gobase-app/models"
	"time"
)

type ApprovalFlowStepApproverRepository struct {
	DB *sql.DB
}

func (r *ApprovalFlowStepApproverRepository) GetStepOptions() ([]models.ApprovalFlowStepOption, error) {
	rows, err := r.DB.Query(`
		SELECT
			s.id,
			s.approval_flow_id,
			s.step_order,
			s.step_name,
			f.flow_code,
			f.flow_name,
			s.is_active
		FROM approval_flow_steps s
		JOIN approval_flows f ON f.id = s.approval_flow_id
		ORDER BY f.flow_name ASC, s.step_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.ApprovalFlowStepOption
	for rows.Next() {
		var item models.ApprovalFlowStepOption
		if err := rows.Scan(
			&item.ID,
			&item.ApprovalFlowID,
			&item.StepOrder,
			&item.StepName,
			&item.FlowCode,
			&item.FlowName,
			&item.IsActive,
		); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func (r *ApprovalFlowStepApproverRepository) ExistsStepByID(stepID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flow_steps
		WHERE id = ?
	`, stepID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepApproverRepository) GetAll() ([]models.ApprovalFlowStepApprover, error) {
	rows, err := r.DB.Query(`
		SELECT
			a.id,
			a.approval_flow_step_id,
			s.step_order,
			s.step_name,
			f.flow_code,
			f.flow_name,
			a.approver_type,
			COALESCE(a.approver_user_id, 0) AS approver_user_id,
			COALESCE(a.approver_role_id, 0) AS approver_role_id,
			COALESCE(a.approver_division_id, 0) AS approver_division_id,
			CASE
				WHEN a.approver_type = 'user' THEN COALESCE(u.name, CONCAT('User #', COALESCE(a.approver_user_id, 0)))
				WHEN a.approver_type = 'role' THEN COALESCE(r.name, CONCAT('Role #', COALESCE(a.approver_role_id, 0)))
				WHEN a.approver_type = 'division' THEN COALESCE(d.name, CONCAT('Division #', COALESCE(a.approver_division_id, 0)))
				ELSE '-'
			END AS approver_label,
			a.is_active,
			a.created_at,
			a.updated_at
		FROM approval_flow_step_approvers a
		JOIN approval_flow_steps s ON s.id = a.approval_flow_step_id
		JOIN approval_flows f ON f.id = s.approval_flow_id
		LEFT JOIN users u ON u.id = a.approver_user_id AND u.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = a.approver_role_id
		LEFT JOIN divisions d ON d.id = a.approver_division_id AND d.deleted_at IS NULL
		ORDER BY f.flow_name ASC, s.step_order ASC, a.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowsOut []models.ApprovalFlowStepApprover
	for rows.Next() {
		var (
			item      models.ApprovalFlowStepApprover
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(
			&item.ID,
			&item.ApprovalFlowStepID,
			&item.StepOrder,
			&item.StepName,
			&item.FlowCode,
			&item.FlowName,
			&item.ApproverType,
			&item.ApproverUserID,
			&item.ApproverRoleID,
			&item.ApproverDivisionID,
			&item.ApproverLabel,
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

		rowsOut = append(rowsOut, item)
	}

	return rowsOut, rows.Err()
}

func (r *ApprovalFlowStepApproverRepository) ExistsDuplicate(stepID int, approverType string, userID, roleID, divisionID int, excludeID int) (bool, error) {
	var (
		column string
		value  int
	)

	switch approverType {
	case "user":
		column = "approver_user_id"
		value = userID
	case "role":
		column = "approver_role_id"
		value = roleID
	case "division":
		column = "approver_division_id"
		value = divisionID
	default:
		return false, fmt.Errorf("tipe approver tidak valid")
	}

	query := `
		SELECT COUNT(1)
		FROM approval_flow_step_approvers
		WHERE approval_flow_step_id = ?
			AND approver_type = ?
			AND ` + column + ` = ?`

	args := []interface{}{stepID, approverType, value}
	if excludeID > 0 {
		query += ` AND id <> ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.DB.QueryRow(query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ApprovalFlowStepApproverRepository) ExistsUserByID(userID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`, userID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepApproverRepository) ExistsRoleByID(roleID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM roles
		WHERE id = ?
	`, roleID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepApproverRepository) ExistsDivisionByID(divisionID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE id = ? AND deleted_at IS NULL
	`, divisionID).Scan(&count)
	return count > 0, err
}

func (r *ApprovalFlowStepApproverRepository) Create(stepID int, approverType string, userID, roleID, divisionID int, isActive bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO approval_flow_step_approvers (
			approval_flow_step_id,
			approver_type,
			approver_user_id,
			approver_role_id,
			approver_division_id,
			is_active,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, stepID, approverType, nullablePositiveInt(userID), nullablePositiveInt(roleID), nullablePositiveInt(divisionID), isActive)
	return err
}

func (r *ApprovalFlowStepApproverRepository) Update(id, stepID int, approverType string, userID, roleID, divisionID int, isActive bool) error {
	_, err := r.DB.Exec(`
		UPDATE approval_flow_step_approvers
		SET
			approval_flow_step_id = ?,
			approver_type = ?,
			approver_user_id = ?,
			approver_role_id = ?,
			approver_division_id = ?,
			is_active = ?,
			updated_at = NOW()
		WHERE id = ?
	`, stepID, approverType, nullablePositiveInt(userID), nullablePositiveInt(roleID), nullablePositiveInt(divisionID), isActive, id)
	return err
}

func (r *ApprovalFlowStepApproverRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM approval_flow_step_approvers
		WHERE id = ?
	`, id)
	return err
}

func (r *ApprovalFlowStepApproverRepository) GetUserOptions() ([]models.LookupOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.LookupOption
	for rows.Next() {
		var item models.LookupOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func (r *ApprovalFlowStepApproverRepository) GetRoleOptions() ([]models.LookupOption, error) {
	rows, err := r.DB.Query(`
		SELECT id, name
		FROM roles
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.LookupOption
	for rows.Next() {
		var item models.LookupOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func (r *ApprovalFlowStepApproverRepository) GetDivisionOptions() ([]models.LookupOption, error) {
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

	var options []models.LookupOption
	for rows.Next() {
		var item models.LookupOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func nullablePositiveInt(value int) interface{} {
	if value <= 0 {
		return nil
	}
	return value
}
