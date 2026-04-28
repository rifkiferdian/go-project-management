package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"gobase-app/config"
	"gobase-app/models"
	"gobase-app/repositories"
	"gobase-app/services"

	"github.com/gin-gonic/gin"
)

func ProjectIndex(c *gin.Context) {
	projectRepo := &repositories.ProjectRepository{DB: config.DB}
	projectService := &services.ProjectService{Repo: projectRepo}

	renderProjectPage(c, projectService, models.Project{}, "")
}

func ProjectStore(c *gin.Context) {
	type projectForm struct {
		Name         string `form:"name" binding:"required"`
		Description  string `form:"description"`
		OwnerID      int    `form:"owner_id" binding:"required"`
		StatusID     int    `form:"status_id" binding:"required"`
		TicketPrefix string `form:"ticket_prefix" binding:"required"`
		StatusType   string `form:"status_type"`
		Type         string `form:"type"`
	}

	var (
		form        projectForm
		projectRepo = &repositories.ProjectRepository{DB: config.DB}
		projectSvc  = &services.ProjectService{Repo: projectRepo}
	)

	if err := c.ShouldBind(&form); err != nil {
		renderProjectPage(c, projectSvc, models.Project{}, "Form tidak lengkap")
		return
	}

	input := models.ProjectCreateInput{
		Name:         strings.TrimSpace(form.Name),
		Description:  strings.TrimSpace(form.Description),
		OwnerID:      form.OwnerID,
		StatusID:     form.StatusID,
		TicketPrefix: strings.TrimSpace(form.TicketPrefix),
		StatusType:   form.StatusType,
		Type:         form.Type,
	}

	if err := projectSvc.CreateProject(input); err != nil {
		renderProjectPage(c, projectSvc, models.Project{
			Name:         input.Name,
			Description:  input.Description,
			OwnerID:      input.OwnerID,
			StatusID:     input.StatusID,
			TicketPrefix: input.TicketPrefix,
			StatusType:   input.StatusType,
			Type:         input.Type,
		}, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/projects")
}

func ProjectUpdate(c *gin.Context) {
	type projectUpdateForm struct {
		ID           int    `form:"project_id" binding:"required"`
		Name         string `form:"name" binding:"required"`
		Description  string `form:"description"`
		OwnerID      int    `form:"owner_id" binding:"required"`
		StatusID     int    `form:"status_id" binding:"required"`
		TicketPrefix string `form:"ticket_prefix" binding:"required"`
		StatusType   string `form:"status_type"`
		Type         string `form:"type"`
	}

	var (
		form        projectUpdateForm
		projectRepo = &repositories.ProjectRepository{DB: config.DB}
		projectSvc  = &services.ProjectService{Repo: projectRepo}
	)

	if err := c.ShouldBind(&form); err != nil {
		renderProjectPage(c, projectSvc, models.Project{ID: form.ID}, "Form tidak lengkap")
		return
	}

	input := models.ProjectUpdateInput{
		ID:           form.ID,
		Name:         strings.TrimSpace(form.Name),
		Description:  strings.TrimSpace(form.Description),
		OwnerID:      form.OwnerID,
		StatusID:     form.StatusID,
		TicketPrefix: strings.TrimSpace(form.TicketPrefix),
		StatusType:   form.StatusType,
		Type:         form.Type,
	}

	if err := projectSvc.UpdateProject(input); err != nil {
		renderProjectPage(c, projectSvc, models.Project{
			ID:           input.ID,
			Name:         input.Name,
			Description:  input.Description,
			OwnerID:      input.OwnerID,
			StatusID:     input.StatusID,
			TicketPrefix: input.TicketPrefix,
			StatusType:   input.StatusType,
			Type:         input.Type,
		}, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/projects")
}

func ProjectDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid project id")
		return
	}

	projectRepo := &repositories.ProjectRepository{DB: config.DB}
	projectService := &services.ProjectService{Repo: projectRepo}

	if err := projectService.DeleteProject(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/projects")
}

func renderProjectPage(c *gin.Context, projectService *services.ProjectService, old models.Project, message string) {
	projects, err := projectService.GetProjects()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	statuses, err := projectService.GetStatusOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userRepo := &repositories.UserRepository{DB: config.DB}
	users, err := userRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "project.html", gin.H{
		"Title":    "Daftar Project",
		"Page":     "project",
		"projects": projects,
		"statuses": statuses,
		"users":    users,
		"Old":      old,
		"Error":    message,
	})
}
