package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gobase-app/config"
	helpers "gobase-app/helper"
	"gobase-app/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func DashboardIndex(c *gin.Context) {
	var (
		totalProjects               int
		totalTickets                int
		totalInProgressProjects     int
		totalImplementationProjects int
	)

	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM projects WHERE deleted_at IS NULL`).Scan(&totalProjects); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM tickets WHERE deleted_at IS NULL`).Scan(&totalTickets); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`
		SELECT COUNT(1)
		FROM projects p
		JOIN project_statuses ps ON ps.id = p.status_id
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
			AND LOWER(ps.name) = 'in progress'
	`).Scan(&totalInProgressProjects); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`
		SELECT COUNT(1)
		FROM projects p
		JOIN project_statuses ps ON ps.id = p.status_id
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
			AND LOWER(ps.name) = 'implementation'
	`).Scan(&totalImplementationProjects); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	kanbanStatuses, err := getProjectKanbanStatuses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	kanbanProjects, err := getProjectKanbanProjects()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	statusChartItems, statusChartTotal, statusChartGradient, err := getProjectStatusComposition()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	divisionChartItems, err := getProjectDivisionComposition()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	urgentProjects, err := getUrgentProjects(5)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	requestQueueProjects, err := getRequestQueueProjects(5)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	recentTicketActivities, err := getRecentTicketActivities()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "dashboard.html", gin.H{
		"Title":                       "Dashboard",
		"Page":                        "dashboard",
		"TotalProjects":               totalProjects,
		"TotalTickets":                totalTickets,
		"TotalInProgressProjects":     totalInProgressProjects,
		"TotalImplementationProjects": totalImplementationProjects,
		"StatusChartItems":            statusChartItems,
		"StatusChartTotal":            statusChartTotal,
		"StatusChartGradient":         statusChartGradient,
		"DivisionChartItems":          divisionChartItems,
		"UrgentProjects":              urgentProjects,
		"RequestQueueProjects":        requestQueueProjects,
		"KanbanStatuses":              kanbanStatuses,
		"KanbanProjects":              kanbanProjects,
		"RecentTicketActivities":      recentTicketActivities,
	})

}

func getProjectStatusComposition() ([]models.ProjectStatusChartItem, int, string, error) {
	rows, err := config.DB.Query(`
		SELECT
			ps.name,
			COALESCE(NULLIF(TRIM(ps.color), ''), '#cecece') AS color,
			COUNT(p.id) AS total
		FROM projects p
		JOIN project_statuses ps ON ps.id = p.status_id
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
		GROUP BY ps.id, ps.name, ps.color
		ORDER BY total DESC, ps.name ASC
	`)
	if err != nil {
		return nil, 0, "", err
	}
	defer rows.Close()

	var (
		items []models.ProjectStatusChartItem
		total int
	)

	for rows.Next() {
		var item models.ProjectStatusChartItem
		if err := rows.Scan(&item.Name, &item.Color, &item.Count); err != nil {
			return nil, 0, "", err
		}
		total += item.Count
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, "", err
	}

	if total == 0 || len(items) == 0 {
		return []models.ProjectStatusChartItem{}, 0, "conic-gradient(#e2e8f0 0 100%)", nil
	}

	assigned := 0
	for i := range items {
		if i == len(items)-1 {
			items[i].Percent = 100 - assigned
		} else {
			percent := int((float64(items[i].Count) / float64(total)) * 100.0)
			items[i].Percent = percent
			assigned += percent
		}
	}

	var (
		from  = 0
		parts []string
	)
	for _, item := range items {
		to := from + item.Percent
		if to > 100 {
			to = 100
		}
		if to < from {
			to = from
		}
		parts = append(parts, fmt.Sprintf("%s %d%% %d%%", item.Color, from, to))
		from = to
	}
	if from < 100 && len(items) > 0 {
		lastColor := items[len(items)-1].Color
		parts = append(parts, fmt.Sprintf("%s %d%% 100%%", lastColor, from))
	}

	return items, total, "conic-gradient(" + strings.Join(parts, ", ") + ")", nil
}

func getProjectDivisionComposition() ([]models.ProjectDivisionChartItem, error) {
	rows, err := config.DB.Query(`
		SELECT
			d.name,
			COUNT(DISTINCT pd.project_id) AS total
		FROM project_divisions pd
		JOIN divisions d ON d.id = pd.division_id
		JOIN projects p ON p.id = pd.project_id
		WHERE d.deleted_at IS NULL
			AND p.deleted_at IS NULL
		GROUP BY d.id, d.name
		ORDER BY total DESC, d.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ProjectDivisionChartItem
	for rows.Next() {
		var item models.ProjectDivisionChartItem
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return []models.ProjectDivisionChartItem{}, nil
	}

	maxCount := 0
	for _, item := range items {
		if item.Count > maxCount {
			maxCount = item.Count
		}
	}
	if maxCount <= 0 {
		for i := range items {
			items[i].WidthPercent = 0
		}
		return items, nil
	}

	for i := range items {
		items[i].WidthPercent = int((float64(items[i].Count) / float64(maxCount)) * 100.0)
		if items[i].WidthPercent == 0 && items[i].Count > 0 {
			items[i].WidthPercent = 4
		}
		if items[i].WidthPercent > 100 {
			items[i].WidthPercent = 100
		}
	}

	return items, nil
}

func getUrgentProjects(limit int) ([]models.DashboardProjectListItem, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := config.DB.Query(`
		SELECT
			p.id,
			p.name,
			COALESCE(NULLIF(GROUP_CONCAT(DISTINCT d.name ORDER BY d.name SEPARATOR ', '), ''), '-') AS request_division,
			ps.name AS status_name,
			COALESCE(ps.color, '#64748b') AS status_color,
			1 AS high_priority_ticket_count
		FROM projects p
		JOIN project_statuses ps ON ps.id = p.status_id
		LEFT JOIN project_priorities pp ON pp.id = p.priority_id AND pp.deleted_at IS NULL
		LEFT JOIN project_divisions pd ON pd.project_id = p.id
		LEFT JOIN divisions d ON d.id = pd.division_id AND d.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
			AND LOWER(COALESCE(pp.name, '')) IN ('high', 'urgent', 'critical', 'tinggi')
		GROUP BY p.id, p.name, ps.name, ps.color, p.created_at
		ORDER BY p.created_at ASC, p.id ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.DashboardProjectListItem
	for rows.Next() {
		var item models.DashboardProjectListItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RequestDivision,
			&item.StatusName,
			&item.StatusColor,
			&item.HighPriorityTicketCount,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func getRequestQueueProjects(limit int) ([]models.DashboardProjectListItem, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := config.DB.Query(`
		SELECT
			p.id,
			p.name,
			COALESCE(NULLIF(GROUP_CONCAT(DISTINCT d.name ORDER BY d.name SEPARATOR ', '), ''), '-') AS request_division,
			ps.name AS status_name,
			COALESCE(ps.color, '#64748b') AS status_color
		FROM projects p
		JOIN project_statuses ps ON ps.id = p.status_id
		LEFT JOIN project_divisions pd ON pd.project_id = p.id
		LEFT JOIN divisions d ON d.id = pd.division_id AND d.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
			AND LOWER(TRIM(COALESCE(ps.name, ''))) = 'request received'
		GROUP BY p.id, p.name, ps.name, ps.color, p.created_at
		ORDER BY p.created_at ASC, p.id ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.DashboardProjectListItem
	for rows.Next() {
		var item models.DashboardProjectListItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RequestDivision,
			&item.StatusName,
			&item.StatusColor,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func getProjectKanbanStatuses() ([]models.ProjectStatusOption, error) {
	rows, err := config.DB.Query(`
		SELECT id, name, COALESCE(NULLIF(TRIM(color), ''), '#cecece') AS color
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
		var item models.ProjectStatusOption
		if err := rows.Scan(&item.ID, &item.Name, &item.Color); err != nil {
			return nil, err
		}
		statuses = append(statuses, item)
	}

	return statuses, rows.Err()
}

func getProjectKanbanProjects() ([]models.Project, error) {
	rows, err := config.DB.Query(`
			SELECT
				p.id,
				p.name,
				COALESCE(NULLIF(GROUP_CONCAT(DISTINCT d.name ORDER BY d.name SEPARATOR ', '), ''), '-') AS request_division,
				COALESCE(owner.name, '-') AS owner_name,
				COALESCE(dev.name, '-') AS developer_name,
				COALESCE(pp.name, '-') AS priority_name,
				COALESCE(NULLIF(TRIM(pp.color), ''), '#cecece') AS priority_color,
				p.status_id,
				ps.name AS status_name,
				COALESCE(NULLIF(TRIM(ps.color), ''), '#cecece') AS status_color,
				p.created_at
			FROM projects p
			JOIN project_statuses ps ON ps.id = p.status_id
			LEFT JOIN users owner ON owner.id = p.owner_id AND owner.deleted_at IS NULL
			LEFT JOIN users dev ON dev.id = p.developer_id AND dev.deleted_at IS NULL
			LEFT JOIN project_priorities pp ON pp.id = p.priority_id AND pp.deleted_at IS NULL
			LEFT JOIN project_divisions pd ON pd.project_id = p.id
			LEFT JOIN divisions d ON d.id = pd.division_id AND d.deleted_at IS NULL
			WHERE p.deleted_at IS NULL
				AND ps.deleted_at IS NULL
			GROUP BY
				p.id, p.name, owner.name, dev.name, pp.name, pp.color, ps.id, ps.name, ps.color, p.created_at
			ORDER BY p.created_at DESC, p.id DESC
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
			&project.RequestDivision,
			&project.OwnerName,
			&project.DeveloperName,
			&project.PriorityName,
			&project.PriorityColor,
			&project.StatusID,
			&project.StatusName,
			&project.StatusColor,
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

func getRecentTicketActivities() ([]models.TicketActivityItem, error) {
	rows, err := config.DB.Query(`
		SELECT
			ta.id,
			t.id AS ticket_id,
			t.code,
			t.name AS ticket_name,
			p.name AS project_name,
			u.name AS user_name,
			old_ts.name AS old_status_name,
			COALESCE(old_ts.color, '#cecece') AS old_status_color,
			new_ts.name AS new_status_name,
			COALESCE(new_ts.color, '#cecece') AS new_status_color,
			ta.created_at
		FROM ticket_activities ta
		JOIN tickets t ON t.id = ta.ticket_id
		JOIN projects p ON p.id = t.project_id
		JOIN users u ON u.id = ta.user_id
		JOIN ticket_statuses old_ts ON old_ts.id = ta.old_status_id
		JOIN ticket_statuses new_ts ON new_ts.id = ta.new_status_id
		WHERE t.deleted_at IS NULL
			AND p.deleted_at IS NULL
		ORDER BY ta.created_at DESC, ta.id DESC
			LIMIT 10
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.TicketActivityItem
	for rows.Next() {
		var (
			item         models.TicketActivityItem
			createdAtRaw sql.NullTime
		)

		if err := rows.Scan(
			&item.ID,
			&item.TicketID,
			&item.TicketCode,
			&item.TicketName,
			&item.ProjectName,
			&item.UserName,
			&item.OldStatusName,
			&item.OldStatusColor,
			&item.NewStatusName,
			&item.NewStatusColor,
			&createdAtRaw,
		); err != nil {
			return nil, err
		}

		item.UserInitials = dashboardInitials(item.UserName)
		item.CreatedAtDisplay = dashboardTimestamp(createdAtRaw)
		item.CreatedAtRelative = dashboardRelativeTime(createdAtRaw)
		activities = append(activities, item)
	}

	return activities, rows.Err()
}

func dashboardInitials(name string) string {
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

func dashboardTimestamp(value sql.NullTime) string {
	if !value.Valid {
		return "-"
	}
	return value.Time.Format("02 Jan 06")
}

func dashboardRelativeTime(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}

	diff := time.Since(value.Time)
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
