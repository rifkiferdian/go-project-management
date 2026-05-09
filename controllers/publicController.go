package controllers

import (
	"net/http"

	"gobase-app/config"

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
