package controllers

import (
	"database/sql"
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

	implementationProjects, err := getImplementationProjects()
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
		"ImplementationProjects":      implementationProjects,
		"RecentTicketActivities":      recentTicketActivities,
	})

}

func getImplementationProjects() ([]models.Project, error) {
	rows, err := config.DB.Query(`
		SELECT
			p.id,
			p.name,
			COALESCE(p.description, '') AS description,
			u.name AS owner_name,
			ps.name AS status_name,
			ps.color AS status_color,
			p.ticket_prefix,
			p.type,
			COUNT(DISTINCT t.id) AS ticket_count,
			p.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		JOIN project_statuses ps ON ps.id = p.status_id
		LEFT JOIN tickets t ON t.project_id = p.id AND t.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
			AND ps.deleted_at IS NULL
			AND LOWER(ps.name) = 'implementation'
		GROUP BY
			p.id, p.name, p.description, u.name, ps.name, ps.color,
			p.ticket_prefix, p.type, p.created_at
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
			&project.OwnerName,
			&project.StatusName,
			&project.StatusColor,
			&project.TicketPrefix,
			&project.Type,
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
		LIMIT 8
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
