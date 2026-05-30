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
	renderTicketPage(c, selectedProjectIDFromQuery(c), "", "", nil, nil)
}

func TicketStore(c *gin.Context) {
	svc := managementService()
	selectedProjectID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("filter_project_id")))
	if selectedProjectID < 0 {
		selectedProjectID = 0
	}

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
		renderTicketPage(c, selectedProjectID, err.Error(), "ticketCreateModal", &input, nil)
		return
	}

	redirectURL := "/tickets"
	if selectedProjectID > 0 {
		redirectURL += "?project_id=" + strconv.Itoa(selectedProjectID)
	}
	c.Redirect(http.StatusSeeOther, redirectURL)
}

func TicketApplyTemplate(c *gin.Context) {
	svc := managementService()
	selectedProjectID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("filter_project_id")))
	if selectedProjectID < 0 {
		selectedProjectID = 0
	}

	projectID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("project_id")))
	templateSetID, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("template_set_id")))
	input := models.TicketTemplateApplyInput{
		ProjectID:     projectID,
		TemplateSetID: templateSetID,
	}

	if err := svc.ApplyTicketTemplateToProject(input, currentSessionUserID(c)); err != nil {
		renderTicketPage(c, selectedProjectID, err.Error(), "ticketTemplateApplyModal", nil, &input)
		return
	}

	redirectURL := "/tickets"
	if selectedProjectID > 0 {
		redirectURL += "?project_id=" + strconv.Itoa(selectedProjectID)
	}
	c.Redirect(http.StatusSeeOther, redirectURL)
}

func renderTicketPage(c *gin.Context, selectedProjectID int, message, openModal string, ticketOld *models.RoadmapTicketCreateInput, templateApplyOld *models.TicketTemplateApplyInput) {
	svc := managementService()

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

	templateSvc := ticketTemplateService()
	templateSetOptions, err := templateSvc.GetSets()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "tickets.html", gin.H{
		"Title":   "Tickets",
		"Page":    "ticket",
		"Tickets": tickets,
		"Columns": columns,

		"TicketError":        message,
		"OpenModal":          openModal,
		"TicketCreateOld":    ticketOld,
		"TemplateApplyOld":   templateApplyOld,
		"UserOptions":        userOptions,
		"EpicOptions":        epicOptions,
		"TemplateSetOptions": templateSetOptions,
		"ProjectOptions":     projectOptions,
		"SelectedProjectID":  selectedProjectID,
		"ProjectLabel":       resolveRoadmapProjectLabel(projectOptions, selectedProjectID),
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

	selectedProjectID := selectedProjectIDFromQuery(c)
	if selectedProjectID <= 0 {
		selectedProjectID = pageData.Ticket.ProjectID
	}

	Render(c, "ticket_detail.html", gin.H{
		"Title":             "View Ticket",
		"Page":              "ticket",
		"Ticket":            pageData.Ticket,
		"Comments":          pageData.Comments,
		"Todos":             pageData.Todos,
		"Activities":        pageData.Activities,
		"Hours":             pageData.Hours,
		"Subscribers":       pageData.Subscribers,
		"Attachments":       pageData.Attachments,
		"CurrentUserID":     currentSessionUserID(c),
		"SelectedProjectID": selectedProjectID,
		"BackToTicketsURL":  ticketListURL(selectedProjectID),
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

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
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

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
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

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
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

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
}

func TicketTodoStore(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	if err := managementService().CreateTicketTodo(ticketID, currentSessionUserID(c), c.PostForm("content")); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
}

func TicketTodoUpdate(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	todoID, err := strconv.Atoi(c.Param("todoId"))
	if err != nil || todoID <= 0 {
		c.String(http.StatusBadRequest, "todo tidak valid")
		return
	}

	isDone := strings.TrimSpace(c.PostForm("is_done")) != ""
	if err := managementService().UpdateTicketTodo(ticketID, todoID, currentSessionUserID(c), c.PostForm("content"), isDone); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
}

func TicketTodoDelete(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ticketID <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	todoID, err := strconv.Atoi(c.Param("todoId"))
	if err != nil || todoID <= 0 {
		c.String(http.StatusBadRequest, "todo tidak valid")
		return
	}

	if err := managementService().DeleteTicketTodo(ticketID, todoID, currentSessionUserID(c)); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, ticketDetailURL(ticketID, selectedProjectIDFromQuery(c)))
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

func TicketDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "ticket tidak valid")
		return
	}

	if err := managementService().DeleteTicket(id); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, ticketListURL(selectedProjectIDFromQuery(c)))
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

type roadmapProjectPeriodViewRow struct {
	ProjectName    string
	StartDateLabel string
	EndDateLabel   string
	DurationLabel  string
	ShowBar        bool
	BarLeftPx      int
	BarWidthPx     int
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
	projectPeriods, err := svc.GetRoadmapProjectPeriods()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalProjects := len(projectOptions)
	filteredEpics, filteredTickets := filterRoadmapByProject(epics, tickets, selectedProjectID)
	projectPeriodWeeks, projectPeriodYearGroups, projectPeriodRows, projectPeriodTimelineWidth, projectPeriodColumnWidth, projectPeriodCurrentMarkerLeft, projectPeriodCurrentMarkerWidth, projectPeriodRangeLabel := buildRoadmapProjectPeriodTimeline(projectPeriods, roadmapNow(), format)
	weeks, rows, timelineWidth, currentMarkerLeft, currentMarkerWidth, columnWidth := svc.BuildRoadmapTimeline(filteredEpics, filteredTickets, roadmapNow(), format)
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
		"Title":                           "Road Map",
		"Page":                            "roadmap",
		"Format":                          format,
		"Epics":                           filteredEpics,
		"Tickets":                         filteredTickets,
		"Weeks":                           weeks,
		"YearGroups":                      yearGroups,
		"Rows":                            rows,
		"TimelineWidth":                   timelineWidth,
		"CurrentMarkerLeft":               currentMarkerLeft,
		"CurrentMarkerWidth":              currentMarkerWidth,
		"ColumnWidth":                     columnWidth,
		"TotalProjects":                   totalProjects,
		"ProjectLabel":                    projectLabel,
		"SelectedProjectID":               selectedProjectID,
		"ProjectPeriodWeeks":              projectPeriodWeeks,
		"ProjectPeriodYearGroups":         projectPeriodYearGroups,
		"ProjectPeriodRows":               projectPeriodRows,
		"ProjectPeriodTimelineWidth":      projectPeriodTimelineWidth,
		"ProjectPeriodColumnWidth":        projectPeriodColumnWidth,
		"ProjectPeriodCurrentMarkerLeft":  projectPeriodCurrentMarkerLeft,
		"ProjectPeriodCurrentMarkerWidth": projectPeriodCurrentMarkerWidth,
		"ProjectPeriodRange":              projectPeriodRangeLabel,
		"RoadmapError":                    message,
		"OpenModal":                       openModal,
		"EpicOld":                         epicOld,
		"TicketOld":                       ticketOld,
		"ProjectOptions":                  projectOptions,
		"EpicOptions":                     epicOptions,
		"UserOptions":                     userOptions,
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

func buildRoadmapProjectPeriodTimeline(periods []models.RoadmapProjectPeriod, now time.Time, format string) ([]models.RoadmapWeek, []models.RoadmapYearGroup, []roadmapProjectPeriodViewRow, int, int, int, int, string) {
	location := now.Location()
	if location == nil {
		location = time.UTC
	}

	type parsedPeriod struct {
		models.RoadmapProjectPeriod
		start time.Time
		end   time.Time
	}

	parsed := make([]parsedPeriod, 0, len(periods))
	var minStart, maxEnd time.Time

	for _, period := range periods {
		start, errStart := time.ParseInLocation("2006-01-02", strings.TrimSpace(period.StartsAtISO), location)
		end, errEnd := time.ParseInLocation("2006-01-02", strings.TrimSpace(period.EndsAtISO), location)
		if errStart != nil || errEnd != nil || end.Before(start) {
			continue
		}

		parsed = append(parsed, parsedPeriod{
			RoadmapProjectPeriod: period,
			start:                start,
			end:                  end,
		})
		if minStart.IsZero() || start.Before(minStart) {
			minStart = start
		}
		if maxEnd.IsZero() || end.After(maxEnd) {
			maxEnd = end
		}
	}

	if minStart.IsZero() || maxEnd.IsZero() {
		start := roadmapProjectStartOfWeek(now)
		minStart = start
		maxEnd = start.AddDate(0, 0, 83)
	}

	baseStart := minStart.AddDate(0, 0, -7)
	baseEnd := maxEnd.AddDate(0, 0, 28)
	var rangeStart, rangeEnd time.Time
	switch format {
	case "day":
		rangeStart = roadmapProjectStartOfDay(baseStart)
		rangeEnd = roadmapProjectStartOfDay(baseEnd)
	case "month":
		rangeStart = roadmapProjectFirstOfMonth(baseStart)
		rangeEnd = roadmapProjectEndOfMonth(baseEnd)
	default:
		rangeStart = roadmapProjectStartOfWeek(baseStart)
		rangeEnd = roadmapProjectEndOfWeek(baseEnd)
	}

	columns, columnWidth := buildRoadmapProjectColumns(rangeStart, rangeEnd, format)
	timelineWidth := len(columns) * columnWidth
	yearGroups := buildRoadmapYearGroups(columns, columnWidth)

	rows := make([]roadmapProjectPeriodViewRow, 0, len(parsed))
	for _, item := range parsed {
		barLeftPx, barWidthPx := roadmapProjectBarMetrics(item.start, item.end, rangeStart, format, columnWidth)
		durationDays := roadmapProjectDaysBetween(item.start, item.end) + 1
		if durationDays < 1 {
			durationDays = 1
		}

		rows = append(rows, roadmapProjectPeriodViewRow{
			ProjectName:    item.ProjectName,
			StartDateLabel: item.StartsAt,
			EndDateLabel:   item.EndsAt,
			DurationLabel:  fmt.Sprintf("%d hari", durationDays),
			ShowBar:        barWidthPx > 0,
			BarLeftPx:      barLeftPx,
			BarWidthPx:     barWidthPx,
		})
	}

	rangeLabel := fmt.Sprintf("%s - %s", minStart.Format("02 Jan 2006"), maxEnd.Format("02 Jan 2006"))
	currentMarkerLeft, currentMarkerWidth := roadmapProjectCurrentMarkerMetrics(now, rangeStart, format, columnWidth)
	return columns, yearGroups, rows, timelineWidth, columnWidth, currentMarkerLeft, currentMarkerWidth, rangeLabel
}

func buildRoadmapProjectColumns(start, end time.Time, format string) ([]models.RoadmapWeek, int) {
	switch format {
	case "day":
		var columns []models.RoadmapWeek
		for cursor := start; !cursor.After(end); cursor = cursor.AddDate(0, 0, 1) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("02 Jan 06"),
			})
		}
		return columns, 78
	case "month":
		var columns []models.RoadmapWeek
		for cursor := roadmapProjectFirstOfMonth(start); !cursor.After(end); cursor = cursor.AddDate(0, 1, 0) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("Jan 06"),
			})
		}
		return columns, 102
	default:
		var columns []models.RoadmapWeek
		for cursor := start; !cursor.After(end); cursor = cursor.AddDate(0, 0, 7) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("02 Jan 06"),
			})
		}
		return columns, 72
	}
}

func roadmapProjectBarMetrics(start, end, rangeStart time.Time, format string, columnWidth int) (int, int) {
	switch format {
	case "day":
		offset := roadmapProjectDaysBetween(rangeStart, roadmapProjectStartOfDay(start))
		span := roadmapProjectDaysBetween(roadmapProjectStartOfDay(start), roadmapProjectStartOfDay(end)) + 1
		return offset * columnWidth, roadmapProjectMax(span, 1) * columnWidth
	case "month":
		offset := roadmapProjectMonthsBetween(roadmapProjectFirstOfMonth(rangeStart), roadmapProjectFirstOfMonth(start))
		span := roadmapProjectMonthsBetween(roadmapProjectFirstOfMonth(start), roadmapProjectFirstOfMonth(end)) + 1
		return offset * columnWidth, roadmapProjectMax(span, 1) * columnWidth
	default:
		offset := roadmapProjectDaysBetween(roadmapProjectStartOfWeek(rangeStart), roadmapProjectStartOfWeek(start)) / 7
		span := (roadmapProjectDaysBetween(roadmapProjectStartOfWeek(start), roadmapProjectStartOfWeek(end)) / 7) + 1
		return offset * columnWidth, roadmapProjectMax(span, 1) * columnWidth
	}
}

func roadmapProjectCurrentMarkerMetrics(now, rangeStart time.Time, format string, columnWidth int) (int, int) {
	switch format {
	case "day":
		return roadmapProjectDaysBetween(rangeStart, roadmapProjectStartOfDay(now)) * columnWidth, columnWidth
	case "month":
		return roadmapProjectMonthsBetween(roadmapProjectFirstOfMonth(rangeStart), roadmapProjectFirstOfMonth(now)) * columnWidth, columnWidth
	default:
		return (roadmapProjectDaysBetween(roadmapProjectStartOfWeek(rangeStart), roadmapProjectStartOfWeek(now)) / 7) * columnWidth, columnWidth
	}
}

func roadmapProjectStartOfWeek(value time.Time) time.Time {
	normalized := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
	weekday := int(normalized.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return normalized.AddDate(0, 0, -(weekday - 1))
}

func roadmapProjectEndOfWeek(value time.Time) time.Time {
	return roadmapProjectStartOfWeek(value).AddDate(0, 0, 6)
}

func roadmapProjectStartOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func roadmapProjectFirstOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func roadmapProjectEndOfMonth(value time.Time) time.Time {
	return roadmapProjectFirstOfMonth(value).AddDate(0, 1, -1)
}

func roadmapProjectDaysBetween(start, end time.Time) int {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	return int(endDate.Sub(startDate).Hours() / 24)
}

func roadmapProjectMonthsBetween(start, end time.Time) int {
	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month())
}

func roadmapProjectMax(a, b int) int {
	if a > b {
		return a
	}
	return b
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

func selectedProjectIDFromQuery(c *gin.Context) int {
	selectedProjectID, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("project_id", "0")))
	if selectedProjectID < 0 {
		return 0
	}
	return selectedProjectID
}

func ticketListURL(selectedProjectID int) string {
	redirectURL := "/tickets"
	if selectedProjectID > 0 {
		redirectURL += "?project_id=" + strconv.Itoa(selectedProjectID)
	}
	return redirectURL
}

func ticketDetailURL(ticketID, selectedProjectID int) string {
	redirectURL := "/tickets/" + strconv.Itoa(ticketID)
	if selectedProjectID > 0 {
		redirectURL += "?project_id=" + strconv.Itoa(selectedProjectID)
	}
	return redirectURL
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
