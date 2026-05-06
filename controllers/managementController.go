package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gobase-app/config"
	"gobase-app/models"
	"gobase-app/repositories"
	"gobase-app/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TicketIndex(c *gin.Context) {
	svc := managementService()
	selectedProjectID, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("project_id", "0")))
	if selectedProjectID < 0 {
		selectedProjectID = 0
	}

	tickets, err := svc.GetTickets(selectedProjectID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	columns, err := svc.GetBoardColumns(selectedProjectID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	projectOptions, err := svc.GetRoadmapProjectOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "tickets.html", gin.H{
		"Title":   "Tickets",
		"Page":    "ticket",
		"Tickets": tickets,
		"Columns": columns,

		"ProjectOptions":    projectOptions,
		"SelectedProjectID": selectedProjectID,
		"ProjectLabel":      resolveRoadmapProjectLabel(projectOptions, selectedProjectID),
	})
}

func TicketShow(c *gin.Context) {
	svc := managementService()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code_error": http.StatusBadRequest,
			"error":      "ticket tidak valid",
		})
		return
	}

	pageData, err := svc.GetTicketDetailPage(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "ticket tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.HTML(statusCode, "error.html", gin.H{
			"code_error": statusCode,
			"error":      err.Error(),
		})
		return
	}

	Render(c, "ticket_detail.html", gin.H{
		"Title":         "View Ticket",
		"Page":          "ticket",
		"Ticket":        pageData.Ticket,
		"Comments":      pageData.Comments,
		"Activities":    pageData.Activities,
		"Hours":         pageData.Hours,
		"Subscribers":   pageData.Subscribers,
		"Attachments":   pageData.Attachments,
		"CurrentUserID": currentSessionUserID(c),
	})
}

func TicketAttachmentStore(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		c.String(http.StatusBadRequest, "file attachment wajib dipilih")
		return
	}
	if file.Size <= 0 {
		c.String(http.StatusBadRequest, "file attachment kosong")
		return
	}
	if file.Size > 10*1024*1024 {
		c.String(http.StatusBadRequest, "file attachment maksimal 10 MB")
		return
	}

	uploadDir := filepath.Join("assets", "uploads", "tickets", strconv.Itoa(ticketID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	safeName := safeUploadFilename(file.Filename)
	storedName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), safeName)
	destination := filepath.Join(uploadDir, storedName)
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	publicPath := "/assets/uploads/tickets/" + strconv.Itoa(ticketID) + "/" + storedName
	input := models.TicketAttachmentCreateInput{
		TicketID:     ticketID,
		UserID:       currentSessionUserID(c),
		OriginalName: filepath.Base(file.Filename),
		FileName:     storedName,
		FilePath:     publicPath,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
	}

	if err := managementService().CreateTicketAttachment(input); err != nil {
		_ = os.Remove(destination)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/tickets/"+strconv.Itoa(ticketID))
}

func TicketContentUpdate(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	if err := managementService().UpdateTicketContent(ticketID, c.PostForm("content")); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/tickets/"+strconv.Itoa(ticketID))
}

func TicketCommentStore(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	if err := managementService().CreateTicketComment(ticketID, currentSessionUserID(c), c.PostForm("content")); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/tickets/"+strconv.Itoa(ticketID))
}

func TicketCommentUpdate(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil || commentID <= 0 {
		c.String(http.StatusBadRequest, "comment tidak valid")
		return
	}

	if err := managementService().UpdateTicketComment(ticketID, commentID, currentSessionUserID(c), c.PostForm("content")); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/tickets/"+strconv.Itoa(ticketID))
}

func TicketEdit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code_error": http.StatusBadRequest,
			"error":      "ticket tidak valid",
		})
		return
	}

	renderTicketEditPage(c, id, nil, "")
}

func TicketUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("ticket_id")))
	statusID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("status_id")))
	priorityID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("priority_id")))
	typeID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("type_id")))
	ownerID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("owner_id")))
	responsibleID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("responsible_id")))
	epicID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("epic_id")))

	input := models.TicketUpdateInput{
		ID:            id,
		Name:          c.PostForm("name"),
		Content:       c.PostForm("content"),
		StatusID:      statusID,
		PriorityID:    priorityID,
		TypeID:        typeID,
		OwnerID:       ownerID,
		ResponsibleID: responsibleID,
		EpicID:        epicID,
		Estimation:    c.PostForm("estimation"),
		StartsAt:      c.PostForm("starts_at"),
		EndsAt:        c.PostForm("ends_at"),
	}

	svc := managementService()
	savedInput, err := svc.UpdateTicket(input, currentSessionUserID(c))
	if err != nil {
		renderTicketEditPage(c, input.ID, &savedInput, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/tickets/"+strconv.Itoa(input.ID))
}

func BoardIndex(c *gin.Context) {
	svc := managementService()

	columns, err := svc.GetBoardColumns(0)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalTickets := 0
	for _, column := range columns {
		totalTickets += column.TicketCount
	}

	Render(c, "board.html", gin.H{
		"Title":        "Board",
		"Page":         "board",
		"Columns":      columns,
		"TotalStatus":  len(columns),
		"TotalTickets": totalTickets,
	})
}

func RoadMapIndex(c *gin.Context) {
	renderRoadMapPage(c, "", "", nil, nil)
}

func RoadMapEpicStore(c *gin.Context) {
	svc := managementService()
	projectID, _ := strconv.Atoi(c.PostForm("project_id"))
	input := models.RoadmapEpicCreateInput{
		ProjectID: projectID,
		Name:      c.PostForm("name"),
		StartsAt:  c.PostForm("starts_at"),
		EndsAt:    c.PostForm("ends_at"),
	}

	if err := svc.CreateRoadmapEpic(input); err != nil {
		renderRoadMapPage(c, err.Error(), "epicModal", input, nil)
		return
	}

	c.Redirect(http.StatusSeeOther, "/road-map?format="+normalizeRoadmapFormat(c.DefaultPostForm("format", "week")))
}

func RoadMapTicketStore(c *gin.Context) {
	svc := managementService()
	projectID, _ := strconv.Atoi(c.PostForm("project_id"))
	resourceUserID, _ := strconv.Atoi(c.PostForm("resource_user_id"))
	estimation, _ := strconv.ParseFloat(strings.TrimSpace(c.PostForm("estimation")), 64)
	input := models.RoadmapTicketCreateInput{
		ProjectID:      projectID,
		EpicID:         parseOptionalInt(c.PostForm("epic_id")),
		Name:           c.PostForm("name"),
		ResourceUserID: resourceUserID,
		Estimation:     estimation,
		StartsAt:       c.PostForm("starts_at"),
		EndsAt:         c.PostForm("ends_at"),
	}

	if err := svc.CreateRoadmapTicket(input); err != nil {
		renderRoadMapPage(c, err.Error(), "ticketModal", nil, input)
		return
	}

	c.Redirect(http.StatusSeeOther, "/road-map?format="+normalizeRoadmapFormat(c.DefaultPostForm("format", "week")))
}

func renderRoadMapPage(c *gin.Context, message, openModal string, epicOld interface{}, ticketOld interface{}) {
	svc := managementService()
	format := normalizeRoadmapFormat(c.DefaultQuery("format", "week"))
	selectedProjectID, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("project_id", "0")))

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

	projectOptions, err := svc.GetRoadmapProjectOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalProjects := len(projectOptions)
	filteredEpics, filteredTickets := filterRoadmapByProject(epics, tickets, selectedProjectID)
	weeks, rows, timelineWidth, currentMarkerLeft, currentMarkerWidth, columnWidth := svc.BuildRoadmapTimeline(filteredEpics, filteredTickets, time.Now(), format)
	yearGroups := buildRoadmapYearGroups(weeks, columnWidth)
	projectLabel := resolveRoadmapProjectLabel(projectOptions, selectedProjectID)

	epicOptions, err := svc.GetRoadmapEpicOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userRepo := &repositories.UserRepository{DB: config.DB}
	userOptions, err := userRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "roadmap.html", gin.H{
		"Title":              "Road Map",
		"Page":               "roadmap",
		"Format":             format,
		"Epics":              filteredEpics,
		"Tickets":            filteredTickets,
		"Weeks":              weeks,
		"YearGroups":         yearGroups,
		"Rows":               rows,
		"TimelineWidth":      timelineWidth,
		"CurrentMarkerLeft":  currentMarkerLeft,
		"CurrentMarkerWidth": currentMarkerWidth,
		"ColumnWidth":        columnWidth,
		"TotalProjects":      totalProjects,
		"ProjectLabel":       projectLabel,
		"SelectedProjectID":  selectedProjectID,
		"RoadmapError":       message,
		"OpenModal":          openModal,
		"EpicOld":            epicOld,
		"TicketOld":          ticketOld,
		"ProjectOptions":     projectOptions,
		"EpicOptions":        epicOptions,
		"UserOptions":        userOptions,
	})
}

func filterRoadmapByProject(epics []models.RoadmapEpic, tickets []models.RoadmapTicket, projectID int) ([]models.RoadmapEpic, []models.RoadmapTicket) {
	if projectID <= 0 {
		return []models.RoadmapEpic{}, []models.RoadmapTicket{}
	}

	filteredEpics := make([]models.RoadmapEpic, 0, len(epics))
	allowedEpicIDs := make(map[int]struct{})
	for _, epic := range epics {
		if epic.ProjectID == projectID {
			filteredEpics = append(filteredEpics, epic)
			allowedEpicIDs[epic.ID] = struct{}{}
		}
	}

	filteredTickets := make([]models.RoadmapTicket, 0, len(tickets))
	for _, ticket := range tickets {
		if ticket.ProjectID != projectID {
			continue
		}
		if ticket.EpicID == 0 {
			filteredTickets = append(filteredTickets, ticket)
			continue
		}
		if _, ok := allowedEpicIDs[ticket.EpicID]; ok {
			filteredTickets = append(filteredTickets, ticket)
		}
	}

	return filteredEpics, filteredTickets
}

func resolveRoadmapProjectLabel(options []models.ProjectOption, selectedProjectID int) string {
	if selectedProjectID <= 0 {
		return ""
	}

	for _, option := range options {
		if option.ID == selectedProjectID {
			return option.Name
		}
	}
	return ""
}

func buildRoadmapYearGroups(weeks []models.RoadmapWeek, columnWidth int) []models.RoadmapYearGroup {
	if len(weeks) == 0 {
		return nil
	}

	groups := make([]models.RoadmapYearGroup, 0, len(weeks))
	current := models.RoadmapYearGroup{Label: weeks[0].YearLabel, Count: 0}

	for _, week := range weeks {
		if week.YearLabel != current.Label {
			current.WidthPx = current.Count * columnWidth
			groups = append(groups, current)
			current = models.RoadmapYearGroup{Label: week.YearLabel, Count: 0}
		}
		current.Count++
	}

	current.WidthPx = current.Count * columnWidth
	groups = append(groups, current)
	return groups
}

func managementService() *services.ManagementService {
	return &services.ManagementService{
		Repo: &repositories.ManagementRepository{DB: config.DB},
	}
}

func renderTicketEditPage(c *gin.Context, id int, old *models.TicketUpdateInput, message string) {
	svc := managementService()
	pageData, err := svc.GetTicketEditPage(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "ticket tidak ditemukan" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "ticket tidak valid" {
			statusCode = http.StatusBadRequest
		}
		c.HTML(statusCode, "error.html", gin.H{
			"code_error": statusCode,
			"error":      err.Error(),
		})
		return
	}

	if old != nil {
		pageData.Form.ID = old.ID
		pageData.Form.Name = old.Name
		pageData.Form.Content = old.Content
		pageData.Form.StatusID = old.StatusID
		pageData.Form.PriorityID = old.PriorityID
		pageData.Form.TypeID = old.TypeID
		pageData.Form.OwnerID = old.OwnerID
		pageData.Form.ResponsibleID = old.ResponsibleID
		pageData.Form.EpicID = old.EpicID
		pageData.Form.Estimation = old.Estimation
		pageData.Form.StartsAt = old.StartsAt
		pageData.Form.EndsAt = old.EndsAt
	}

	Render(c, "ticket_edit.html", gin.H{
		"Title": "Edit Ticket",
		"Page":  "ticket",
		"Form":  pageData.Form,
		"Error": message,

		"StatusOptions":   pageData.StatusOptions,
		"PriorityOptions": pageData.PriorityOptions,
		"TypeOptions":     pageData.TypeOptions,
		"UserOptions":     pageData.UserOptions,
		"EpicOptions":     pageData.EpicOptions,
	})
}

func normalizeRoadmapFormat(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "day":
		return "day"
	case "month":
		return "month"
	default:
		return "week"
	}
}

func currentSessionUserID(c *gin.Context) int {
	session := sessions.Default(c)

	switch value := session.Get("user_id").(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	case string:
		id, _ := strconv.Atoi(strings.TrimSpace(value))
		return id
	}

	switch value := session.Get("user").(type) {
	case models.SessionUser:
		return value.UserID
	case map[string]interface{}:
		switch raw := value["user_id"].(type) {
		case int:
			return raw
		case int64:
			return int(raw)
		case float64:
			return int(raw)
		case string:
			id, _ := strconv.Atoi(strings.TrimSpace(raw))
			return id
		}
	}

	return 0
}

func safeUploadFilename(name string) string {
	base := strings.TrimSpace(filepath.Base(name))
	if base == "" || base == "." {
		return "attachment"
	}

	extension := filepath.Ext(base)
	stem := strings.TrimSuffix(base, extension)
	var builder strings.Builder
	for _, r := range stem {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '-' || r == '_':
			builder.WriteRune(r)
		default:
			builder.WriteRune('-')
		}
	}

	safeStem := strings.Trim(builder.String(), "-_")
	if safeStem == "" {
		safeStem = "attachment"
	}
	return safeStem + strings.ToLower(extension)
}
