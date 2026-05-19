package controllers

import (
	"database/sql"
	"gobase-app/config"
	"gobase-app/models"
)

func userDivisionIDSet(userID int) (map[int]bool, error) {
	result := map[int]bool{}
	if userID <= 0 {
		return result, nil
	}

	rows, err := config.DB.Query(`
		SELECT DISTINCT ud.division_id
		FROM user_divisions ud
		JOIN divisions d ON d.id = ud.division_id AND d.deleted_at IS NULL
		JOIN users u ON u.id = ud.user_id AND u.deleted_at IS NULL
		WHERE ud.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var divisionID int
		if err := rows.Scan(&divisionID); err != nil {
			return nil, err
		}
		if divisionID > 0 {
			result[divisionID] = true
		}
	}

	return result, rows.Err()
}

func userManageableProjectIDSet(userID int) (map[int]bool, error) {
	result := map[int]bool{}
	if userID <= 0 {
		return result, nil
	}

	rows, err := config.DB.Query(`
		SELECT DISTINCT p.id
		FROM projects p
		JOIN project_divisions pd ON pd.project_id = p.id
		JOIN user_divisions ud ON ud.division_id = pd.division_id
		WHERE p.deleted_at IS NULL
			AND ud.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var projectID int
		if err := rows.Scan(&projectID); err != nil {
			return nil, err
		}
		if projectID > 0 {
			result[projectID] = true
		}
	}

	return result, rows.Err()
}

func userCanManageProjectByID(userID, projectID int) (bool, error) {
	if userID <= 0 || projectID <= 0 {
		return false, nil
	}

	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(1)
		FROM projects p
		JOIN project_divisions pd ON pd.project_id = p.id
		JOIN user_divisions ud ON ud.division_id = pd.division_id
		WHERE p.id = ?
			AND p.deleted_at IS NULL
			AND ud.user_id = ?
	`, projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func userCanManageTicketByID(userID, ticketID int) (bool, error) {
	if userID <= 0 || ticketID <= 0 {
		return false, nil
	}

	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(1)
		FROM tickets t
		JOIN projects p ON p.id = t.project_id AND p.deleted_at IS NULL
		JOIN project_divisions pd ON pd.project_id = p.id
		JOIN user_divisions ud ON ud.division_id = pd.division_id
		WHERE t.id = ?
			AND t.deleted_at IS NULL
			AND ud.user_id = ?
	`, ticketID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func userCanAssignProjectDivisions(userID int, divisionIDs []int64) (bool, error) {
	if userID <= 0 || len(divisionIDs) == 0 {
		return false, nil
	}

	userDivisionSet, err := userDivisionIDSet(userID)
	if err != nil {
		return false, err
	}
	if len(userDivisionSet) == 0 {
		return false, nil
	}

	for _, divisionID := range divisionIDs {
		if divisionID <= 0 {
			return false, nil
		}
		if !userDivisionSet[int(divisionID)] {
			return false, nil
		}
	}

	return true, nil
}

func filterDivisionOptionsByUser(options []models.DivisionOption, userID int) ([]models.DivisionOption, error) {
	if userID <= 0 {
		return []models.DivisionOption{}, nil
	}

	allowedSet, err := userDivisionIDSet(userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.DivisionOption, 0, len(options))
	for _, option := range options {
		if allowedSet[option.ID] {
			result = append(result, option)
		}
	}
	return result, nil
}

func filterProjectOptionsByUser(options []models.ProjectOption, userID int) ([]models.ProjectOption, map[int]bool, error) {
	allowedSet, err := userManageableProjectIDSet(userID)
	if err != nil {
		return nil, nil, err
	}

	result := make([]models.ProjectOption, 0, len(options))
	for _, option := range options {
		if allowedSet[option.ID] {
			result = append(result, option)
		}
	}

	return result, allowedSet, nil
}

func ticketProjectID(ticketID int) (int, error) {
	if ticketID <= 0 {
		return 0, sql.ErrNoRows
	}

	var projectID int
	err := config.DB.QueryRow(`
		SELECT project_id
		FROM tickets
		WHERE id = ? AND deleted_at IS NULL
	`, ticketID).Scan(&projectID)
	if err != nil {
		return 0, err
	}
	return projectID, nil
}
