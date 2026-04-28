package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"gobase-app/config"
	"gobase-app/repositories"
	"gobase-app/services"

	"github.com/gin-gonic/gin"
)

type statusPageData struct {
	Title          string
	Page           string
	EntityLabel    string
	EntityPlural   string
	Rows           interface{}
	Error          string
	CreateAction   string
	UpdateAction   string
	DeleteBasePath string
}

func referentialService() *services.ReferentialService {
	return &services.ReferentialService{
		Repo: &repositories.ReferentialRepository{DB: config.DB},
	}
}

func ActivityIndex(c *gin.Context) {
	svc := referentialService()
	rows, err := svc.GetActivities()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "activities.html", gin.H{
		"Title": "Activities",
		"Page":  "activity",
		"Rows":  rows,
	})
}

func ActivityStore(c *gin.Context) {
	svc := referentialService()
	if err := svc.CreateActivity(c.PostForm("name"), c.PostForm("description")); err != nil {
		renderActivitiesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/activities")
}

func ActivityUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderActivitiesWithError(c, "activity tidak valid")
		return
	}

	svc := referentialService()
	if err := svc.UpdateActivity(id, c.PostForm("name"), c.PostForm("description")); err != nil {
		renderActivitiesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/activities")
}

func ActivityDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid activity id")
		return
	}
	svc := referentialService()
	if err := svc.DeleteActivity(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/activities")
}

func ProjectStatusIndex(c *gin.Context) {
	svc := referentialService()
	rows, err := svc.GetProjectStatuses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	renderStatusPage(c, statusPageData{
		Title:          "Project statuses",
		Page:           "projectStatus",
		EntityLabel:    "Project status",
		EntityPlural:   "Project statuses",
		Rows:           rows,
		CreateAction:   "/project-statuses",
		UpdateAction:   "/project-statuses/update",
		DeleteBasePath: "/project-statuses/delete/",
	})
}

func ProjectStatusStore(c *gin.Context) {
	svc := referentialService()
	if err := svc.CreateProjectStatus(c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderProjectStatusesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/project-statuses")
}

func ProjectStatusUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderProjectStatusesWithError(c, "project status tidak valid")
		return
	}

	svc := referentialService()
	if err := svc.UpdateProjectStatus(id, c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderProjectStatusesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/project-statuses")
}

func ProjectStatusDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid project status id")
		return
	}
	svc := referentialService()
	if err := svc.DeleteProjectStatus(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/project-statuses")
}

func TicketPriorityIndex(c *gin.Context) {
	svc := referentialService()
	rows, err := svc.GetTicketPriorities()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	renderStatusPage(c, statusPageData{
		Title:          "Ticket priorities",
		Page:           "ticketPriority",
		EntityLabel:    "Ticket priority",
		EntityPlural:   "Ticket priorities",
		Rows:           rows,
		CreateAction:   "/ticket-priorities",
		UpdateAction:   "/ticket-priorities/update",
		DeleteBasePath: "/ticket-priorities/delete/",
	})
}

func TicketPriorityStore(c *gin.Context) {
	svc := referentialService()
	if err := svc.CreateTicketPriority(c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderTicketPrioritiesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-priorities")
}

func TicketPriorityUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderTicketPrioritiesWithError(c, "ticket priority tidak valid")
		return
	}
	svc := referentialService()
	if err := svc.UpdateTicketPriority(id, c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderTicketPrioritiesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-priorities")
}

func TicketPriorityDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid ticket priority id")
		return
	}
	svc := referentialService()
	if err := svc.DeleteTicketPriority(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-priorities")
}

func TicketStatusIndex(c *gin.Context) {
	svc := referentialService()
	rows, err := svc.GetTicketStatuses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	projects, err := svc.GetProjectOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "ticket_statuses.html", gin.H{
		"Title":    "Ticket statuses",
		"Page":     "ticketStatus",
		"Rows":     rows,
		"Projects": projects,
	})
}

func TicketStatusStore(c *gin.Context) {
	order, _ := strconv.Atoi(c.PostForm("order"))
	svc := referentialService()
	if err := svc.CreateTicketStatus(c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default"), order, parseOptionalInt(c.PostForm("project_id"))); err != nil {
		renderTicketStatusesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-statuses")
}

func TicketStatusUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderTicketStatusesWithError(c, "ticket status tidak valid")
		return
	}
	order, _ := strconv.Atoi(c.PostForm("order"))
	svc := referentialService()
	if err := svc.UpdateTicketStatus(id, c.PostForm("name"), c.PostForm("color"), checkboxOn(c, "is_default"), order, parseOptionalInt(c.PostForm("project_id"))); err != nil {
		renderTicketStatusesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-statuses")
}

func TicketStatusDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid ticket status id")
		return
	}
	svc := referentialService()
	if err := svc.DeleteTicketStatus(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-statuses")
}

func TicketTypeIndex(c *gin.Context) {
	svc := referentialService()
	rows, err := svc.GetTicketTypes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "ticket_types.html", gin.H{
		"Title": "Ticket types",
		"Page":  "ticketType",
		"Rows":  rows,
	})
}

func TicketTypeStore(c *gin.Context) {
	svc := referentialService()
	if err := svc.CreateTicketType(c.PostForm("name"), c.PostForm("icon"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderTicketTypesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-types")
}

func TicketTypeUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderTicketTypesWithError(c, "ticket type tidak valid")
		return
	}
	svc := referentialService()
	if err := svc.UpdateTicketType(id, c.PostForm("name"), c.PostForm("icon"), c.PostForm("color"), checkboxOn(c, "is_default")); err != nil {
		renderTicketTypesWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-types")
}

func TicketTypeDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid ticket type id")
		return
	}
	svc := referentialService()
	if err := svc.DeleteTicketType(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-types")
}

func renderActivitiesWithError(c *gin.Context, message string) {
	svc := referentialService()
	rows, err := svc.GetActivities()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	Render(c, "activities.html", gin.H{
		"Title": "Activities",
		"Page":  "activity",
		"Rows":  rows,
		"Error": message,
	})
}

func renderProjectStatusesWithError(c *gin.Context, message string) {
	svc := referentialService()
	rows, err := svc.GetProjectStatuses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	renderStatusPage(c, statusPageData{
		Title:          "Project statuses",
		Page:           "projectStatus",
		EntityLabel:    "Project status",
		EntityPlural:   "Project statuses",
		Rows:           rows,
		Error:          message,
		CreateAction:   "/project-statuses",
		UpdateAction:   "/project-statuses/update",
		DeleteBasePath: "/project-statuses/delete/",
	})
}

func renderTicketPrioritiesWithError(c *gin.Context, message string) {
	svc := referentialService()
	rows, err := svc.GetTicketPriorities()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	renderStatusPage(c, statusPageData{
		Title:          "Ticket priorities",
		Page:           "ticketPriority",
		EntityLabel:    "Ticket priority",
		EntityPlural:   "Ticket priorities",
		Rows:           rows,
		Error:          message,
		CreateAction:   "/ticket-priorities",
		UpdateAction:   "/ticket-priorities/update",
		DeleteBasePath: "/ticket-priorities/delete/",
	})
}

func renderTicketStatusesWithError(c *gin.Context, message string) {
	svc := referentialService()
	rows, err := svc.GetTicketStatuses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	projects, err := svc.GetProjectOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	Render(c, "ticket_statuses.html", gin.H{
		"Title":    "Ticket statuses",
		"Page":     "ticketStatus",
		"Rows":     rows,
		"Projects": projects,
		"Error":    message,
	})
}

func renderTicketTypesWithError(c *gin.Context, message string) {
	svc := referentialService()
	rows, err := svc.GetTicketTypes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	Render(c, "ticket_types.html", gin.H{
		"Title": "Ticket types",
		"Page":  "ticketType",
		"Rows":  rows,
		"Error": message,
	})
}

func renderStatusPage(c *gin.Context, data statusPageData) {
	Render(c, "statuses.html", gin.H{
		"Title":          data.Title,
		"Page":           data.Page,
		"EntityLabel":    data.EntityLabel,
		"EntityPlural":   data.EntityPlural,
		"Rows":           data.Rows,
		"Error":          data.Error,
		"CreateAction":   data.CreateAction,
		"UpdateAction":   data.UpdateAction,
		"DeleteBasePath": data.DeleteBasePath,
	})
}

func checkboxOn(c *gin.Context, key string) bool {
	return strings.TrimSpace(c.PostForm(key)) != ""
}

func parseOptionalInt(value string) *int {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
}
