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

func ticketTemplateService() *services.TicketTemplateService {
	return &services.TicketTemplateService{
		Repo: &repositories.TicketTemplateRepository{DB: config.DB},
	}
}

func TicketTemplateIndex(c *gin.Context) {
	data, err := ticketTemplatePageData("")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	Render(c, "ticket_templates.html", data)
}

func TicketTemplateSetStore(c *gin.Context) {
	svc := ticketTemplateService()
	if err := svc.CreateSet(
		c.PostForm("name"),
		c.PostForm("purpose"),
		c.PostForm("description"),
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateSetUpdate(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		renderTicketTemplateWithError(c, "set template tidak valid")
		return
	}

	svc := ticketTemplateService()
	if err := svc.UpdateSet(
		id,
		c.PostForm("name"),
		c.PostForm("purpose"),
		c.PostForm("description"),
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateSetDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid template set id")
		return
	}

	svc := ticketTemplateService()
	if err := svc.DeleteSet(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateEpicStore(c *gin.Context) {
	setID, err := strconv.Atoi(strings.TrimSpace(c.PostForm("set_id")))
	if err != nil {
		renderTicketTemplateWithError(c, "set template wajib dipilih")
		return
	}
	sortOrder, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("sort_order")))

	svc := ticketTemplateService()
	if err := svc.CreateEpic(
		setID,
		c.PostForm("name"),
		c.PostForm("description"),
		parseOptionalInt(c.PostForm("start_offset_days")),
		parseOptionalInt(c.PostForm("due_offset_days")),
		sortOrder,
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateEpicUpdate(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		renderTicketTemplateWithError(c, "epic template tidak valid")
		return
	}
	setID, err := strconv.Atoi(strings.TrimSpace(c.PostForm("set_id")))
	if err != nil {
		renderTicketTemplateWithError(c, "set template wajib dipilih")
		return
	}
	sortOrder, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("sort_order")))

	svc := ticketTemplateService()
	if err := svc.UpdateEpic(
		id,
		setID,
		c.PostForm("name"),
		c.PostForm("description"),
		parseOptionalInt(c.PostForm("start_offset_days")),
		parseOptionalInt(c.PostForm("due_offset_days")),
		sortOrder,
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateEpicDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid template epic id")
		return
	}

	svc := ticketTemplateService()
	if err := svc.DeleteEpic(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateItemStore(c *gin.Context) {
	setID, err := strconv.Atoi(strings.TrimSpace(c.PostForm("set_id")))
	if err != nil {
		renderTicketTemplateWithError(c, "set template wajib dipilih")
		return
	}

	sortOrder, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("sort_order")))
	svc := ticketTemplateService()
	if err := svc.CreateItem(
		setID,
		c.PostForm("title"),
		c.PostForm("description"),
		parseOptionalInt(c.PostForm("template_epic_id")),
		parseOptionalInt(c.PostForm("default_type_id")),
		parseOptionalInt(c.PostForm("default_priority_id")),
		parseOptionalInt(c.PostForm("default_status_id")),
		parseOptionalInt(c.PostForm("default_owner_id")),
		parseOptionalInt(c.PostForm("default_responsible_id")),
		parseOptionalFloat(c.PostForm("estimation")),
		parseOptionalInt(c.PostForm("start_offset_days")),
		parseOptionalInt(c.PostForm("due_offset_days")),
		sortOrder,
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateItemUpdate(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		renderTicketTemplateWithError(c, "item template tidak valid")
		return
	}

	setID, err := strconv.Atoi(strings.TrimSpace(c.PostForm("set_id")))
	if err != nil {
		renderTicketTemplateWithError(c, "set template wajib dipilih")
		return
	}

	sortOrder, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("sort_order")))
	svc := ticketTemplateService()
	if err := svc.UpdateItem(
		id,
		setID,
		c.PostForm("title"),
		c.PostForm("description"),
		parseOptionalInt(c.PostForm("template_epic_id")),
		parseOptionalInt(c.PostForm("default_type_id")),
		parseOptionalInt(c.PostForm("default_priority_id")),
		parseOptionalInt(c.PostForm("default_status_id")),
		parseOptionalInt(c.PostForm("default_owner_id")),
		parseOptionalInt(c.PostForm("default_responsible_id")),
		parseOptionalFloat(c.PostForm("estimation")),
		parseOptionalInt(c.PostForm("start_offset_days")),
		parseOptionalInt(c.PostForm("due_offset_days")),
		sortOrder,
		checkboxOn(c, "is_active"),
	); err != nil {
		renderTicketTemplateWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func TicketTemplateItemDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid template item id")
		return
	}

	svc := ticketTemplateService()
	if err := svc.DeleteItem(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/ticket-templates")
}

func renderTicketTemplateWithError(c *gin.Context, message string) {
	data, err := ticketTemplatePageData(message)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	Render(c, "ticket_templates.html", data)
}

func ticketTemplatePageData(message string) (gin.H, error) {
	svc := ticketTemplateService()

	sets, err := svc.GetSets()
	if err != nil {
		return nil, err
	}
	epics, err := svc.GetEpics()
	if err != nil {
		return nil, err
	}
	items, err := svc.GetItems()
	if err != nil {
		return nil, err
	}
	typeOptions, err := svc.GetTicketTypeOptions()
	if err != nil {
		return nil, err
	}
	priorityOptions, err := svc.GetTicketPriorityOptions()
	if err != nil {
		return nil, err
	}
	statusOptions, err := svc.GetTicketStatusOptions()
	if err != nil {
		return nil, err
	}
	userOptions, err := svc.GetUserOptions()
	if err != nil {
		return nil, err
	}
	epicOptions, err := svc.GetEpicOptions()
	if err != nil {
		return nil, err
	}

	return gin.H{
		"Title":           "Ticket templates",
		"Page":            "ticketTemplate",
		"Error":           message,
		"Sets":            sets,
		"Epics":           epics,
		"Items":           items,
		"TypeOptions":     typeOptions,
		"PriorityOptions": priorityOptions,
		"StatusOptions":   statusOptions,
		"UserOptions":     userOptions,
		"EpicOptions":     epicOptions,
	}, nil
}

func parseOptionalFloat(value string) *float64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &parsed
}
