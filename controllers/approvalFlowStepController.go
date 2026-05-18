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
	rows, err := stepSvc.GetApprovalFlowSteps()
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
		"Title":       "Approval flow steps",
		"Page":        "approvalFlowStep",
		"Rows":        rows,
		"FlowOptions": flows,
	})
}

func ApprovalFlowStepStore(c *gin.Context) {
	flowID, _ := strconv.Atoi(c.PostForm("approval_flow_id"))
	stepOrder, _ := strconv.Atoi(c.PostForm("step_order"))

	stepSvc := approvalFlowStepService()
	if err := stepSvc.CreateApprovalFlowStep(flowID, stepOrder, c.PostForm("step_name"), c.PostForm("approval_rule"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepsWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flow-steps")
}

func ApprovalFlowStepUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderApprovalFlowStepsWithError(c, "approval flow step tidak valid")
		return
	}

	flowID, _ := strconv.Atoi(c.PostForm("approval_flow_id"))
	stepOrder, _ := strconv.Atoi(c.PostForm("step_order"))

	stepSvc := approvalFlowStepService()
	if err := stepSvc.UpdateApprovalFlowStep(id, flowID, stepOrder, c.PostForm("step_name"), c.PostForm("approval_rule"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepsWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flow-steps")
}

func ApprovalFlowStepDelete(c *gin.Context) {
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

	c.Redirect(http.StatusSeeOther, "/approval-flow-steps")
}

func approvalFlowStepService() *services.ApprovalFlowStepService {
	return &services.ApprovalFlowStepService{
		Repo: &repositories.ApprovalFlowStepRepository{DB: config.DB},
	}
}

func renderApprovalFlowStepsWithError(c *gin.Context, message string) {
	stepSvc := approvalFlowStepService()
	rows, err := stepSvc.GetApprovalFlowSteps()
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
		"Title":       "Approval flow steps",
		"Page":        "approvalFlowStep",
		"Rows":        rows,
		"FlowOptions": flows,
		"Error":       message,
	})
}
