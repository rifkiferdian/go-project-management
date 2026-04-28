package repositories

import (
	"database/sql"
	"fmt"
	helpers "gobase-app/helper"
	"gobase-app/models"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ManagementRepository struct {
	DB *sql.DB
}

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

func (r *ManagementRepository) GetTickets() ([]models.TicketListItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			t.id,
			t.code,
			t.name,
			p.name AS project_name,
			ts.name AS status_name,
			COALESCE(ts.color, '#cecece') AS status_color,
			tp.name AS priority_name,
			COALESCE(tp.color, '#cecece') AS priority_color,
			tt.name AS type_name,
			COALESCE(tt.color, '#cecece') AS type_color,
			owner.name AS owner_name,
			COALESCE(responsible.name, '-') AS responsible_name,
			COALESCE(t.estimation, 0) AS estimation,
			t.starts_at,
			t.ends_at,
			t.updated_at
		FROM tickets t
		JOIN projects p ON p.id = t.project_id
		JOIN ticket_statuses ts ON ts.id = t.status_id
		JOIN ticket_priorities tp ON tp.id = t.priority_id
		JOIN ticket_types tt ON tt.id = t.type_id
		JOIN users owner ON owner.id = t.owner_id
		LEFT JOIN users responsible ON responsible.id = t.responsible_id
		WHERE t.deleted_at IS NULL
		ORDER BY t.updated_at DESC, t.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketListItem
	for rows.Next() {
		var (
			item        models.TicketListItem
			startsAtRaw sql.NullTime
			endsAtRaw   sql.NullTime
			updatedAt   time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.ProjectName,
			&item.StatusName,
			&item.StatusColor,
			&item.PriorityName,
			&item.PriorityColor,
			&item.TypeName,
			&item.TypeColor,
			&item.OwnerName,
			&item.ResponsibleName,
			&item.Estimation,
			&startsAtRaw,
			&endsAtRaw,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		item.EstimationText = formatHours(item.Estimation)
		item.StartsAtDisplay = formatOptionalDate(startsAtRaw)
		item.EndsAtDisplay = formatOptionalDate(endsAtRaw)
		item.UpdatedAtDisplay = updatedAt.Format("02 Jan 2006 15:04")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetTicketDetailPage(id int) (models.TicketDetailPage, error) {
	detail, err := r.GetTicketDetail(id)
	if err != nil {
		return models.TicketDetailPage{}, err
	}

	comments, err := r.GetTicketComments(id)
	if err != nil {
		return models.TicketDetailPage{}, err
	}

	activities, err := r.GetTicketActivities(id)
	if err != nil {
		return models.TicketDetailPage{}, err
	}

	hours, err := r.GetTicketHours(id)
	if err != nil {
		return models.TicketDetailPage{}, err
	}

	subscribers, err := r.GetTicketSubscribers(id)
	if err != nil {
		return models.TicketDetailPage{}, err
	}

	return models.TicketDetailPage{
		Ticket:      detail,
		Comments:    comments,
		Activities:  activities,
		Hours:       hours,
		Subscribers: subscribers,
	}, nil
}

func (r *ManagementRepository) GetTicketDetail(id int) (models.TicketDetail, error) {
	var (
		item         models.TicketDetail
		contentRaw   string
		startsAtRaw  sql.NullTime
		endsAtRaw    sql.NullTime
		createdAtRaw sql.NullTime
		updatedAtRaw sql.NullTime
	)

	err := r.DB.QueryRow(`
		SELECT
			t.id,
			t.code,
			t.name,
			t.content,
			p.name AS project_name,
			ts.name AS status_name,
			COALESCE(ts.color, '#cecece') AS status_color,
			tp.name AS priority_name,
			COALESCE(tp.color, '#cecece') AS priority_color,
			tt.name AS type_name,
			COALESCE(tt.color, '#cecece') AS type_color,
			owner.name AS owner_name,
			COALESCE(responsible.name, '-') AS responsible_name,
			COALESCE(e.name, '-') AS epic_name,
			COALESCE(t.estimation, 0) AS estimation,
			COALESCE((
				SELECT SUM(th.value)
				FROM ticket_hours th
				WHERE th.ticket_id = t.id
			), 0) AS logged_hours,
			COALESCE((
				SELECT COUNT(1)
				FROM ticket_subscribers sub
				WHERE sub.ticket_id = t.id
			), 0) AS subscribers_count,
			t.starts_at,
			t.ends_at,
			t.created_at,
			t.updated_at
		FROM tickets t
		JOIN projects p ON p.id = t.project_id
		JOIN ticket_statuses ts ON ts.id = t.status_id
		JOIN ticket_priorities tp ON tp.id = t.priority_id
		JOIN ticket_types tt ON tt.id = t.type_id
		JOIN users owner ON owner.id = t.owner_id
		LEFT JOIN users responsible ON responsible.id = t.responsible_id
		LEFT JOIN epics e ON e.id = t.epic_id
		WHERE t.id = ? AND t.deleted_at IS NULL
	`, id).Scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&contentRaw,
		&item.ProjectName,
		&item.StatusName,
		&item.StatusColor,
		&item.PriorityName,
		&item.PriorityColor,
		&item.TypeName,
		&item.TypeColor,
		&item.OwnerName,
		&item.ResponsibleName,
		&item.EpicName,
		&item.Estimation,
		&item.LoggedHours,
		&item.SubscribersCount,
		&startsAtRaw,
		&endsAtRaw,
		&createdAtRaw,
		&updatedAtRaw,
	)
	if err != nil {
		return models.TicketDetail{}, err
	}

	item.ContentText = plainTextFromHTML(contentRaw)
	item.OwnerInitials = initialsOrFallback(item.OwnerName)
	item.ResponsibleInitials = initialsOrFallback(item.ResponsibleName)
	item.EstimationText = formatHoursLabel(item.Estimation)
	item.LoggedHoursText = formatHoursLabel(item.LoggedHours)
	item.LoggedPercent = percentFromHours(item.LoggedHours, item.Estimation)
	item.StartsAtDisplay = formatOptionalDate(startsAtRaw)
	item.EndsAtDisplay = formatOptionalDate(endsAtRaw)
	item.CreatedAtDisplay = formatTimestamp(createdAtRaw)
	item.CreatedAtRelative = relativeTime(createdAtRaw)
	item.UpdatedAtDisplay = formatTimestamp(updatedAtRaw)
	item.UpdatedAtRelative = relativeTime(updatedAtRaw)

	return item, nil
}

func (r *ManagementRepository) GetTicketComments(ticketID int) ([]models.TicketCommentItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			tc.id,
			u.name,
			tc.content,
			tc.created_at
		FROM ticket_comments tc
		JOIN users u ON u.id = tc.user_id
		WHERE tc.ticket_id = ? AND tc.deleted_at IS NULL
		ORDER BY tc.created_at ASC, tc.id ASC
	`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketCommentItem
	for rows.Next() {
		var (
			item         models.TicketCommentItem
			createdAtRaw sql.NullTime
		)
		if err := rows.Scan(&item.ID, &item.UserName, &item.Content, &createdAtRaw); err != nil {
			return nil, err
		}
		item.UserInitials = initialsOrFallback(item.UserName)
		item.CreatedAtDisplay = formatTimestamp(createdAtRaw)
		item.CreatedAtRelative = relativeTime(createdAtRaw)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetTicketActivities(ticketID int) ([]models.TicketActivityItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			ta.id,
			u.name,
			old_ts.name AS old_status_name,
			new_ts.name AS new_status_name,
			ta.created_at
		FROM ticket_activities ta
		JOIN users u ON u.id = ta.user_id
		JOIN ticket_statuses old_ts ON old_ts.id = ta.old_status_id
		JOIN ticket_statuses new_ts ON new_ts.id = ta.new_status_id
		WHERE ta.ticket_id = ?
		ORDER BY ta.created_at DESC, ta.id DESC
	`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketActivityItem
	for rows.Next() {
		var (
			item         models.TicketActivityItem
			createdAtRaw sql.NullTime
		)
		if err := rows.Scan(&item.ID, &item.UserName, &item.OldStatusName, &item.NewStatusName, &createdAtRaw); err != nil {
			return nil, err
		}
		item.UserInitials = initialsOrFallback(item.UserName)
		item.CreatedAtDisplay = formatTimestamp(createdAtRaw)
		item.CreatedAtRelative = relativeTime(createdAtRaw)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetTicketHours(ticketID int) ([]models.TicketHourItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			th.id,
			u.name,
			COALESCE(a.name, '-') AS activity_name,
			COALESCE(th.comment, '') AS comment,
			th.value,
			th.created_at
		FROM ticket_hours th
		JOIN users u ON u.id = th.user_id
		LEFT JOIN activities a ON a.id = th.activity_id
		WHERE th.ticket_id = ?
		ORDER BY th.created_at DESC, th.id DESC
	`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketHourItem
	for rows.Next() {
		var (
			item         models.TicketHourItem
			createdAtRaw sql.NullTime
		)
		if err := rows.Scan(&item.ID, &item.UserName, &item.ActivityName, &item.Comment, &item.Value, &createdAtRaw); err != nil {
			return nil, err
		}
		item.UserInitials = initialsOrFallback(item.UserName)
		item.ValueText = formatHours(item.Value)
		item.CreatedAtDisplay = formatTimestamp(createdAtRaw)
		item.CreatedAtRelative = relativeTime(createdAtRaw)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetTicketSubscribers(ticketID int) ([]models.TicketSubscriberItem, error) {
	rows, err := r.DB.Query(`
		SELECT
			u.id,
			u.name
		FROM ticket_subscribers ts
		JOIN users u ON u.id = ts.user_id
		WHERE ts.ticket_id = ?
		ORDER BY u.name ASC
	`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TicketSubscriberItem
	for rows.Next() {
		var item models.TicketSubscriberItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		item.Initials = initialsOrFallback(item.Name)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetBoardColumns() ([]models.BoardColumn, error) {
	rows, err := r.DB.Query(`
		SELECT
			ts.id,
			ts.name,
			COALESCE(ts.color, '#cecece') AS color,
			COALESCE(p.name, 'Global') AS scope_label,
			ts.` + "`order`" + `
		FROM ticket_statuses ts
		LEFT JOIN projects p ON p.id = ts.project_id
		WHERE ts.deleted_at IS NULL
		ORDER BY ts.project_id IS NOT NULL, ts.` + "`order`" + `, ts.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []models.BoardColumn
	for rows.Next() {
		var item models.BoardColumn
		if err := rows.Scan(&item.ID, &item.Name, &item.Color, &item.ScopeLabel, &item.Order); err != nil {
			return nil, err
		}
		columns = append(columns, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	tickets, err := r.getBoardTickets()
	if err != nil {
		return nil, err
	}

	grouped := map[int][]models.BoardTicket{}
	for _, ticket := range tickets {
		grouped[ticket.StatusID] = append(grouped[ticket.StatusID], ticket)
	}

	for i := range columns {
		columns[i].Tickets = grouped[columns[i].ID]
		columns[i].TicketCount = len(columns[i].Tickets)
	}

	return columns, nil
}

func (r *ManagementRepository) GetRoadmapEpics() ([]models.RoadmapEpic, error) {
	rows, err := r.DB.Query(`
		SELECT
			e.id,
			e.project_id,
			e.name,
			p.name AS project_name,
			e.starts_at,
			e.ends_at,
			(
				SELECT COUNT(1)
				FROM sprints s
				WHERE s.epic_id = e.id AND s.deleted_at IS NULL
			) AS sprint_count,
			(
				SELECT COUNT(1)
				FROM tickets t
				WHERE t.epic_id = e.id AND t.deleted_at IS NULL
			) AS ticket_count,
			(
				SELECT COUNT(1)
				FROM tickets t
				JOIN ticket_statuses ts ON ts.id = t.status_id
				WHERE t.epic_id = e.id
					AND t.deleted_at IS NULL
					AND LOWER(ts.name) IN ('done', 'closed')
			) AS done_count
		FROM epics e
		JOIN projects p ON p.id = e.project_id
		WHERE e.deleted_at IS NULL
		ORDER BY e.starts_at ASC, e.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RoadmapEpic
	for rows.Next() {
		var (
			item     models.RoadmapEpic
			startsAt time.Time
			endsAt   time.Time
		)
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.Name,
			&item.ProjectName,
			&startsAt,
			&endsAt,
			&item.SprintCount,
			&item.TicketCount,
			&item.DoneCount,
		); err != nil {
			return nil, err
		}
		item.StartsAtISO = startsAt.Format("2006-01-02")
		item.EndsAtISO = endsAt.Format("2006-01-02")
		item.StartsAt = startsAt.Format("02 Jan 2006")
		item.EndsAt = endsAt.Format("02 Jan 2006")
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetRoadmapSprints() ([]models.RoadmapSprint, error) {
	rows, err := r.DB.Query(`
		SELECT
			s.id,
			s.name,
			p.name AS project_name,
			COALESCE(s.epic_id, 0) AS epic_id,
			COALESCE(e.name, '-') AS epic_name,
			s.starts_at,
			s.ends_at,
			s.started_at,
			s.ended_at,
			(
				SELECT COUNT(1)
				FROM tickets t
				WHERE t.sprint_id = s.id AND t.deleted_at IS NULL
			) AS ticket_count,
			(
				SELECT COUNT(1)
				FROM tickets t
				JOIN ticket_statuses ts ON ts.id = t.status_id
				WHERE t.sprint_id = s.id
					AND t.deleted_at IS NULL
					AND LOWER(ts.name) IN ('done', 'closed')
			) AS done_count
		FROM sprints s
		JOIN projects p ON p.id = s.project_id
		LEFT JOIN epics e ON e.id = s.epic_id
		WHERE s.deleted_at IS NULL
		ORDER BY s.starts_at ASC, s.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RoadmapSprint
	for rows.Next() {
		var (
			item         models.RoadmapSprint
			startsAt     time.Time
			endsAt       time.Time
			startedAtRaw sql.NullTime
			endedAtRaw   sql.NullTime
		)
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.ProjectName,
			&item.EpicID,
			&item.EpicName,
			&startsAt,
			&endsAt,
			&startedAtRaw,
			&endedAtRaw,
			&item.TicketCount,
			&item.DoneCount,
		); err != nil {
			return nil, err
		}
		item.StartsAtISO = startsAt.Format("2006-01-02")
		item.EndsAtISO = endsAt.Format("2006-01-02")
		item.StartsAt = startsAt.Format("02 Jan 2006")
		item.EndsAt = endsAt.Format("02 Jan 2006")
		switch {
		case endedAtRaw.Valid:
			item.StateLabel = "Finished"
		case startedAtRaw.Valid:
			item.StateLabel = "Active"
		default:
			item.StateLabel = "Planned"
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) CountRoadmapProjects() (int, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM projects WHERE deleted_at IS NULL`).Scan(&count)
	return count, err
}

func (r *ManagementRepository) GetRoadmapProjectOptions() ([]models.ProjectOption, error) {
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

	var items []models.ProjectOption
	for rows.Next() {
		var item models.ProjectOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) GetRoadmapEpicOptions() ([]models.RoadmapEpicOption, error) {
	rows, err := r.DB.Query(`
		SELECT e.id, e.name, e.project_id, p.name
		FROM epics e
		JOIN projects p ON p.id = e.project_id
		WHERE e.deleted_at IS NULL
		ORDER BY p.name, e.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RoadmapEpicOption
	for rows.Next() {
		var item models.RoadmapEpicOption
		if err := rows.Scan(&item.ID, &item.Name, &item.ProjectID, &item.ProjectName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) CreateRoadmapEpic(input models.RoadmapEpicCreateInput) error {
	_, err := r.DB.Exec(`
		INSERT INTO epics (project_id, name, starts_at, ends_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, input.ProjectID, input.Name, input.StartsAt, input.EndsAt)
	return err
}

func (r *ManagementRepository) CreateRoadmapTicket(input models.RoadmapTicketCreateInput) error {
	statusID, err := r.defaultTicketStatusID(input.ProjectID)
	if err != nil {
		return err
	}

	typeID, err := r.defaultTicketTypeID()
	if err != nil {
		return err
	}

	priorityID, err := r.defaultTicketPriorityID()
	if err != nil {
		return err
	}

	code, err := r.nextTicketCode(input.ProjectID)
	if err != nil {
		return err
	}

	nextOrder, err := r.nextTicketOrder(input.ProjectID)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`
		INSERT INTO tickets (
			name, content, owner_id, responsible_id, status_id, project_id,
			created_at, updated_at, code, type_id, `+"`order`"+`, priority_id, estimation, starts_at, ends_at, epic_id
		)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?)
	`, input.Name, input.Name, input.ResourceUserID, input.ResourceUserID, statusID, input.ProjectID, code, typeID, nextOrder, priorityID, nullableEstimation(input.Estimation), input.StartsAt, input.EndsAt, nullableIntPointer(input.EpicID))
	return err
}

func (r *ManagementRepository) GetRoadmapTickets() ([]models.RoadmapTicket, error) {
	rows, err := r.DB.Query(`
		SELECT
			t.id,
			COALESCE(t.epic_id, 0) AS epic_id,
			t.project_id,
			t.name,
			p.name AS project_name,
			COALESCE(responsible.name, owner.name, '-') AS resource_name,
			t.starts_at,
			t.ends_at,
			ts.name AS status_name
		FROM tickets t
		JOIN projects p ON p.id = t.project_id
		JOIN ticket_statuses ts ON ts.id = t.status_id
		JOIN users owner ON owner.id = t.owner_id
		LEFT JOIN users responsible ON responsible.id = t.responsible_id
		WHERE t.deleted_at IS NULL
		ORDER BY t.epic_id ASC, t.` + "`order`" + ` ASC, t.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RoadmapTicket
	for rows.Next() {
		var (
			item        models.RoadmapTicket
			startsAtRaw sql.NullTime
			endsAtRaw   sql.NullTime
			statusName  string
		)
		if err := rows.Scan(&item.ID, &item.EpicID, &item.ProjectID, &item.Name, &item.ProjectName, &item.ResourceName, &startsAtRaw, &endsAtRaw, &statusName); err != nil {
			return nil, err
		}
		item.Progress = roadmapTicketProgress(statusName)
		item.StartsAtISO = optionalDateISO(startsAtRaw)
		item.EndsAtISO = optionalDateISO(endsAtRaw)
		item.StartsAt = formatOptionalDate(startsAtRaw)
		item.EndsAt = formatOptionalDate(endsAtRaw)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ManagementRepository) getBoardTickets() ([]models.BoardTicket, error) {
	rows, err := r.DB.Query(`
		SELECT
			t.id,
			t.code,
			t.name,
			p.name AS project_name,
			tp.name AS priority_name,
			COALESCE(tp.color, '#cecece') AS priority_color,
			tt.name AS type_name,
			COALESCE(tt.color, '#cecece') AS type_color,
			COALESCE(responsible.name, '-') AS responsible_name,
			COALESCE(t.estimation, 0) AS estimation,
			t.status_id
		FROM tickets t
		JOIN projects p ON p.id = t.project_id
		JOIN ticket_priorities tp ON tp.id = t.priority_id
		JOIN ticket_types tt ON tt.id = t.type_id
		LEFT JOIN users responsible ON responsible.id = t.responsible_id
		WHERE t.deleted_at IS NULL
		ORDER BY t.project_id ASC, t.` + "`order`" + ` ASC, t.updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.BoardTicket
	for rows.Next() {
		var (
			item       models.BoardTicket
			estimation float64
		)
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.ProjectName,
			&item.PriorityName,
			&item.PriorityColor,
			&item.TypeName,
			&item.TypeColor,
			&item.ResponsibleName,
			&estimation,
			&item.StatusID,
		); err != nil {
			return nil, err
		}
		item.EstimationText = formatHours(estimation)
		items = append(items, item)
	}

	return items, rows.Err()
}

func formatHours(value float64) string {
	if value <= 0 {
		return "-"
	}
	return trimFloat(value) + "h"
}

func trimFloat(value float64) string {
	if value == float64(int64(value)) {
		return strconv.Itoa(int(value))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
}

func roadmapTicketProgress(statusName string) int {
	switch strings.ToLower(strings.TrimSpace(statusName)) {
	case "done", "closed":
		return 100
	case "in progress", "review", "testing":
		return 50
	default:
		return 1
	}
}

func (r *ManagementRepository) defaultTicketStatusID(projectID int) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		SELECT id
		FROM ticket_statuses
		WHERE deleted_at IS NULL
			AND (project_id IS NULL OR project_id = ?)
		ORDER BY
			CASE WHEN project_id = ? THEN 0 ELSE 1 END,
			CASE WHEN is_default = 1 THEN 0 ELSE 1 END,
			`+"`order`"+` ASC,
			id ASC
		LIMIT 1
	`, projectID, projectID).Scan(&id)
	return id, err
}

func (r *ManagementRepository) defaultTicketTypeID() (int, error) {
	var id int
	err := r.DB.QueryRow(`
		SELECT id
		FROM ticket_types
		WHERE deleted_at IS NULL
		ORDER BY CASE WHEN is_default = 1 THEN 0 ELSE 1 END, id ASC
		LIMIT 1
	`).Scan(&id)
	return id, err
}

func (r *ManagementRepository) defaultTicketPriorityID() (int, error) {
	var id int
	err := r.DB.QueryRow(`
		SELECT id
		FROM ticket_priorities
		WHERE deleted_at IS NULL
		ORDER BY CASE WHEN is_default = 1 THEN 0 ELSE 1 END, id ASC
		LIMIT 1
	`).Scan(&id)
	return id, err
}

func (r *ManagementRepository) nextTicketCode(projectID int) (string, error) {
	var (
		prefix string
		maxSeq int
	)

	if err := r.DB.QueryRow(`SELECT ticket_prefix FROM projects WHERE id = ? AND deleted_at IS NULL`, projectID).Scan(&prefix); err != nil {
		return "", err
	}

	if err := r.DB.QueryRow(`
		SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(code, '-', -1) AS UNSIGNED)), 0)
		FROM tickets
		WHERE project_id = ? AND deleted_at IS NULL
	`, projectID).Scan(&maxSeq); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%d", strings.ToUpper(strings.TrimSpace(prefix)), maxSeq+1), nil
}

func (r *ManagementRepository) nextTicketOrder(projectID int) (int, error) {
	var maxOrder int
	if err := r.DB.QueryRow(`SELECT COALESCE(MAX(`+"`order`"+`), 0) FROM tickets WHERE project_id = ? AND deleted_at IS NULL`, projectID).Scan(&maxOrder); err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

func nullableEstimation(value float64) interface{} {
	if value <= 0 {
		return nil
	}
	return value
}

func nullableIntPointer(value *int) interface{} {
	if value == nil || *value <= 0 {
		return nil
	}
	return *value
}

func formatOptionalDate(value sql.NullTime) string {
	if !value.Valid {
		return "-"
	}
	return value.Time.Format("02 Jan 2006")
}

func optionalDateISO(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format("2006-01-02")
}

func formatTimestamp(value sql.NullTime) string {
	if !value.Valid {
		return "-"
	}
	return value.Time.Format("2006-01-02 3:04 PM")
}

func relativeTime(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}

	now := time.Now()
	diff := now.Sub(value.Time)
	if diff < 0 {
		diff = 0
	}

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return strconv.Itoa(int(diff.Minutes())) + " minutes ago"
	case diff < 24*time.Hour:
		return strconv.Itoa(int(diff.Hours())) + " hours ago"
	default:
		return strconv.Itoa(int(diff.Hours()/24)) + " days ago"
	}
}

func percentFromHours(total, estimation float64) int {
	if estimation <= 0 {
		if total > 0 {
			return 100
		}
		return 0
	}
	value := int((total / estimation) * 100)
	if value < 0 {
		return 0
	}
	return value
}

func plainTextFromHTML(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n\n",
		"</div>", "\n",
	)
	normalized := replacer.Replace(value)
	cleaned := htmlTagPattern.ReplaceAllString(normalized, "")
	cleaned = html.UnescapeString(cleaned)
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return "-"
	}
	return cleaned
}

func formatHoursLabel(value float64) string {
	if value <= 0 {
		return "-"
	}
	if value == 1 {
		return "1 hour"
	}
	return trimFloat(value) + " hours"
}

func initialsOrFallback(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed == "-" {
		return "NA"
	}
	value := helpers.Initials(trimmed)
	if value == "" {
		return "NA"
	}
	return value
}
