package controllers

import (
	"github.com/gin-gonic/gin"
	"gobase-app/config"
	"net/http"
)

func DashboardIndex(c *gin.Context) {
	var (
		totalProjects int
		totalTickets  int
		totalUsers    int
		totalRoles    int
	)

	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM projects WHERE deleted_at IS NULL`).Scan(&totalProjects); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM tickets WHERE deleted_at IS NULL`).Scan(&totalTickets); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE deleted_at IS NULL`).Scan(&totalUsers); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := config.DB.QueryRow(`SELECT COUNT(1) FROM roles`).Scan(&totalRoles); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "dashboard.html", gin.H{
		"Title":         "Dashboard",
		"Page":          "dashboard",
		"TotalProjects": totalProjects,
		"TotalTickets":  totalTickets,
		"TotalUsers":    totalUsers,
		"TotalRoles":    totalRoles,
	})

}
