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
	rows, err := svc.GetApprovers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stepOptions, userOptions, roleOptions, divisionOptions, err := loadStepApproverOptions(svc)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_step_approvers.html", gin.H{
		"Title":           "Approval flow step approvers",
		"Page":            "approvalFlowStepApprover",
		"Rows":            rows,
		"StepOptions":     stepOptions,
		"UserOptions":     userOptions,
		"RoleOptions":     roleOptions,
		"DivisionOptions": divisionOptions,
	})
}

func ApprovalFlowStepApproverStore(c *gin.Context) {
	stepID, _ := strconv.Atoi(c.PostForm("approval_flow_step_id"))
	userID, _ := strconv.Atoi(c.PostForm("approver_user_id"))
	roleID, _ := strconv.Atoi(c.PostForm("approver_role_id"))
	divisionID, _ := strconv.Atoi(c.PostForm("approver_division_id"))

	svc := approvalFlowStepApproverService()
	if err := svc.CreateApprover(stepID, c.PostForm("approver_type"), userID, roleID, divisionID, checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepApproversWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flow-step-approvers")
}

func ApprovalFlowStepApproverUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderApprovalFlowStepApproversWithError(c, "approval flow step approver tidak valid")
		return
	}

	stepID, _ := strconv.Atoi(c.PostForm("approval_flow_step_id"))
	userID, _ := strconv.Atoi(c.PostForm("approver_user_id"))
	roleID, _ := strconv.Atoi(c.PostForm("approver_role_id"))
	divisionID, _ := strconv.Atoi(c.PostForm("approver_division_id"))

	svc := approvalFlowStepApproverService()
	if err := svc.UpdateApprover(id, stepID, c.PostForm("approver_type"), userID, roleID, divisionID, checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowStepApproversWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flow-step-approvers")
}

func ApprovalFlowStepApproverDelete(c *gin.Context) {
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

	c.Redirect(http.StatusSeeOther, "/approval-flow-step-approvers")
}

func approvalFlowStepApproverService() *services.ApprovalFlowStepApproverService {
	return &services.ApprovalFlowStepApproverService{
		Repo: &repositories.ApprovalFlowStepApproverRepository{DB: config.DB},
	}
}

func renderApprovalFlowStepApproversWithError(c *gin.Context, message string) {
	svc := approvalFlowStepApproverService()
	rows, err := svc.GetApprovers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stepOptions, userOptions, roleOptions, divisionOptions, err := loadStepApproverOptions(svc)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flow_step_approvers.html", gin.H{
		"Title":           "Approval flow step approvers",
		"Page":            "approvalFlowStepApprover",
		"Rows":            rows,
		"StepOptions":     stepOptions,
		"UserOptions":     userOptions,
		"RoleOptions":     roleOptions,
		"DivisionOptions": divisionOptions,
		"Error":           message,
	})
}

func loadStepApproverOptions(svc *services.ApprovalFlowStepApproverService) (interface{}, interface{}, interface{}, interface{}, error) {
	stepOptions, err := svc.GetStepOptions()
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
