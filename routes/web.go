package routes

import (
	"gobase-app/controllers"
	"gobase-app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(r *gin.Engine) {
	r.Use(middleware.UserMiddleware())

	r.GET("/", controllers.LoginPage)
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginPost)
	r.POST("/register", controllers.CreateUser)
	r.GET("/logout", controllers.Logout)

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired(), middleware.PermissionContext())
	{
		auth.GET("/dashboard", controllers.DashboardIndex)
		auth.GET("/projects", middleware.RequirePermission("List projects"), controllers.ProjectIndex)
		auth.GET("/tickets", middleware.RequirePermission("List tickets"), controllers.TicketIndex)
		auth.GET("/tickets/:id", middleware.RequirePermission("List tickets"), controllers.TicketShow)
		auth.GET("/board", middleware.RequirePermission("List tickets"), controllers.BoardIndex)
		auth.GET("/road-map", middleware.RequirePermission("List sprints"), controllers.RoadMapIndex)
		auth.POST("/road-map/epics", middleware.RequirePermission("List sprints"), controllers.RoadMapEpicStore)
		auth.POST("/road-map/tickets", middleware.RequirePermission("Create ticket"), controllers.RoadMapTicketStore)
		auth.POST("/projects", middleware.RequirePermission("Create project"), controllers.ProjectStore)
		auth.POST("/projects/update", middleware.RequirePermission("Update project"), controllers.ProjectUpdate)
		auth.GET("/projects/delete/:id", middleware.RequirePermission("Delete project"), controllers.ProjectDelete)
		auth.GET("/activities", middleware.RequirePermission("List activities"), controllers.ActivityIndex)
		auth.POST("/activities", middleware.RequirePermission("Create activity"), controllers.ActivityStore)
		auth.POST("/activities/update", middleware.RequirePermission("Update activity"), controllers.ActivityUpdate)
		auth.GET("/activities/delete/:id", middleware.RequirePermission("Delete activity"), controllers.ActivityDelete)
		auth.GET("/project-statuses", middleware.RequirePermission("List project statuses"), controllers.ProjectStatusIndex)
		auth.POST("/project-statuses", middleware.RequirePermission("Create project status"), controllers.ProjectStatusStore)
		auth.POST("/project-statuses/update", middleware.RequirePermission("Update project status"), controllers.ProjectStatusUpdate)
		auth.GET("/project-statuses/delete/:id", middleware.RequirePermission("Delete project status"), controllers.ProjectStatusDelete)
		auth.GET("/ticket-statuses", middleware.RequirePermission("List ticket statuses"), controllers.TicketStatusIndex)
		auth.POST("/ticket-statuses", middleware.RequirePermission("Create ticket status"), controllers.TicketStatusStore)
		auth.POST("/ticket-statuses/update", middleware.RequirePermission("Update ticket status"), controllers.TicketStatusUpdate)
		auth.GET("/ticket-statuses/delete/:id", middleware.RequirePermission("Delete ticket status"), controllers.TicketStatusDelete)
		auth.GET("/ticket-types", middleware.RequirePermission("List ticket types"), controllers.TicketTypeIndex)
		auth.POST("/ticket-types", middleware.RequirePermission("Create ticket type"), controllers.TicketTypeStore)
		auth.POST("/ticket-types/update", middleware.RequirePermission("Update ticket type"), controllers.TicketTypeUpdate)
		auth.GET("/ticket-types/delete/:id", middleware.RequirePermission("Delete ticket type"), controllers.TicketTypeDelete)
		auth.GET("/ticket-priorities", middleware.RequirePermission("List ticket priorities"), controllers.TicketPriorityIndex)
		auth.POST("/ticket-priorities", middleware.RequirePermission("Create ticket priority"), controllers.TicketPriorityStore)
		auth.POST("/ticket-priorities/update", middleware.RequirePermission("Update ticket priority"), controllers.TicketPriorityUpdate)
		auth.GET("/ticket-priorities/delete/:id", middleware.RequirePermission("Delete ticket priority"), controllers.TicketPriorityDelete)
		auth.GET("/users", middleware.RequirePermission("List users"), controllers.UserIndex)
		auth.POST("/users", middleware.RequirePermission("Create user"), controllers.UserStore)
		auth.POST("/users/update", middleware.RequirePermission("Update user"), controllers.UserUpdate)
		auth.GET("/users/delete/:id", middleware.RequirePermission("Delete user"), controllers.UserDelete)
		auth.GET("/role", middleware.RequirePermission("List roles"), controllers.RoleIndex)
		auth.GET("/roleForm", controllers.RoleFormIndex)
		auth.GET("/role/:id/edit", middleware.RequirePermission("Update role"), controllers.RoleEdit)
		auth.POST("/role", middleware.RequirePermission("Create role"), controllers.RoleStore)
		auth.POST("/role/update", middleware.RequirePermission("Update role"), controllers.RoleUpdate)
		auth.GET("/role/delete/:id", middleware.RequirePermission("Delete role"), controllers.RoleDelete)
	}
}
