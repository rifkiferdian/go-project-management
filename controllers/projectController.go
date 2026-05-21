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
		DeveloperID  int    `form:"developer_id" binding:"required"`
		StartDate    string `form:"start_date"`
		EndDate      string `form:"end_date"`
		StatusID     int    `form:"status_id" binding:"required"`
		PriorityID   int    `form:"priority_id" binding:"required"`
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

	divisionIDs, err := parseIDList(c.PostFormArray("request_divisions"))
	if err != nil {
		renderProjectPage(c, projectSvc, models.Project{}, "divisi requester tidak valid")
		return
	}

	input := models.ProjectCreateInput{
		Name:         strings.TrimSpace(form.Name),
		Description:  strings.TrimSpace(form.Description),
		OwnerID:      form.OwnerID,
		DeveloperID:  form.DeveloperID,
		StartDate:    strings.TrimSpace(form.StartDate),
		EndDate:      strings.TrimSpace(form.EndDate),
		DivisionIDs:  divisionIDs,
		StatusID:     form.StatusID,
		PriorityID:   form.PriorityID,
		TicketPrefix: strings.TrimSpace(form.TicketPrefix),
		StatusType:   form.StatusType,
		Type:         form.Type,
	}

	if err := projectSvc.CreateProject(input); err != nil {
		renderProjectPage(c, projectSvc, models.Project{
			Name:               input.Name,
			Description:        input.Description,
			OwnerID:            input.OwnerID,
			DeveloperID:        input.DeveloperID,
			StartDate:          input.StartDate,
			EndDate:            input.EndDate,
			RequestDivisionIDs: ints64ToInts(input.DivisionIDs),
			StatusID:           input.StatusID,
			PriorityID:         input.PriorityID,
			TicketPrefix:       input.TicketPrefix,
			StatusType:         input.StatusType,
			Type:               input.Type,
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
		DeveloperID  int    `form:"developer_id" binding:"required"`
		StartDate    string `form:"start_date"`
		EndDate      string `form:"end_date"`
		StatusID     int    `form:"status_id" binding:"required"`
		PriorityID   int    `form:"priority_id" binding:"required"`
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

	divisionIDs, err := parseIDList(c.PostFormArray("request_divisions"))
	if err != nil {
		renderProjectPage(c, projectSvc, models.Project{ID: form.ID}, "divisi requester tidak valid")
		return
	}

	input := models.ProjectUpdateInput{
		ID:           form.ID,
		Name:         strings.TrimSpace(form.Name),
		Description:  strings.TrimSpace(form.Description),
		OwnerID:      form.OwnerID,
		DeveloperID:  form.DeveloperID,
		StartDate:    strings.TrimSpace(form.StartDate),
		EndDate:      strings.TrimSpace(form.EndDate),
		DivisionIDs:  divisionIDs,
		StatusID:     form.StatusID,
		PriorityID:   form.PriorityID,
		TicketPrefix: strings.TrimSpace(form.TicketPrefix),
		StatusType:   form.StatusType,
		Type:         form.Type,
	}

	if err := projectSvc.UpdateProject(input); err != nil {
		renderProjectPage(c, projectSvc, models.Project{
			ID:                 input.ID,
			Name:               input.Name,
			Description:        input.Description,
			OwnerID:            input.OwnerID,
			DeveloperID:        input.DeveloperID,
			StartDate:          input.StartDate,
			EndDate:            input.EndDate,
			RequestDivisionIDs: ints64ToInts(input.DivisionIDs),
			StatusID:           input.StatusID,
			PriorityID:         input.PriorityID,
			TicketPrefix:       input.TicketPrefix,
			StatusType:         input.StatusType,
			Type:               input.Type,
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
	divisions, err := projectService.GetDivisionOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	priorities, err := projectService.GetPriorityOptions()
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
	developerUsers, err := userRepo.GetByDivisionName("Audit & Sistem (IT)")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "project.html", gin.H{
		"Title":          "Daftar Project",
		"Page":           "project",
		"projects":       projects,
		"statuses":       statuses,
		"users":          users,
		"developerUsers": developerUsers,
		"divisions":      divisions,
		"priorities":     priorities,
		"Old":            old,
		"Error":          message,
	})
}

func ints64ToInts(values []int64) []int {
	result := make([]int, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		result = append(result, int(value))
	}
	return result
}
