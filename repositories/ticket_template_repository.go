package repositories

import (
	"database/sql"
	"fmt"
	"gobase-app/models"
	"time"
)

type TicketTemplateRepository struct {
	DB *sql.DB
}

func (r *TicketTemplateRepository) GetSets() ([]models.TicketTemplateSet, error) {
	rows, err := r.DB.Query(`
		SELECT
			s.id,
			s.name,
			s.purpose,
			COALESCE(s.description, '') AS description,
			s.is_active,
			(
				SELECT COUNT(1)
				FROM ticket_template_epics e
				WHERE e.set_id = s.id AND e.deleted_at IS NULL
			) AS epic_count,
			(
				SELECT COUNT(1)
				FROM ticket_template_items i
				WHERE i.set_id = s.id AND i.deleted_at IS NULL
			) AS item_count,
			s.created_at
		FROM ticket_template_sets s
		WHERE s.deleted_at IS NULL
		ORDER BY s.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []models.TicketTemplateSet
	for rows.Next() {
		var (
			item      models.TicketTemplateSet
			createdAt time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Purpose,
			&item.Description,
			&item.IsActive,
			&item.EpicCount,
			&item.ItemCount,
			&createdAt,
		); err != nil {
			return nil, err
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		sets = append(sets, item)
	}

	return sets, rows.Err()
}

func (r *TicketTemplateRepository) CreateSet(name, purpose, description string, isActive bool) error {
	_, err := r.DB.Exec(`
		INSERT INTO ticket_template_sets (name, purpose, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, name, purpose, templateNullableString(description), isActive)
	return err
}

func (r *TicketTemplateRepository) UpdateSet(id int, name, purpose, description string, isActive bool) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_template_sets
		SET name = ?, purpose = ?, description = ?, is_active = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, name, purpose, templateNullableString(description), isActive, id)
	return err
}

func (r *TicketTemplateRepository) DeleteSet(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ticket_template_sets
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ticket_template_items
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE set_id = ? AND deleted_at IS NULL
	`, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ticket_template_epics
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE set_id = ? AND deleted_at IS NULL
	`, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TicketTemplateRepository) SetExists(id int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_sets
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) ExistsSetByNamePurpose(name, purpose string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_sets
		WHERE name = ? AND purpose = ? AND deleted_at IS NULL
	`, name, purpose).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) ExistsSetByNamePurposeExceptID(name, purpose string, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_sets
		WHERE name = ? AND purpose = ? AND id <> ? AND deleted_at IS NULL
	`, name, purpose, excludeID).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) GetEpics() ([]models.TicketTemplateEpic, error) {
	rows, err := r.DB.Query(`
		SELECT
			e.id,
			e.set_id,
			s.name AS set_name,
			s.purpose AS set_purpose,
			e.name,
			COALESCE(e.description, '') AS description,
			e.start_offset_days,
			e.due_offset_days,
			e.sort_order,
			e.is_active,
			e.created_at
		FROM ticket_template_epics e
		INNER JOIN ticket_template_sets s ON s.id = e.set_id AND s.deleted_at IS NULL
		WHERE e.deleted_at IS NULL
		ORDER BY s.name ASC, e.sort_order ASC, e.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketTemplateEpic
	for rows.Next() {
		var (
			item            models.TicketTemplateEpic
			startOffsetDays sql.NullInt64
			dueOffsetDays   sql.NullInt64
			createdAt       time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.SetID,
			&item.SetName,
			&item.SetPurpose,
			&item.Name,
			&item.Description,
			&startOffsetDays,
			&dueOffsetDays,
			&item.SortOrder,
			&item.IsActive,
			&createdAt,
		); err != nil {
			return nil, err
		}
		if startOffsetDays.Valid {
			item.StartOffsetDays = int(startOffsetDays.Int64)
		}
		if dueOffsetDays.Valid {
			item.DueOffsetDays = int(dueOffsetDays.Int64)
		}
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *TicketTemplateRepository) CreateEpic(
	setID int,
	name, description string,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	_, err := r.DB.Exec(`
		INSERT INTO ticket_template_epics (
			set_id, name, description, start_offset_days, due_offset_days, sort_order, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		setID,
		name,
		templateNullableString(description),
		templateNullableInt(startOffsetDays),
		templateNullableInt(dueOffsetDays),
		sortOrder,
		isActive,
	)
	return err
}

func (r *TicketTemplateRepository) UpdateEpic(
	id, setID int,
	name, description string,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_template_epics
		SET
			set_id = ?,
			name = ?,
			description = ?,
			start_offset_days = ?,
			due_offset_days = ?,
			sort_order = ?,
			is_active = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`,
		setID,
		name,
		templateNullableString(description),
		templateNullableInt(startOffsetDays),
		templateNullableInt(dueOffsetDays),
		sortOrder,
		isActive,
		id,
	)
	return err
}

func (r *TicketTemplateRepository) DeleteEpic(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ticket_template_epics
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ticket_template_items
		SET template_epic_id = NULL, updated_at = NOW()
		WHERE template_epic_id = ? AND deleted_at IS NULL
	`, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TicketTemplateRepository) EpicExists(id int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_epics
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) ExistsEpicBySetAndName(setID int, name string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_epics
		WHERE set_id = ? AND name = ? AND deleted_at IS NULL
	`, setID, name).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) ExistsEpicBySetAndNameExceptID(setID int, name string, excludeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_epics
		WHERE set_id = ? AND name = ? AND id <> ? AND deleted_at IS NULL
	`, setID, name, excludeID).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) EpicBelongsToSet(epicID, setID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_epics
		WHERE id = ? AND set_id = ? AND deleted_at IS NULL
	`, epicID, setID).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) GetItems() ([]models.TicketTemplateItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			i.id,
			i.set_id,
			s.name AS set_name,
			s.purpose AS set_purpose,
			i.title,
			COALESCE(i.description, '') AS description,
			i.template_epic_id,
			COALESCE(e.name, '') AS template_epic_name,
			i.default_type_id,
			COALESCE(tt.name, '') AS type_name,
			i.default_priority_id,
			COALESCE(tp.name, '') AS priority_name,
			i.default_status_id,
			COALESCE(ts.name, '') AS status_name,
			i.default_owner_id,
			COALESCE(uo.name, '') AS owner_name,
			i.default_responsible_id,
			COALESCE(ur.name, '') AS responsible_name,
			i.estimation,
			i.start_offset_days,
			i.due_offset_days,
			i.sort_order,
			i.is_active,
			i.created_at
		FROM ticket_template_items i
		INNER JOIN ticket_template_sets s ON s.id = i.set_id AND s.deleted_at IS NULL
		LEFT JOIN ticket_template_epics e ON e.id = i.template_epic_id AND e.deleted_at IS NULL
		LEFT JOIN ticket_types tt ON tt.id = i.default_type_id
		LEFT JOIN ticket_priorities tp ON tp.id = i.default_priority_id
		LEFT JOIN ticket_statuses ts ON ts.id = i.default_status_id
		LEFT JOIN users uo ON uo.id = i.default_owner_id
		LEFT JOIN users ur ON ur.id = i.default_responsible_id
		WHERE i.deleted_at IS NULL
		ORDER BY s.name ASC, i.sort_order ASC, i.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketTemplateItem
	for rows.Next() {
		var (
			item              models.TicketTemplateItem
			templateEpicID    sql.NullInt64
			defaultTypeID     sql.NullInt64
			defaultPriorityID sql.NullInt64
			defaultStatusID   sql.NullInt64
			defaultOwnerID    sql.NullInt64
			defaultRespID     sql.NullInt64
			estimation        sql.NullFloat64
			startOffsetDays   sql.NullInt64
			dueOffsetDays     sql.NullInt64
			createdAt         time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.SetID,
			&item.SetName,
			&item.SetPurpose,
			&item.Title,
			&item.Description,
			&templateEpicID,
			&item.TemplateEpicName,
			&defaultTypeID,
			&item.DefaultTypeName,
			&defaultPriorityID,
			&item.DefaultPriorityName,
			&defaultStatusID,
			&item.DefaultStatusName,
			&defaultOwnerID,
			&item.DefaultOwnerName,
			&defaultRespID,
			&item.DefaultResponsibleName,
			&estimation,
			&startOffsetDays,
			&dueOffsetDays,
			&item.SortOrder,
			&item.IsActive,
			&createdAt,
		); err != nil {
			return nil, err
		}

		if templateEpicID.Valid {
			item.TemplateEpicID = int(templateEpicID.Int64)
		}
		if defaultTypeID.Valid {
			item.DefaultTypeID = int(defaultTypeID.Int64)
		}
		if defaultPriorityID.Valid {
			item.DefaultPriorityID = int(defaultPriorityID.Int64)
		}
		if defaultStatusID.Valid {
			item.DefaultStatusID = int(defaultStatusID.Int64)
		}
		if defaultOwnerID.Valid {
			item.DefaultOwnerID = int(defaultOwnerID.Int64)
		}
		if defaultRespID.Valid {
			item.DefaultResponsibleID = int(defaultRespID.Int64)
		}
		if estimation.Valid {
			item.Estimation = estimation.Float64
			item.EstimationText = fmt.Sprintf("%.2f", estimation.Float64)
		} else {
			item.EstimationText = "-"
		}
		if startOffsetDays.Valid {
			item.StartOffsetDays = int(startOffsetDays.Int64)
		}
		if dueOffsetDays.Valid {
			item.DueOffsetDays = int(dueOffsetDays.Int64)
		}

		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *TicketTemplateRepository) CreateItem(
	setID int,
	title, description string,
	templateEpicID, defaultTypeID, defaultPriorityID, defaultStatusID, defaultOwnerID, defaultResponsibleID *int,
	estimation *float64,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	_, err := r.DB.Exec(`
		INSERT INTO ticket_template_items (
			set_id, title, description, template_epic_id, default_type_id, default_priority_id, default_status_id,
			default_owner_id, default_responsible_id, estimation, start_offset_days, due_offset_days,
			sort_order, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		setID,
		title,
		templateNullableString(description),
		templateNullableInt(templateEpicID),
		templateNullableInt(defaultTypeID),
		templateNullableInt(defaultPriorityID),
		templateNullableInt(defaultStatusID),
		templateNullableInt(defaultOwnerID),
		templateNullableInt(defaultResponsibleID),
		templateNullableFloat(estimation),
		templateNullableInt(startOffsetDays),
		templateNullableInt(dueOffsetDays),
		sortOrder,
		isActive,
	)
	return err
}

func (r *TicketTemplateRepository) UpdateItem(
	id, setID int,
	title, description string,
	templateEpicID, defaultTypeID, defaultPriorityID, defaultStatusID, defaultOwnerID, defaultResponsibleID *int,
	estimation *float64,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_template_items
		SET
			set_id = ?,
			title = ?,
			description = ?,
			template_epic_id = ?,
			default_type_id = ?,
			default_priority_id = ?,
			default_status_id = ?,
			default_owner_id = ?,
			default_responsible_id = ?,
			estimation = ?,
			start_offset_days = ?,
			due_offset_days = ?,
			sort_order = ?,
			is_active = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`,
		setID,
		title,
		templateNullableString(description),
		templateNullableInt(templateEpicID),
		templateNullableInt(defaultTypeID),
		templateNullableInt(defaultPriorityID),
		templateNullableInt(defaultStatusID),
		templateNullableInt(defaultOwnerID),
		templateNullableInt(defaultResponsibleID),
		templateNullableFloat(estimation),
		templateNullableInt(startOffsetDays),
		templateNullableInt(dueOffsetDays),
		sortOrder,
		isActive,
		id,
	)
	return err
}

func (r *TicketTemplateRepository) DeleteItem(id int) error {
	_, err := r.DB.Exec(`
		UPDATE ticket_template_items
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, id)
	return err
}

func (r *TicketTemplateRepository) ItemExists(id int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM ticket_template_items
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&count)
	return count > 0, err
}

func (r *TicketTemplateRepository) GetTicketTypeOptions() ([]models.TicketTemplateOption, error) {
	return r.getSimpleOptions(`
		SELECT id, name
		FROM ticket_types
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
}

func (r *TicketTemplateRepository) GetTicketPriorityOptions() ([]models.TicketTemplateOption, error) {
	return r.getSimpleOptions(`
		SELECT id, name
		FROM ticket_priorities
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
}

func (r *TicketTemplateRepository) GetTicketStatusOptions() ([]models.TicketTemplateOption, error) {
	return r.getSimpleOptions(`
		SELECT ts.id,
			   CASE
			   		WHEN ts.project_id IS NULL THEN ts.name
			   		ELSE CONCAT(ts.name, ' (', COALESCE(p.name, 'Project'), ')')
			   END AS name
		FROM ticket_statuses ts
		LEFT JOIN projects p ON p.id = ts.project_id
		WHERE ts.deleted_at IS NULL
		ORDER BY ts.project_id IS NOT NULL, ts.` + "`order`" + `, ts.name
	`)
}

func (r *TicketTemplateRepository) GetUserOptions() ([]models.TicketTemplateOption, error) {
	return r.getSimpleOptions(`
		SELECT id, name
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
}

func (r *TicketTemplateRepository) GetEpicOptions() ([]models.TicketTemplateOption, error) {
	return r.getSimpleOptions(`
		SELECT e.id, CONCAT(e.name, ' (', s.name, ')') AS name
		FROM ticket_template_epics e
		INNER JOIN ticket_template_sets s ON s.id = e.set_id
		WHERE e.deleted_at IS NULL
			AND s.deleted_at IS NULL
		ORDER BY s.name ASC, e.sort_order ASC, e.name ASC
	`)
}

func (r *TicketTemplateRepository) getSimpleOptions(query string) ([]models.TicketTemplateOption, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketTemplateOption
	for rows.Next() {
		var item models.TicketTemplateOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func templateNullableInt(value *int) interface{} {
	if value == nil || *value <= 0 {
		return nil
	}
	return *value
}

func templateNullableFloat(value *float64) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func templateNullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
