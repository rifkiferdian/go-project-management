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

func ApprovalFlowIndex(c *gin.Context) {
	flowSvc := approvalFlowService()
	rows, err := flowSvc.GetApprovalFlows()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flows.html", gin.H{
		"Title": "Approval flows",
		"Page":  "approvalFlow",
		"Rows":  rows,
	})
}

func ApprovalFlowStore(c *gin.Context) {
	flowSvc := approvalFlowService()
	if err := flowSvc.CreateApprovalFlow(c.PostForm("flow_code"), c.PostForm("flow_name"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowsWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flows")
}

func ApprovalFlowUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderApprovalFlowsWithError(c, "approval flow tidak valid")
		return
	}

	flowSvc := approvalFlowService()
	if err := flowSvc.UpdateApprovalFlow(id, c.PostForm("flow_code"), c.PostForm("flow_name"), checkboxOn(c, "is_active")); err != nil {
		renderApprovalFlowsWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flows")
}

func ApprovalFlowDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid approval flow id")
		return
	}

	flowSvc := approvalFlowService()
	if err := flowSvc.DeleteApprovalFlow(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/approval-flows")
}

func approvalFlowService() *services.ApprovalFlowService {
	return &services.ApprovalFlowService{
		Repo: &repositories.ApprovalFlowRepository{DB: config.DB},
	}
}

func renderApprovalFlowsWithError(c *gin.Context, message string) {
	flowSvc := approvalFlowService()
	rows, err := flowSvc.GetApprovalFlows()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "approval_flows.html", gin.H{
		"Title": "Approval flows",
		"Page":  "approvalFlow",
		"Rows":  rows,
		"Error": message,
	})
}
