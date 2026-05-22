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

func ApprovalFlowStepIndex(c *gin.Context) {
	stepSvc := approvalFlowStepService()
	selectedFlowID := parsePositiveInt(c.Query("approval_flow_id"))

	rows, err := stepSvc.GetApprovalFlowStepsByFlowID(selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	flows, err := stepSvc.GetFlowOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_steps.html", gin.H{
		"Title":          "Approval flow steps",
		"Page":           "approvalFlowStep",
		"Rows":           rows,
		"FlowOptions":    flows,
		"SelectedFlowID": selectedFlowID,
	})
}

func ApprovalFlowStepStore(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.PostForm("filter_approval_flow_id"))
	flowID, _ := strconv.Atoi(c.PostForm("approval_flow_id"))
	stepOrder, _ := strconv.Atoi(c.PostForm("step_order"))

	stepSvc := approvalFlowStepService()
	if err := stepSvc.CreateApprovalFlowStep(flowID, stepOrder, c.PostForm("step_name"), c.PostForm("approval_rule"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepsWithError(c, err.Error(), selectedFlowID)
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepsRedirectURL(selectedFlowID))
}

func ApprovalFlowStepUpdate(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.PostForm("filter_approval_flow_id"))
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderApprovalFlowStepsWithError(c, "approval flow step tidak valid", selectedFlowID)
		return
	}

	flowID, _ := strconv.Atoi(c.PostForm("approval_flow_id"))
	stepOrder, _ := strconv.Atoi(c.PostForm("step_order"))

	stepSvc := approvalFlowStepService()
	if err := stepSvc.UpdateApprovalFlowStep(id, flowID, stepOrder, c.PostForm("step_name"), c.PostForm("approval_rule"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepsWithError(c, err.Error(), selectedFlowID)
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepsRedirectURL(selectedFlowID))
}

func ApprovalFlowStepDelete(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.Query("approval_flow_id"))
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid approval flow step id")
		return
	}

	stepSvc := approvalFlowStepService()
	if err := stepSvc.DeleteApprovalFlowStep(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepsRedirectURL(selectedFlowID))
}

func approvalFlowStepService() *services.ApprovalFlowStepService {
	return &services.ApprovalFlowStepService{
		Repo: &repositories.ApprovalFlowStepRepository{DB: config.DB},
	}
}

func renderApprovalFlowStepsWithError(c *gin.Context, message string, selectedFlowID int) {
	stepSvc := approvalFlowStepService()
	rows, err := stepSvc.GetApprovalFlowStepsByFlowID(selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	flows, err := stepSvc.GetFlowOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_steps.html", gin.H{
		"Title":          "Approval flow steps",
		"Page":           "approvalFlowStep",
		"Rows":           rows,
		"FlowOptions":    flows,
		"SelectedFlowID": selectedFlowID,
		"Error":          message,
	})
}

func approvalFlowStepsRedirectURL(selectedFlowID int) string {
	if selectedFlowID > 0 {
		return "/approval-flow-steps?approval_flow_id=" + strconv.Itoa(selectedFlowID)
	}
	return "/approval-flow-steps"
}

func parsePositiveInt(value string) int {
	number, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || number <= 0 {
		return 0
	}
	return number
}
