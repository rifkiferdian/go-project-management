package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"gobase-app/config"
	"gobase-app/models"

	"github.com/gin-gonic/gin"
)

func PublicHomePage(c *gin.Context) {
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

	statusChartItems, statusChartTotal, _, err := getProjectStatusComposition()
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

	Render(c, "public_home.html", gin.H{
		"Title":                       "Home",
		"Page":                        "home",
		"TotalProjects":               totalProjects,
		"TotalTickets":                totalTickets,
		"TotalInProgressProjects":     totalInProgressProjects,
		"TotalImplementationProjects": totalImplementationProjects,
		"StatusChartItems":            statusChartItems,
		"StatusChartTotal":            statusChartTotal,
		"DivisionChartItems":          divisionChartItems,
		"UrgentProjects":              urgentProjects,
		"RequestQueueProjects":        requestQueueProjects,
		"KanbanStatuses":              kanbanStatuses,
		"KanbanProjects":              kanbanProjects,
		"RecentTicketActivities":      recentTicketActivities,
	})
}

func PublicProjectDetailPage(c *gin.Context) {
	projectID, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil || projectID <= 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code_error": http.StatusBadRequest,
			"error":      "project tidak valid",
		})
		return
	}

	svc := managementService()

	var (
		projectDescription string
		ownerID            int
		ownerName          sql.NullString
		developerID        sql.NullInt64
		developerName      sql.NullString
		statusID           int
		statusName         sql.NullString
		statusColor        sql.NullString
		priorityID         sql.NullInt64
		priorityName       sql.NullString
		priorityColor      sql.NullString
	)
	if err := config.DB.QueryRow(`
		SELECT
			COALESCE(p.description, '') AS description,
			p.owner_id,
			COALESCE(owner.name, '') AS owner_name,
			COALESCE(p.developer_id, 0) AS developer_id,
			COALESCE(dev.name, '') AS developer_name,
			p.status_id,
			COALESCE(ps.name, '') AS status_name,
			COALESCE(NULLIF(TRIM(ps.color), ''), '#cecece') AS status_color,
			COALESCE(p.priority_id, 0) AS priority_id,
			COALESCE(pp.name, '') AS priority_name,
			COALESCE(NULLIF(TRIM(pp.color), ''), '#cecece') AS priority_color
		FROM projects p
		JOIN users owner ON owner.id = p.owner_id
		LEFT JOIN users dev ON dev.id = p.developer_id
		LEFT JOIN project_statuses ps ON ps.id = p.status_id
		LEFT JOIN project_priorities pp ON pp.id = p.priority_id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`, projectID).Scan(
		&projectDescription,
		&ownerID,
		&ownerName,
		&developerID,
		&developerName,
		&statusID,
		&statusName,
		&statusColor,
		&priorityID,
		&priorityName,
		&priorityColor,
	); err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"code_error": http.StatusNotFound,
				"error":      "project tidak ditemukan",
			})
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	projectOptions, err := svc.GetRoadmapProjectOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	projectName := resolveRoadmapProjectLabel(projectOptions, projectID)
	if projectName == "" {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"code_error": http.StatusNotFound,
			"error":      "project tidak ditemukan",
		})
		return
	}

	columns, err := svc.GetBoardColumns(projectID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalTickets := 0
	for _, column := range columns {
		totalTickets += column.TicketCount
	}

	format := normalizeRoadmapFormat(c.DefaultQuery("format", "week"))
	epics, err := svc.GetRoadmapEpics()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tickets, err := svc.GetRoadmapTickets()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	filteredEpics, filteredTickets := filterRoadmapByProject(epics, tickets, projectID)
	weeks, rows, timelineWidth, currentMarkerLeft, currentMarkerWidth, columnWidth := svc.BuildRoadmapTimeline(filteredEpics, filteredTickets, roadmapNow(), format)
	yearGroups := buildRoadmapYearGroups(weeks, columnWidth)
	recentTicketActivities, err := getRecentTicketActivitiesByProject(projectID, 0)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "public_project_detail.html", gin.H{
		"Title":                  "Detail Project",
		"Page":                   "home",
		"ProjectID":              projectID,
		"ProjectName":            projectName,
		"ProjectDescription":     strings.TrimSpace(projectDescription),
		"ProjectOwnerID":         ownerID,
		"ProjectOwnerName":       strings.TrimSpace(ownerName.String),
		"ProjectDeveloperID":     int(developerID.Int64),
		"ProjectDeveloperName":   strings.TrimSpace(developerName.String),
		"ProjectStatusID":        statusID,
		"ProjectStatusName":      strings.TrimSpace(statusName.String),
		"ProjectStatusColor":     strings.TrimSpace(statusColor.String),
		"ProjectPriorityID":      int(priorityID.Int64),
		"ProjectPriorityName":    strings.TrimSpace(priorityName.String),
		"ProjectPriorityColor":   strings.TrimSpace(priorityColor.String),
		"Columns":                columns,
		"TotalStatus":            len(columns),
		"TotalTickets":           totalTickets,
		"Format":                 format,
		"Rows":                   rows,
		"Weeks":                  weeks,
		"YearGroups":             yearGroups,
		"TimelineWidth":          timelineWidth,
		"CurrentMarkerLeft":      currentMarkerLeft,
		"CurrentMarkerWidth":     currentMarkerWidth,
		"ColumnWidth":            columnWidth,
		"RoadmapEpicCount":       len(filteredEpics),
		"RoadmapTicketCount":     len(filteredTickets),
		"RecentTicketActivities": recentTicketActivities,
	})
}

func getRecentTicketActivitiesByProject(projectID int, limit int) ([]models.TicketActivityItem, error) {
	if projectID <= 0 {
		return []models.TicketActivityItem{}, nil
	}

	query := `
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
			AND p.id = ?
		ORDER BY ta.created_at DESC, ta.id DESC
	`
	args := []interface{}{projectID}
	if limit > 0 {
		query += "\nLIMIT ?"
		args = append(args, limit)
	}

	rows, err := config.DB.Query(query, args...)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return []models.TicketActivityItem{}, nil
	}

	return activities, nil
}
