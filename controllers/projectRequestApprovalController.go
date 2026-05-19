package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gobase-app/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type managerProjectRequestItem struct {
	ID                 int64
	RequestNo          string
	ProjectName        string
	ProjectDescription string
	BusinessGoal       string
	RequesterName      string
	RequesterEmployee  string
	DivisionName       string
	TicketPrefix       string
	Status             string
	StatusLabel        string
	FlowName           string
	CurrentStepOrder   int
	CurrentStepName    string
	CurrentRule        string
	CurrentRuleLabel   string
	ApprovedCount      int
	RequiredCount      int
	HasApproved        bool
	UserDecision       string
	UserDecisionLabel  string
	IsCurrentApprover  bool
	CreatedAtDisplay   string
	HasAttachment      bool
	AttachmentPath     string
	AttachmentName     string
	CanTakeAction      bool
}

type projectRequestStepHistory struct {
	StepOrder       int
	StepName        string
	ApprovalRule    string
	ApprovalRuleLbl string
	StepStatus      string
	StepStatusLbl   string
	MasterApprovers []projectRequestMasterApprover
	Decisions       []projectRequestStepDecision
}

type projectRequestMasterApprover struct {
	ApproverType  string
	ApproverLabel string
	IsActive      bool
}

type projectRequestStepDecision struct {
	ApproverName     string
	Decision         string
	DecisionLabel    string
	Note             string
	CreatedAtDisplay string
}

type pendingDecisionContext struct {
	RequestID          int64
	RequestStatus      string
	CurrentStepOrder   int
	StepID             int64
	StepStatus         string
	StepRule           string
	RequestNo          string
	ProjectName        string
	ProjectDescription sql.NullString
	BusinessGoal       sql.NullString
	DivisionID         int64
	TicketPrefix       string
}

func ProjectRequestManageIndex(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	keyword := strings.TrimSpace(c.Query("q"))
	rows, err := getManagerProjectRequestRows(userID, keyword)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "project_requests_manage.html", gin.H{
		"Title":   "Project Requests",
		"Page":    "projectRequestManage",
		"Rows":    rows,
		"Keyword": keyword,
		"Error":   strings.TrimSpace(c.Query("error")),
		"Success": strings.TrimSpace(c.Query("success")),
	})
}

func ProjectRequestManageDetail(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requestID, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || requestID <= 0 {
		c.String(http.StatusBadRequest, "invalid project request id")
		return
	}

	row, err := getManagerProjectRequestDetail(userID, requestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, "project request tidak ditemukan atau Anda tidak punya akses")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stepHistories, err := getProjectRequestStepHistories(requestID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "project_request_manage_detail.html", gin.H{
		"Title":         "Project Request Detail",
		"Page":          "projectRequestManage",
		"Row":           row,
		"StepHistories": stepHistories,
		"Error":         strings.TrimSpace(c.Query("error")),
		"Success":       strings.TrimSpace(c.Query("success")),
	})
}

func ProjectRequestApprove(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requestID, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || requestID <= 0 {
		redirectProjectRequestManage(c, "Invalid project request id", "")
		return
	}

	note := strings.TrimSpace(c.PostForm("note"))
	if err := applyProjectRequestDecision(requestID, userID, "approved", note); err != nil {
		redirectAfterProjectRequestAction(c, requestID, err.Error(), "")
		return
	}

	redirectAfterProjectRequestAction(c, requestID, "", "Request berhasil di-approve")
}

func ProjectRequestReject(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requestID, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || requestID <= 0 {
		redirectProjectRequestManage(c, "Invalid project request id", "")
		return
	}

	reason := strings.TrimSpace(c.PostForm("note"))
	if reason == "" {
		redirectAfterProjectRequestAction(c, requestID, "Alasan reject wajib diisi", "")
		return
	}

	if err := applyProjectRequestDecision(requestID, userID, "rejected", reason); err != nil {
		redirectAfterProjectRequestAction(c, requestID, err.Error(), "")
		return
	}

	redirectAfterProjectRequestAction(c, requestID, "", "Request berhasil di-reject")
}

func ProjectRequestDelete(c *gin.Context) {
	requestID, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || requestID <= 0 {
		redirectProjectRequestManage(c, "Invalid project request id", "")
		return
	}

	if err := deleteProjectRequest(requestID); err != nil {
		redirectProjectRequestManage(c, err.Error(), "")
		return
	}

	redirectProjectRequestManage(c, "", "Project request berhasil dihapus")
}

func redirectProjectRequestManage(c *gin.Context, errMsg, successMsg string) {
	values := url.Values{}
	if strings.TrimSpace(errMsg) != "" {
		values.Set("error", errMsg)
	}
	if strings.TrimSpace(successMsg) != "" {
		values.Set("success", successMsg)
	}

	target := "/project-requests/manage"
	if encoded := values.Encode(); encoded != "" {
		target += "?" + encoded
	}
	c.Redirect(http.StatusSeeOther, target)
}

func redirectProjectRequestDetail(c *gin.Context, requestID int64, errMsg, successMsg string) {
	values := url.Values{}
	if strings.TrimSpace(errMsg) != "" {
		values.Set("error", errMsg)
	}
	if strings.TrimSpace(successMsg) != "" {
		values.Set("success", successMsg)
	}

	target := fmt.Sprintf("/project-requests/manage/%d", requestID)
	if encoded := values.Encode(); encoded != "" {
		target += "?" + encoded
	}
	c.Redirect(http.StatusSeeOther, target)
}

func redirectAfterProjectRequestAction(c *gin.Context, requestID int64, errMsg, successMsg string) {
	if strings.ToLower(strings.TrimSpace(c.PostForm("redirect_to"))) == "detail" {
		redirectProjectRequestDetail(c, requestID, errMsg, successMsg)
		return
	}
	redirectProjectRequestManage(c, errMsg, successMsg)
}

func getManagerProjectRequestRows(userID int, keyword string) ([]managerProjectRequestItem, error) {
	query := `
		SELECT
			pr.id,
			pr.request_no,
			pr.project_name,
			COALESCE(pr.project_description, '') AS project_description,
			COALESCE(pr.business_goal, '') AS business_goal,
			pr.requester_name,
			COALESCE(pr.requester_employee_id, '-') AS requester_employee_id,
			COALESCE(d.name, '-') AS division_name,
			pr.requested_ticket_prefix,
			pr.status,
			COALESCE(f.flow_name, '-') AS flow_name,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			COALESCE((
				SELECT COUNT(DISTINCT pa.approver_user_id)
				FROM project_request_approvals pa
				WHERE pa.project_request_id = pr.id
					AND pa.approval_flow_step_id = ss.approval_flow_step_id
					AND pa.decision = 'approved'
			), 0) AS approved_count,
			COALESCE((
				SELECT COUNT(1)
				FROM approval_flow_step_approvers a
				WHERE a.approval_flow_step_id = ss.approval_flow_step_id
					AND a.is_active = 1
			), 0) AS required_count,
			CASE WHEN EXISTS (
				SELECT 1
				FROM approval_flow_step_approvers a_self
				WHERE a_self.approval_flow_step_id = ss.approval_flow_step_id
					AND a_self.is_active = 1
					AND (
						(a_self.approver_type = 'user' AND a_self.approver_user_id = ?)
						OR (
							a_self.approver_type = 'role'
							AND EXISTS (
								SELECT 1
								FROM model_has_roles mhr_self
								WHERE mhr_self.model_id = ?
									AND mhr_self.model_type = ?
									AND mhr_self.role_id = a_self.approver_role_id
							)
						)
						OR (
							a_self.approver_type = 'division'
							AND EXISTS (
								SELECT 1
								FROM user_divisions ud_self
								WHERE ud_self.user_id = ?
									AND ud_self.division_id = a_self.approver_division_id
							)
						)
					)
			) THEN 1 ELSE 0 END AS is_current_approver,
			CASE WHEN EXISTS (
				SELECT 1
				FROM project_request_approvals pa_self
				WHERE pa_self.project_request_id = pr.id
					AND pa_self.approval_flow_step_id = ss.approval_flow_step_id
					AND pa_self.approver_user_id = ?
			) THEN 1 ELSE 0 END AS has_approved,
			COALESCE((
				SELECT pa_self2.decision
				FROM project_request_approvals pa_self2
				WHERE pa_self2.project_request_id = pr.id
					AND pa_self2.approval_flow_step_id = ss.approval_flow_step_id
					AND pa_self2.approver_user_id = ?
				ORDER BY pa_self2.id DESC
				LIMIT 1
			), '') AS user_decision,
			pr.created_at,
			CASE WHEN COUNT(pra.id) > 0 THEN 1 ELSE 0 END AS has_attachment,
			COALESCE(MAX(pra.file_path), '') AS attachment_path,
			COALESCE(MAX(pra.original_name), '') AS attachment_name
		FROM project_requests pr
		JOIN project_request_step_states ss
			ON ss.project_request_id = pr.id
			AND ss.step_order = pr.current_step_order
		JOIN approval_flows f ON f.id = pr.approval_flow_id
		LEFT JOIN divisions d ON d.id = pr.request_division_id
		LEFT JOIN project_request_attachments pra ON pra.project_request_id = pr.id
		WHERE 1 = 1
	`
	args := []interface{}{userID, userID, userModelType, userID, userID, userID}

	if keyword != "" {
		query += `
			AND (
				pr.request_no LIKE ?
				OR pr.project_name LIKE ?
				OR pr.requester_name LIKE ?
				OR pr.requester_employee_id LIKE ?
			)
		`
		like := "%" + keyword + "%"
		args = append(args, like, like, like, like)
	}

	query += `
		GROUP BY
			pr.id,
			pr.request_no,
			pr.project_name,
			pr.project_description,
			pr.business_goal,
			pr.requester_name,
			pr.requester_employee_id,
			d.name,
			pr.requested_ticket_prefix,
			pr.status,
			f.flow_name,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			pr.created_at
		ORDER BY pr.created_at DESC
		LIMIT 300
	`

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]managerProjectRequestItem, 0)
	for rows.Next() {
		var (
			item          managerProjectRequestItem
			createdAt     time.Time
			hasAttachment int
			hasApproved   int
		)

		if err := rows.Scan(
			&item.ID,
			&item.RequestNo,
			&item.ProjectName,
			&item.ProjectDescription,
			&item.BusinessGoal,
			&item.RequesterName,
			&item.RequesterEmployee,
			&item.DivisionName,
			&item.TicketPrefix,
			&item.Status,
			&item.FlowName,
			&item.CurrentStepOrder,
			&item.CurrentStepName,
			&item.CurrentRule,
			&item.ApprovedCount,
			&item.RequiredCount,
			&item.IsCurrentApprover,
			&hasApproved,
			&item.UserDecision,
			&createdAt,
			&hasAttachment,
			&item.AttachmentPath,
			&item.AttachmentName,
		); err != nil {
			return nil, err
		}

		item.StatusLabel = humanizeProjectRequestStatus(item.Status)
		item.CurrentRuleLabel = humanizeApprovalRule(item.CurrentRule)
		item.HasApproved = hasApproved > 0
		item.UserDecisionLabel = humanizeProjectRequestDecision(item.UserDecision)
		item.HasAttachment = hasAttachment > 0
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04")
		item.CanTakeAction = item.Status == "pending" && strings.TrimSpace(item.UserDecision) == "" && item.IsCurrentApprover
		result = append(result, item)
	}

	return result, rows.Err()
}

func getManagerProjectRequestDetail(userID int, requestID int64) (managerProjectRequestItem, error) {
	query := `
		SELECT
			pr.id,
			pr.request_no,
			pr.project_name,
			COALESCE(pr.project_description, '') AS project_description,
			COALESCE(pr.business_goal, '') AS business_goal,
			pr.requester_name,
			COALESCE(pr.requester_employee_id, '-') AS requester_employee_id,
			COALESCE(d.name, '-') AS division_name,
			pr.requested_ticket_prefix,
			pr.status,
			COALESCE(f.flow_name, '-') AS flow_name,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			COALESCE((
				SELECT COUNT(DISTINCT pa.approver_user_id)
				FROM project_request_approvals pa
				WHERE pa.project_request_id = pr.id
					AND pa.approval_flow_step_id = ss.approval_flow_step_id
					AND pa.decision = 'approved'
			), 0) AS approved_count,
			COALESCE((
				SELECT COUNT(1)
				FROM approval_flow_step_approvers a
				WHERE a.approval_flow_step_id = ss.approval_flow_step_id
					AND a.is_active = 1
			), 0) AS required_count,
			CASE WHEN EXISTS (
				SELECT 1
				FROM approval_flow_step_approvers a_self
				WHERE a_self.approval_flow_step_id = ss.approval_flow_step_id
					AND a_self.is_active = 1
					AND (
						(a_self.approver_type = 'user' AND a_self.approver_user_id = ?)
						OR (
							a_self.approver_type = 'role'
							AND EXISTS (
								SELECT 1
								FROM model_has_roles mhr_self
								WHERE mhr_self.model_id = ?
									AND mhr_self.model_type = ?
									AND mhr_self.role_id = a_self.approver_role_id
							)
						)
						OR (
							a_self.approver_type = 'division'
							AND EXISTS (
								SELECT 1
								FROM user_divisions ud_self
								WHERE ud_self.user_id = ?
									AND ud_self.division_id = a_self.approver_division_id
							)
						)
					)
			) THEN 1 ELSE 0 END AS is_current_approver,
			CASE WHEN EXISTS (
				SELECT 1
				FROM project_request_approvals pa_self
				WHERE pa_self.project_request_id = pr.id
					AND pa_self.approval_flow_step_id = ss.approval_flow_step_id
					AND pa_self.approver_user_id = ?
			) THEN 1 ELSE 0 END AS has_approved,
			COALESCE((
				SELECT pa_self2.decision
				FROM project_request_approvals pa_self2
				WHERE pa_self2.project_request_id = pr.id
					AND pa_self2.approval_flow_step_id = ss.approval_flow_step_id
					AND pa_self2.approver_user_id = ?
				ORDER BY pa_self2.id DESC
				LIMIT 1
			), '') AS user_decision,
			pr.created_at,
			CASE WHEN COUNT(pra.id) > 0 THEN 1 ELSE 0 END AS has_attachment,
			COALESCE(MAX(pra.file_path), '') AS attachment_path,
			COALESCE(MAX(pra.original_name), '') AS attachment_name
		FROM project_requests pr
		JOIN project_request_step_states ss
			ON ss.project_request_id = pr.id
			AND ss.step_order = pr.current_step_order
		JOIN approval_flows f ON f.id = pr.approval_flow_id
		LEFT JOIN divisions d ON d.id = pr.request_division_id
		LEFT JOIN project_request_attachments pra ON pra.project_request_id = pr.id
		WHERE pr.id = ?
		GROUP BY
			pr.id,
			pr.request_no,
			pr.project_name,
			pr.project_description,
			pr.business_goal,
			pr.requester_name,
			pr.requester_employee_id,
			d.name,
			pr.requested_ticket_prefix,
			pr.status,
			f.flow_name,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			pr.created_at
		LIMIT 1
	`

	var (
		item          managerProjectRequestItem
		createdAt     time.Time
		hasAttachment int
		hasApproved   int
	)
	err := config.DB.QueryRow(
		query,
		userID,
		userID,
		userModelType,
		userID,
		userID,
		userID,
		requestID,
	).Scan(
		&item.ID,
		&item.RequestNo,
		&item.ProjectName,
		&item.ProjectDescription,
		&item.BusinessGoal,
		&item.RequesterName,
		&item.RequesterEmployee,
		&item.DivisionName,
		&item.TicketPrefix,
		&item.Status,
		&item.FlowName,
		&item.CurrentStepOrder,
		&item.CurrentStepName,
		&item.CurrentRule,
		&item.ApprovedCount,
		&item.RequiredCount,
		&item.IsCurrentApprover,
		&hasApproved,
		&item.UserDecision,
		&createdAt,
		&hasAttachment,
		&item.AttachmentPath,
		&item.AttachmentName,
	)
	if err != nil {
		return item, err
	}

	item.StatusLabel = humanizeProjectRequestStatus(item.Status)
	item.CurrentRuleLabel = humanizeApprovalRule(item.CurrentRule)
	item.HasApproved = hasApproved > 0
	item.UserDecisionLabel = humanizeProjectRequestDecision(item.UserDecision)
	item.HasAttachment = hasAttachment > 0
	item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04")
	item.CanTakeAction = item.Status == "pending" && strings.TrimSpace(item.UserDecision) == "" && item.IsCurrentApprover

	return item, nil
}

func getProjectRequestStepHistories(requestID int64) ([]projectRequestStepHistory, error) {
	rows, err := config.DB.Query(`
		SELECT
			ss.approval_flow_step_id,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			ss.status
		FROM project_request_step_states ss
		WHERE ss.project_request_id = ?
		ORDER BY ss.step_order ASC
	`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]projectRequestStepHistory, 0)
	for rows.Next() {
		var (
			stepID int64
			item   projectRequestStepHistory
		)

		if err := rows.Scan(
			&stepID,
			&item.StepOrder,
			&item.StepName,
			&item.ApprovalRule,
			&item.StepStatus,
		); err != nil {
			return nil, err
		}

		decisions, err := getProjectRequestStepDecisions(requestID, stepID)
		if err != nil {
			return nil, err
		}
		masterApprovers, err := getProjectRequestStepMasterApprovers(stepID)
		if err != nil {
			return nil, err
		}

		item.ApprovalRuleLbl = humanizeApprovalRule(item.ApprovalRule)
		item.StepStatusLbl = humanizeProjectRequestStepStatus(item.StepStatus)
		item.MasterApprovers = masterApprovers
		item.Decisions = decisions
		result = append(result, item)
	}

	return result, rows.Err()
}

func getProjectRequestStepDecisions(requestID, stepID int64) ([]projectRequestStepDecision, error) {
	rows, err := config.DB.Query(`
		SELECT
			COALESCE(u.name, CONCAT('User#', pa.approver_user_id)) AS approver_name,
			pa.decision,
			COALESCE(pa.note, '') AS note,
			pa.created_at
		FROM project_request_approvals pa
		LEFT JOIN users u ON u.id = pa.approver_user_id
		WHERE pa.project_request_id = ?
			AND pa.approval_flow_step_id = ?
		ORDER BY pa.created_at ASC, pa.id ASC
	`, requestID, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]projectRequestStepDecision, 0)
	for rows.Next() {
		var (
			item      projectRequestStepDecision
			createdAt time.Time
		)

		if err := rows.Scan(
			&item.ApproverName,
			&item.Decision,
			&item.Note,
			&createdAt,
		); err != nil {
			return nil, err
		}

		item.DecisionLabel = humanizeProjectRequestDecision(item.Decision)
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04")
		result = append(result, item)
	}

	return result, rows.Err()
}

func getProjectRequestStepMasterApprovers(stepID int64) ([]projectRequestMasterApprover, error) {
	rows, err := config.DB.Query(`
		SELECT
			a.approver_type,
			CASE
				WHEN a.approver_type = 'user' THEN COALESCE(u.name, CONCAT('User#', a.approver_user_id))
				WHEN a.approver_type = 'role' THEN COALESCE(r.name, CONCAT('Role#', a.approver_role_id))
				WHEN a.approver_type = 'division' THEN COALESCE(d.name, CONCAT('Division#', a.approver_division_id))
				ELSE '-'
			END AS approver_label,
			a.is_active
		FROM approval_flow_step_approvers a
		LEFT JOIN users u ON u.id = a.approver_user_id
		LEFT JOIN roles r ON r.id = a.approver_role_id
		LEFT JOIN divisions d ON d.id = a.approver_division_id
		WHERE a.approval_flow_step_id = ?
		ORDER BY a.id ASC
	`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]projectRequestMasterApprover, 0)
	for rows.Next() {
		var (
			item     projectRequestMasterApprover
			isActive int
		)

		if err := rows.Scan(
			&item.ApproverType,
			&item.ApproverLabel,
			&isActive,
		); err != nil {
			return nil, err
		}

		item.IsActive = isActive == 1
		result = append(result, item)
	}

	return result, rows.Err()
}

func deleteProjectRequest(requestID int64) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists int
	if err := tx.QueryRow(`SELECT COUNT(1) FROM project_requests WHERE id = ?`, requestID).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("Project request tidak ditemukan")
	}

	paths := make([]string, 0)
	rows, err := tx.Query(`
		SELECT COALESCE(file_path, '')
		FROM project_request_attachments
		WHERE project_request_id = ?
	`, requestID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			rows.Close()
			return err
		}
		paths = append(paths, p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	if _, err := tx.Exec(`DELETE FROM project_requests WHERE id = ?`, requestID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, p := range paths {
		_ = removeProjectRequestFileByPublicPath(p)
	}
	_ = os.RemoveAll(filepath.Join("assets", "uploads", "project_requests", strconv.FormatInt(requestID, 10)))

	return nil
}

func removeProjectRequestFileByPublicPath(publicPath string) error {
	publicPath = strings.TrimSpace(publicPath)
	if publicPath == "" {
		return nil
	}

	relative := strings.TrimPrefix(publicPath, "/")
	relative = filepath.FromSlash(relative)
	return os.Remove(relative)
}

func applyProjectRequestDecision(requestID int64, userID int, decision, note string) error {
	decision = strings.ToLower(strings.TrimSpace(decision))
	if decision != "approved" && decision != "rejected" {
		return errors.New("Decision tidak valid")
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx, err := lockPendingProjectRequestDecision(tx, requestID)
	if err != nil {
		return err
	}

	if strings.ToLower(strings.TrimSpace(ctx.RequestStatus)) != "pending" {
		return errors.New("Request ini sudah diproses")
	}
	if strings.ToLower(strings.TrimSpace(ctx.StepStatus)) != "pending" {
		return errors.New("Step approval aktif sudah tidak pending")
	}

	allowed, err := canUserApproveStep(tx, ctx.StepID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("Anda tidak terdaftar sebagai approver untuk step ini")
	}

	already, err := hasUserDecisionOnStep(tx, requestID, ctx.StepID, userID)
	if err != nil {
		return err
	}
	if already {
		return errors.New("Anda sudah mengambil keputusan pada step ini")
	}

	if _, err := tx.Exec(`
		INSERT INTO project_request_approvals (
			project_request_id,
			approval_flow_step_id,
			approver_user_id,
			decision,
			note,
			created_at
		) VALUES (?, ?, ?, ?, ?, NOW())
	`, requestID, ctx.StepID, userID, decision, toNullableText(note)); err != nil {
		return err
	}

	if decision == "rejected" {
		if err := rejectProjectRequest(tx, requestID, ctx.StepID, userID, note); err != nil {
			return err
		}
		return tx.Commit()
	}

	completeStep, err := isStepApprovalCompleted(tx, requestID, ctx.StepID, ctx.StepRule)
	if err != nil {
		return err
	}
	if !completeStep {
		return tx.Commit()
	}

	if _, err := tx.Exec(`
		UPDATE project_request_step_states
		SET status = 'approved', decided_by = ?, decided_at = NOW(), updated_at = NOW()
		WHERE project_request_id = ?
			AND approval_flow_step_id = ?
			AND status = 'pending'
	`, userID, requestID, ctx.StepID); err != nil {
		return err
	}

	nextStepOrder, hasNext, err := nextPendingStepOrder(tx, requestID)
	if err != nil {
		return err
	}

	if hasNext {
		if _, err := tx.Exec(`
			UPDATE project_requests
			SET current_step_order = ?, updated_at = NOW()
			WHERE id = ?
				AND status = 'pending'
		`, nextStepOrder, requestID); err != nil {
			return err
		}
		return tx.Commit()
	}

	if _, err := tx.Exec(`
		UPDATE project_requests
		SET status = 'approved', final_decided_by = ?, final_decided_at = NOW(), updated_at = NOW()
		WHERE id = ?
			AND status = 'pending'
	`, userID, requestID); err != nil {
		return err
	}

	if err := syncApprovedProjectRequest(tx, ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func lockPendingProjectRequestDecision(tx *sql.Tx, requestID int64) (pendingDecisionContext, error) {
	var ctx pendingDecisionContext
	err := tx.QueryRow(`
		SELECT
			pr.id,
			pr.status,
			pr.current_step_order,
			ss.approval_flow_step_id,
			ss.status,
			ss.approval_rule,
			pr.request_no,
			pr.project_name,
			pr.project_description,
			pr.business_goal,
			pr.request_division_id,
			pr.requested_ticket_prefix
		FROM project_requests pr
		JOIN project_request_step_states ss
			ON ss.project_request_id = pr.id
			AND ss.step_order = pr.current_step_order
		WHERE pr.id = ?
		FOR UPDATE
	`, requestID).Scan(
		&ctx.RequestID,
		&ctx.RequestStatus,
		&ctx.CurrentStepOrder,
		&ctx.StepID,
		&ctx.StepStatus,
		&ctx.StepRule,
		&ctx.RequestNo,
		&ctx.ProjectName,
		&ctx.ProjectDescription,
		&ctx.BusinessGoal,
		&ctx.DivisionID,
		&ctx.TicketPrefix,
	)
	if err == sql.ErrNoRows {
		return ctx, errors.New("Project request tidak ditemukan")
	}
	return ctx, err
}

func canUserApproveStep(tx *sql.Tx, stepID int64, userID int) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flow_step_approvers a
		WHERE a.approval_flow_step_id = ?
			AND a.is_active = 1
			AND (
				(a.approver_type = 'user' AND a.approver_user_id = ?)
				OR (
					a.approver_type = 'role'
					AND EXISTS (
						SELECT 1
						FROM model_has_roles mhr
						WHERE mhr.model_id = ?
							AND mhr.model_type = ?
							AND mhr.role_id = a.approver_role_id
					)
				)
				OR (
					a.approver_type = 'division'
					AND EXISTS (
						SELECT 1
						FROM user_divisions ud
						WHERE ud.user_id = ?
							AND ud.division_id = a.approver_division_id
					)
				)
			)
	`, stepID, userID, userID, userModelType, userID).Scan(&count)
	return count > 0, err
}

func hasUserDecisionOnStep(tx *sql.Tx, requestID, stepID int64, userID int) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM project_request_approvals
		WHERE project_request_id = ?
			AND approval_flow_step_id = ?
			AND approver_user_id = ?
	`, requestID, stepID, userID).Scan(&count)
	return count > 0, err
}

func rejectProjectRequest(tx *sql.Tx, requestID, stepID int64, userID int, reason string) error {
	if _, err := tx.Exec(`
		UPDATE project_request_step_states
		SET status = 'rejected', decided_by = ?, decided_at = NOW(), updated_at = NOW()
		WHERE project_request_id = ?
			AND approval_flow_step_id = ?
			AND status = 'pending'
	`, userID, requestID, stepID); err != nil {
		return err
	}

	_, err := tx.Exec(`
		UPDATE project_requests
		SET
			status = 'rejected',
			final_decided_by = ?,
			final_decided_at = NOW(),
			rejection_reason = ?,
			updated_at = NOW()
		WHERE id = ?
			AND status = 'pending'
	`, userID, strings.TrimSpace(reason), requestID)
	return err
}

func isStepApprovalCompleted(tx *sql.Tx, requestID, stepID int64, rule string) (bool, error) {
	rule = strings.ToLower(strings.TrimSpace(rule))
	switch rule {
	case "any":
		var approvedCount int
		if err := tx.QueryRow(`
			SELECT COUNT(1)
			FROM project_request_approvals
			WHERE project_request_id = ?
				AND approval_flow_step_id = ?
				AND decision = 'approved'
		`, requestID, stepID).Scan(&approvedCount); err != nil {
			return false, err
		}
		return approvedCount > 0, nil
	case "all":
		var (
			total     int
			satisfied int
		)
		if err := tx.QueryRow(`
			SELECT
				COUNT(1) AS total,
				COALESCE(SUM(CASE WHEN EXISTS (
					SELECT 1
					FROM project_request_approvals pa
					WHERE pa.project_request_id = ?
						AND pa.approval_flow_step_id = ?
						AND pa.decision = 'approved'
						AND (
							(a.approver_type = 'user' AND pa.approver_user_id = a.approver_user_id)
							OR (
								a.approver_type = 'role'
								AND EXISTS (
									SELECT 1
									FROM model_has_roles mhr
									WHERE mhr.model_id = pa.approver_user_id
										AND mhr.model_type = ?
										AND mhr.role_id = a.approver_role_id
								)
							)
							OR (
								a.approver_type = 'division'
								AND EXISTS (
									SELECT 1
									FROM user_divisions ud
									WHERE ud.user_id = pa.approver_user_id
										AND ud.division_id = a.approver_division_id
								)
							)
						)
				) THEN 1 ELSE 0 END), 0) AS satisfied
			FROM approval_flow_step_approvers a
			WHERE a.approval_flow_step_id = ?
				AND a.is_active = 1
		`, requestID, stepID, userModelType, stepID).Scan(&total, &satisfied); err != nil {
			return false, err
		}
		if total == 0 {
			return false, errors.New("Step approval belum memiliki approver aktif")
		}
		return satisfied >= total, nil
	default:
		return false, errors.New("approval rule tidak valid")
	}
}

func nextPendingStepOrder(tx *sql.Tx, requestID int64) (int, bool, error) {
	var nextStep int
	err := tx.QueryRow(`
		SELECT step_order
		FROM project_request_step_states
		WHERE project_request_id = ?
			AND status = 'pending'
		ORDER BY step_order ASC
		LIMIT 1
	`, requestID).Scan(&nextStep)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return nextStep, true, nil
}

func syncApprovedProjectRequest(tx *sql.Tx, ctx pendingDecisionContext) error {
	var existing int
	if err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM projects
		WHERE ticket_prefix = ?
			AND deleted_at IS NULL
	`, ctx.TicketPrefix).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return fmt.Errorf("Ticket prefix %s sudah dipakai project lain", ctx.TicketPrefix)
	}

	var statusID int64
	if err := tx.QueryRow(`
		SELECT id
		FROM project_statuses
		WHERE LOWER(TRIM(name)) = 'request received'
			AND deleted_at IS NULL
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&statusID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Status project 'Request Received' belum tersedia")
		}
		return err
	}

	var priorityID sql.NullInt64
	if err := tx.QueryRow(`
		SELECT id
		FROM project_priorities
		WHERE is_default = 1
			AND deleted_at IS NULL
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&priorityID); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		priorityID = sql.NullInt64{Valid: false}
	}

	insertRes, err := tx.Exec(`
		INSERT INTO projects (
			name,
			description,
			owner_id,
			developer_id,
			status_id,
			priority_id,
			ticket_prefix,
			status_type,
			type,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'default', 'kanban', NOW(), NOW())
	`,
		ctx.ProjectName,
		toNullableText(ctx.ProjectDescription.String),
		nil,
		nil,
		statusID,
		priorityID,
		ctx.TicketPrefix,
	)
	if err != nil {
		return err
	}

	projectID, err := insertRes.LastInsertId()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO project_divisions (project_id, division_id, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)
	`, projectID, ctx.DivisionID); err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE project_requests
		SET
			project_id = ?,
			status = 'synced_to_project',
			updated_at = NOW()
		WHERE id = ?
	`, projectID, ctx.RequestID)
	return err
}

func sessionUserID(c *gin.Context) int {
	session := sessions.Default(c)

	if v := session.Get("user_id"); v != nil {
		switch id := v.(type) {
		case int:
			return id
		case int64:
			return int(id)
		case float64:
			return int(id)
		}
	}

	if u := session.Get("user"); u != nil {
		switch val := u.(type) {
		case map[string]interface{}:
			if id, ok := val["user_id"]; ok {
				return normalizeSessionID(id)
			}
			if id, ok := val["UserID"]; ok {
				return normalizeSessionID(id)
			}
		case gin.H:
			if id, ok := val["user_id"]; ok {
				return normalizeSessionID(id)
			}
			if id, ok := val["UserID"]; ok {
				return normalizeSessionID(id)
			}
		}
	}

	return 0
}

func normalizeSessionID(value interface{}) int {
	switch id := value.(type) {
	case int:
		return id
	case int64:
		return int(id)
	case float64:
		return int(id)
	default:
		return 0
	}
}

func humanizeApprovalRule(rule string) string {
	switch strings.ToLower(strings.TrimSpace(rule)) {
	case "all":
		return "Semua Approver"
	case "any":
		return "Salah Satu Approver"
	default:
		if strings.TrimSpace(rule) == "" {
			return "-"
		}
		return rule
	}
}

func humanizeProjectRequestDecision(decision string) string {
	switch strings.ToLower(strings.TrimSpace(decision)) {
	case "approved":
		return "Disetujui"
	case "rejected":
		return "Ditolak"
	default:
		return ""
	}
}

func humanizeProjectRequestStepStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "approved":
		return "Disetujui"
	case "rejected":
		return "Ditolak"
	case "pending":
		return "Menunggu"
	case "skipped":
		return "Dilewati"
	default:
		if strings.TrimSpace(status) == "" {
			return "-"
		}
		return status
	}
}

func toNullableText(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}
