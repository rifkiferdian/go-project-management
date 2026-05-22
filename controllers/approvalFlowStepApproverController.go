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

func ApprovalFlowStepApproverIndex(c *gin.Context) {
	svc := approvalFlowStepApproverService()
	selectedFlowID := parsePositiveInt(c.Query("approval_flow_id"))

	rows, err := svc.GetApproversByFlowID(selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stepOptions, userOptions, roleOptions, divisionOptions, err := loadStepApproverOptions(svc, selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	flowOptions, err := approvalFlowStepService().GetFlowOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_step_approvers.html", gin.H{
		"Title":           "Approval flow step approvers",
		"Page":            "approvalFlowStepApprover",
		"Rows":            rows,
		"StepOptions":     stepOptions,
		"FlowOptions":     flowOptions,
		"SelectedFlowID":  selectedFlowID,
		"UserOptions":     userOptions,
		"RoleOptions":     roleOptions,
		"DivisionOptions": divisionOptions,
	})
}

func ApprovalFlowStepApproverStore(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.PostForm("filter_approval_flow_id"))
	stepID, _ := strconv.Atoi(c.PostForm("approval_flow_step_id"))
	userID, _ := strconv.Atoi(c.PostForm("approver_user_id"))
	roleID, _ := strconv.Atoi(c.PostForm("approver_role_id"))
	divisionID, _ := strconv.Atoi(c.PostForm("approver_division_id"))

	svc := approvalFlowStepApproverService()
	if err := svc.CreateApprover(stepID, c.PostForm("approver_type"), userID, roleID, divisionID, checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepApproversWithError(c, err.Error(), selectedFlowID)
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepApproversRedirectURL(selectedFlowID))
}

func ApprovalFlowStepApproverUpdate(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.PostForm("filter_approval_flow_id"))
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderApprovalFlowStepApproversWithError(c, "approval flow step approver tidak valid", selectedFlowID)
		return
	}

	stepID, _ := strconv.Atoi(c.PostForm("approval_flow_step_id"))
	userID, _ := strconv.Atoi(c.PostForm("approver_user_id"))
	roleID, _ := strconv.Atoi(c.PostForm("approver_role_id"))
	divisionID, _ := strconv.Atoi(c.PostForm("approver_division_id"))

	svc := approvalFlowStepApproverService()
	if err := svc.UpdateApprover(id, stepID, c.PostForm("approver_type"), userID, roleID, divisionID, checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepApproversWithError(c, err.Error(), selectedFlowID)
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepApproversRedirectURL(selectedFlowID))
}

func ApprovalFlowStepApproverDelete(c *gin.Context) {
	selectedFlowID := parsePositiveInt(c.Query("approval_flow_id"))
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid approval flow step approver id")
		return
	}

	svc := approvalFlowStepApproverService()
	if err := svc.DeleteApprover(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, approvalFlowStepApproversRedirectURL(selectedFlowID))
}

func approvalFlowStepApproverService() *services.ApprovalFlowStepApproverService {
	return &services.ApprovalFlowStepApproverService{
		Repo: &repositories.ApprovalFlowStepApproverRepository{DB: config.DB},
	}
}

func renderApprovalFlowStepApproversWithError(c *gin.Context, message string, selectedFlowID int) {
	svc := approvalFlowStepApproverService()
	rows, err := svc.GetApproversByFlowID(selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stepOptions, userOptions, roleOptions, divisionOptions, err := loadStepApproverOptions(svc, selectedFlowID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	flowOptions, err := approvalFlowStepService().GetFlowOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_step_approvers.html", gin.H{
		"Title":           "Approval flow step approvers",
		"Page":            "approvalFlowStepApprover",
		"Rows":            rows,
		"StepOptions":     stepOptions,
		"FlowOptions":     flowOptions,
		"SelectedFlowID":  selectedFlowID,
		"UserOptions":     userOptions,
		"RoleOptions":     roleOptions,
		"DivisionOptions": divisionOptions,
		"Error":           message,
	})
}

func loadStepApproverOptions(svc *services.ApprovalFlowStepApproverService, selectedFlowID int) (interface{}, interface{}, interface{}, interface{}, error) {
	stepOptions, err := svc.GetStepOptionsByFlowID(selectedFlowID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	userOptions, err := svc.GetUserOptions()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	roleOptions, err := svc.GetRoleOptions()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	divisionOptions, err := svc.GetDivisionOptions()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return stepOptions, userOptions, roleOptions, divisionOptions, nil
}

func approvalFlowStepApproversRedirectURL(selectedFlowID int) string {
	if selectedFlowID > 0 {
		return "/approval-flow-step-approvers?approval_flow_id=" + strconv.Itoa(selectedFlowID)
	}
	return "/approval-flow-step-approvers"
}
