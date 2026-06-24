package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gobase-app/config"
	"gobase-app/models"

	"github.com/gin-gonic/gin"
)

func PublicProjectRequestPage(c *gin.Context) {
	success := strings.TrimSpace(c.Query("success")) == "1"
	renderPublicProjectRequestPage(c, http.StatusOK, "", success, map[string]string{})
}

func PublicProjectRequestListPage(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("q"))

	rows, err := getPublicProjectRequestRows(keyword)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "public_project_request_list.html", gin.H{
		"Title":   "Daftar Request Project Publik",
		"Rows":    rows,
		"Keyword": keyword,
	})
}

func PublicProjectRequestStore(c *gin.Context) {
	form := map[string]string{
		"project_name":          strings.TrimSpace(c.PostForm("project_name")),
		"project_description":   strings.TrimSpace(c.PostForm("project_description")),
		"business_goal":         strings.TrimSpace(c.PostForm("business_goal")),
		"request_division_id":   strings.TrimSpace(c.PostForm("request_division_id")),
		"requester_name":        strings.TrimSpace(c.PostForm("requester_name")),
		"requester_employee_id": strings.TrimSpace(c.PostForm("requester_employee_id")),
		"approval_flow_id":      strings.TrimSpace(c.PostForm("approval_flow_id")),
	}

	requestDivisionID, err := strconv.Atoi(form["request_division_id"])
	if err != nil || requestDivisionID <= 0 {
		renderPublicProjectRequestPage(c, http.StatusBadRequest, "Divisi requester wajib dipilih", false, form)
		return
	}

	approvalFlowID, err := strconv.Atoi(form["approval_flow_id"])
	if err != nil || approvalFlowID <= 0 {
		renderPublicProjectRequestPage(c, http.StatusBadRequest, "Approval flow wajib dipilih", false, form)
		return
	}

	file, err := c.FormFile("supporting_document")
	if err != nil {
		renderPublicProjectRequestPage(c, http.StatusBadRequest, "Dokumen pendukung wajib diupload (PDF/Image)", false, form)
		return
	}
	if err := validateSupportingDocumentFile(file); err != nil {
		renderPublicProjectRequestPage(c, http.StatusBadRequest, err.Error(), false, form)
		return
	}

	if err := validatePublicProjectRequestForm(form); err != nil {
		renderPublicProjectRequestPage(c, http.StatusBadRequest, err.Error(), false, form)
		return
	}

	if err := storePublicProjectRequest(c, form, requestDivisionID, approvalFlowID, file); err != nil {
		renderPublicProjectRequestPage(c, http.StatusInternalServerError, err.Error(), false, form)
		return
	}

	c.Redirect(http.StatusSeeOther, "/project-requests/new?success=1")
}

func renderPublicProjectRequestPage(c *gin.Context, status int, message string, success bool, form map[string]string) {
	divisions, err := getPublicDivisionOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	flows, err := getPublicApprovalFlowOptions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	form = normalizePublicProjectRequestForm(form)

	data := gin.H{
		"Title":     "Request Project Publik",
		"Divisions": divisions,
		"Flows":     flows,
		"Form":      form,
		"Success":   success,
	}

	if strings.TrimSpace(message) != "" {
		data["Error"] = message
	}

	c.HTML(status, "public_project_request.html", data)
}

type publicProjectRequestListItem struct {
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
	CreatedAtDisplay   string
	HasAttachment      bool
	AttachmentPath     string
	AttachmentName     string
	StepHistorySummary string
}

func normalizePublicProjectRequestForm(form map[string]string) map[string]string {
	result := map[string]string{
		"project_name":          "",
		"project_description":   "",
		"business_goal":         "",
		"request_division_id":   "",
		"requester_name":        "",
		"requester_employee_id": "",
		"approval_flow_id":      "",
	}
	for key, val := range form {
		result[key] = val
	}
	return result
}

func getPublicDivisionOptions() ([]models.DivisionOption, error) {
	rows, err := config.DB.Query(`
		SELECT id, name
		FROM divisions
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.DivisionOption
	for rows.Next() {
		var item models.DivisionOption
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func getPublicApprovalFlowOptions() ([]models.ApprovalFlowOption, error) {
	rows, err := config.DB.Query(`
		SELECT id, flow_code, flow_name, is_active
		FROM approval_flows
		WHERE entity_type = 'project_request'
			AND is_active = 1
		ORDER BY flow_name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.ApprovalFlowOption
	for rows.Next() {
		var item models.ApprovalFlowOption
		if err := rows.Scan(&item.ID, &item.FlowCode, &item.FlowName, &item.IsActive); err != nil {
			return nil, err
		}
		options = append(options, item)
	}

	return options, rows.Err()
}

func getPublicProjectRequestRows(keyword string) ([]publicProjectRequestListItem, error) {
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
			COALESCE(pr.requested_ticket_prefix, '-') AS requested_ticket_prefix,
			pr.status,
			pr.created_at,
			CASE WHEN COUNT(pra.id) > 0 THEN 1 ELSE 0 END AS has_attachment,
			COALESCE(MAX(pra.file_path), '') AS attachment_path,
			COALESCE(MAX(pra.original_name), '') AS attachment_name
		FROM project_requests pr
		LEFT JOIN divisions d ON d.id = pr.request_division_id
		LEFT JOIN project_request_attachments pra ON pra.project_request_id = pr.id
	`
	args := make([]interface{}, 0, 2)
	if keyword != "" {
		query += `
		WHERE pr.request_no LIKE ?
			OR pr.requester_employee_id LIKE ?
		`
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}
	query += `
		GROUP BY
			pr.id, pr.request_no, pr.project_name, pr.requester_name, pr.requester_employee_id,
			pr.project_description, pr.business_goal, d.name, pr.requested_ticket_prefix, pr.status, pr.created_at
		ORDER BY pr.created_at DESC
		LIMIT 200
	`

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]publicProjectRequestListItem, 0)
	for rows.Next() {
		var (
			item          publicProjectRequestListItem
			createdAt     time.Time
			hasAttachment int
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
			&createdAt,
			&hasAttachment,
			&item.AttachmentPath,
			&item.AttachmentName,
		); err != nil {
			return nil, err
		}

		item.StatusLabel = humanizeProjectRequestStatus(item.Status)
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04")
		item.HasAttachment = hasAttachment > 0
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	historyMap, err := getPublicStepHistorySummaryMap(result)
	if err != nil {
		return nil, err
	}
	for i := range result {
		result[i].StepHistorySummary = historyMap[result[i].ID]
	}

	return result, nil
}

func humanizeProjectRequestStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending":
		return "Menunggu Persetujuan"
	case "approved":
		return "Disetujui"
	case "rejected":
		return "Ditolak"
	case "synced_to_project":
		return "Masuk ke Project"
	case "cancelled":
		return "Dibatalkan"
	default:
		if status == "" {
			return "-"
		}
		return status
	}
}

func getPublicStepHistorySummaryMap(rows []publicProjectRequestListItem) (map[int64]string, error) {
	summary := make(map[int64]string, len(rows))
	if len(rows) == 0 {
		return summary, nil
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row.ID > 0 {
			ids = append(ids, row.ID)
		}
	}
	if len(ids) == 0 {
		return summary, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
		SELECT
			ss.project_request_id,
			ss.step_order,
			ss.step_name,
			ss.approval_rule,
			ss.status,
			COALESCE(GROUP_CONCAT(
				DISTINCT CONCAT(
					COALESCE(u.name, CONCAT('User#', pa.approver_user_id)),
					' (',
					CASE
						WHEN pa.decision = 'approved' THEN 'disetujui'
						WHEN pa.decision = 'rejected' THEN 'ditolak'
						ELSE pa.decision
					END,
					')'
				)
				ORDER BY pa.id SEPARATOR ', '
			), '') AS decisions,
			CONCAT_WS(', ',
				NULLIF((
					SELECT GROUP_CONCAT(DISTINCT direct_user.name ORDER BY direct_user.name SEPARATOR ', ')
					FROM approval_flow_step_approvers direct_approver
					JOIN users direct_user ON direct_user.id = direct_approver.approver_user_id
					WHERE direct_approver.approval_flow_step_id = ss.approval_flow_step_id
						AND direct_approver.approver_type = 'user'
						AND direct_approver.is_active = 1
						AND direct_user.deleted_at IS NULL
				), ''),
				NULLIF((
					SELECT GROUP_CONCAT(DISTINCT role_user.name ORDER BY role_user.name SEPARATOR ', ')
					FROM approval_flow_step_approvers role_approver
					JOIN model_has_roles role_member
						ON role_member.role_id = role_approver.approver_role_id
						AND role_member.model_type = ?
					JOIN users role_user ON role_user.id = role_member.model_id
					WHERE role_approver.approval_flow_step_id = ss.approval_flow_step_id
						AND role_approver.approver_type = 'role'
						AND role_approver.is_active = 1
						AND role_user.deleted_at IS NULL
				), ''),
				NULLIF((
					SELECT GROUP_CONCAT(DISTINCT division_user.name ORDER BY division_user.name SEPARATOR ', ')
					FROM approval_flow_step_approvers division_approver
					JOIN user_divisions division_member
						ON division_member.division_id = division_approver.approver_division_id
					JOIN users division_user ON division_user.id = division_member.user_id
					WHERE division_approver.approval_flow_step_id = ss.approval_flow_step_id
						AND division_approver.approver_type = 'division'
						AND division_approver.is_active = 1
						AND division_user.deleted_at IS NULL
				), '')
			) AS approvers
		FROM project_request_step_states ss
		LEFT JOIN project_request_approvals pa
			ON pa.project_request_id = ss.project_request_id
			AND pa.approval_flow_step_id = ss.approval_flow_step_id
		LEFT JOIN users u ON u.id = pa.approver_user_id
		WHERE ss.project_request_id IN (` + strings.Join(placeholders, ",") + `)
		GROUP BY ss.project_request_id, ss.step_order, ss.step_name, ss.approval_rule, ss.status
		ORDER BY ss.project_request_id ASC, ss.step_order ASC
	`

	queryArgs := make([]interface{}, 0, len(args)+1)
	queryArgs = append(queryArgs, userModelType)
	queryArgs = append(queryArgs, args...)

	historyRows, err := config.DB.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer historyRows.Close()

	linesByRequest := make(map[int64][]string, len(rows))
	for historyRows.Next() {
		var (
			requestID  int64
			stepOrder  int
			stepName   string
			stepRule   string
			stepStatus string
			decisions  string
			approvers  string
		)
		if err := historyRows.Scan(&requestID, &stepOrder, &stepName, &stepRule, &stepStatus, &decisions, &approvers); err != nil {
			return nil, err
		}

		line := fmt.Sprintf("Step %d - %s: %s", stepOrder, strings.TrimSpace(stepName), humanizeProjectRequestStepStatus(stepStatus))
		decisions = strings.TrimSpace(decisions)
		approvers = strings.TrimSpace(approvers)
		line += " | " + decisions + " | " + approvers + " | " + strings.TrimSpace(stepRule)
		linesByRequest[requestID] = append(linesByRequest[requestID], line)
	}
	if err := historyRows.Err(); err != nil {
		return nil, err
	}

	for _, row := range rows {
		lines := linesByRequest[row.ID]
		if len(lines) == 0 {
			summary[row.ID] = "Belum ada progres step approval."
			continue
		}
		summary[row.ID] = strings.Join(lines, "\n")
	}

	return summary, nil
}

func validatePublicProjectRequestForm(form map[string]string) error {
	if form["project_name"] == "" {
		return errors.New("Nama project wajib diisi")
	}
	if len(form["project_name"]) > 255 {
		return errors.New("Nama project maksimal 255 karakter")
	}
	if form["requester_name"] == "" {
		return errors.New("Nama requester wajib diisi")
	}
	if form["requester_employee_id"] == "" {
		return errors.New("Employee ID requester wajib diisi")
	}
	return nil
}

func storePublicProjectRequest(c *gin.Context, form map[string]string, requestDivisionID, approvalFlowID int, file *multipart.FileHeader) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if ok, err := existsDivision(tx, requestDivisionID); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("Divisi dengan id %d tidak ditemukan", requestDivisionID)
	}

	if ok, err := existsActiveProjectRequestFlow(tx, approvalFlowID); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("Approval flow dengan id %d tidak aktif atau tidak ditemukan", approvalFlowID)
	}

	if hasSteps, err := hasActiveFlowSteps(tx, approvalFlowID); err != nil {
		return err
	} else if !hasSteps {
		return errors.New("Approval flow belum memiliki step aktif")
	}

	requestNo := generateProjectRequestNo()
	prefix, err := generateDivisionTicketPrefix(tx, requestDivisionID)
	if err != nil {
		return err
	}
	systemEmail := buildSystemRequesterEmail(requestNo)

	res, err := tx.Exec(`
		INSERT INTO project_requests (
			request_no,
			project_name,
			project_description,
			business_goal,
			request_division_id,
			requested_ticket_prefix,
			requester_name,
			requester_email,
			requester_phone,
			requester_employee_id,
			approval_flow_id,
			current_step_order,
			status,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'pending', NOW(), NOW())
	`, requestNo, form["project_name"], nullableString(form["project_description"]), nullableString(form["business_goal"]), requestDivisionID, prefix, form["requester_name"], systemEmail, nil, nullableString(form["requester_employee_id"]), approvalFlowID)
	if err != nil {
		return err
	}

	requestID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	insertStepsResult, err := tx.Exec(`
		INSERT INTO project_request_step_states (
			project_request_id,
			approval_flow_step_id,
			step_order,
			step_name,
			approval_rule,
			status,
			created_at,
			updated_at
		)
		SELECT
			?,
			s.id,
			s.step_order,
			s.step_name,
			s.approval_rule,
			'pending',
			NOW(),
			NOW()
		FROM approval_flow_steps s
		WHERE s.approval_flow_id = ?
			AND s.is_active = 1
		ORDER BY s.step_order ASC
	`, requestID, approvalFlowID)
	if err != nil {
		return err
	}

	affected, err := insertStepsResult.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Gagal membuat snapshot step approval untuk request")
	}

	publicPath, storedName, err := storeProjectRequestSupportingDocument(c, requestID, requestNo, file)
	if err != nil {
		return err
	}
	if err := createProjectRequestAttachment(tx, requestID, file, storedName, publicPath); err != nil {
		_ = removeProjectRequestSupportingDocument(publicPath)
		return err
	}

	if err := enqueueProjectRequestStepNotifications(tx, requestID, 1); err != nil {
		_ = removeProjectRequestSupportingDocument(publicPath)
		return err
	}

	return tx.Commit()
}

func existsDivision(tx *sql.Tx, divisionID int) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM divisions
		WHERE id = ? AND deleted_at IS NULL
	`, divisionID).Scan(&count)
	return count > 0, err
}

func existsActiveProjectRequestFlow(tx *sql.Tx, flowID int) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flows
		WHERE id = ?
			AND entity_type = 'project_request'
			AND is_active = 1
	`, flowID).Scan(&count)
	return count > 0, err
}

func hasActiveFlowSteps(tx *sql.Tx, flowID int) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM approval_flow_steps
		WHERE approval_flow_id = ?
			AND is_active = 1
	`, flowID).Scan(&count)
	return count > 0, err
}

func isPrefixInUse(tx *sql.Tx, prefix string) (bool, error) {
	var projectCount int
	if err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM projects
		WHERE ticket_prefix = ?
			AND deleted_at IS NULL
	`, prefix).Scan(&projectCount); err != nil {
		return false, err
	}
	if projectCount > 0 {
		return true, nil
	}

	var requestCount int
	if err := tx.QueryRow(`
		SELECT COUNT(1)
		FROM project_requests
		WHERE requested_ticket_prefix = ?
			AND status IN ('pending', 'approved', 'synced_to_project')
	`, prefix).Scan(&requestCount); err != nil {
		return false, err
	}

	return requestCount > 0, nil
}

func generateProjectRequestNo() string {
	return fmt.Sprintf("PR-%s-%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)
}

func buildSystemRequesterEmail(requestNo string) string {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(requestNo), " ", ""))
	if normalized == "" {
		normalized = fmt.Sprintf("request-%d", time.Now().Unix())
	}
	return normalized + "@public-request.local"
}

func generateDivisionTicketPrefix(tx *sql.Tx, divisionID int) (string, error) {
	divisionPrefix, err := getDivisionPrefix(tx, divisionID)
	if err != nil {
		return "", err
	}

	startPos := len(divisionPrefix) + 1
	var lastRequestNo int64
	if err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(requested_ticket_prefix, ?) AS UNSIGNED)), 0)
		FROM project_requests
		WHERE requested_ticket_prefix LIKE CONCAT(?, '%')
			AND requested_ticket_prefix REGEXP CONCAT('^', ?, '[0-9]+$')
	`, startPos, divisionPrefix, divisionPrefix).Scan(&lastRequestNo); err != nil {
		return "", err
	}

	var lastProjectNo int64
	if err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(ticket_prefix, ?) AS UNSIGNED)), 0)
		FROM projects
		WHERE ticket_prefix LIKE CONCAT(?, '%')
			AND ticket_prefix REGEXP CONCAT('^', ?, '[0-9]+$')
			AND deleted_at IS NULL
	`, startPos, divisionPrefix, divisionPrefix).Scan(&lastProjectNo); err != nil {
		return "", err
	}

	lastNo := lastRequestNo
	if lastProjectNo > lastNo {
		lastNo = lastProjectNo
	}

	nextNo := lastNo + 1
	candidate := fmt.Sprintf("%s%d", divisionPrefix, nextNo)
	return candidate, nil
}

func getDivisionPrefix(tx *sql.Tx, divisionID int) (string, error) {
	var rawPrefix sql.NullString
	var divisionName sql.NullString
	err := tx.QueryRow(`
		SELECT prefix_division, name
		FROM divisions
		WHERE id = ?
			AND deleted_at IS NULL
	`, divisionID).Scan(&rawPrefix, &divisionName)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("Divisi dengan id %d tidak ditemukan", divisionID)
	}
	if err != nil {
		return "", err
	}

	prefix := strings.ToUpper(strings.TrimSpace(rawPrefix.String))
	if prefix == "" {
		prefix = autoDivisionPrefixFromName(divisionName.String)
		if prefix == "" {
			return "", errors.New("Gagal generate prefix otomatis dari nama divisi")
		}
		if _, err := tx.Exec(`
			UPDATE divisions
			SET prefix_division = ?, updated_at = NOW()
			WHERE id = ?
		`, prefix, divisionID); err != nil {
			return "", err
		}
	}
	if len(prefix) > 10 {
		return "", errors.New("Prefix division maksimal 10 karakter")
	}
	for _, ch := range prefix {
		if (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') {
			return "", errors.New("Prefix division hanya boleh huruf A-Z dan angka 0-9 tanpa spasi")
		}
	}

	return prefix, nil
}

func autoDivisionPrefixFromName(name string) string {
	name = strings.ToUpper(strings.TrimSpace(name))
	if name == "" {
		return ""
	}

	var compact strings.Builder
	compact.Grow(len(name))
	for _, ch := range name {
		if (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			compact.WriteRune(ch)
		}
	}

	cleaned := compact.String()
	if cleaned == "" {
		return ""
	}

	// Simpel: default 2 karakter awal agar mirip contoh AC (Accounting).
	if len(cleaned) > 2 {
		cleaned = cleaned[:2]
	}
	if len(cleaned) > 10 {
		cleaned = cleaned[:10]
	}
	return cleaned
}

func nullableString(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func validateSupportingDocumentFile(file *multipart.FileHeader) error {
	if file == nil {
		return errors.New("Dokumen pendukung tidak valid")
	}
	if file.Size <= 0 {
		return errors.New("Dokumen pendukung kosong")
	}
	if file.Size > 10*1024*1024 {
		return errors.New("Dokumen pendukung maksimal 10 MB")
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))
	switch ext {
	case ".pdf", ".png", ".jpg", ".jpeg", ".webp":
	default:
		return errors.New("Format dokumen hanya boleh PDF/JPG/JPEG/PNG/WEBP")
	}

	mimeType := strings.ToLower(strings.TrimSpace(file.Header.Get("Content-Type")))
	if mimeType == "" {
		return nil
	}

	if strings.HasPrefix(mimeType, "image/") || mimeType == "application/pdf" {
		return nil
	}

	return errors.New("Content type dokumen harus PDF atau image")
}

func storeProjectRequestSupportingDocument(c *gin.Context, requestID int64, requestNo string, file *multipart.FileHeader) (string, string, error) {
	uploadDir := filepath.Join("assets", "uploads", "project_requests", strconv.FormatInt(requestID, 10))
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", "", err
	}

	storedName := buildRequestNoFilename(requestNo, file.Filename)
	destination := filepath.Join(uploadDir, storedName)
	if err := c.SaveUploadedFile(file, destination); err != nil {
		return "", "", err
	}

	publicPath := "/assets/uploads/project_requests/" + strconv.FormatInt(requestID, 10) + "/" + storedName
	return publicPath, storedName, nil
}

func createProjectRequestAttachment(tx *sql.Tx, requestID int64, file *multipart.FileHeader, storedName, publicPath string) error {
	_, err := tx.Exec(`
		INSERT INTO project_request_attachments (
			project_request_id,
			original_name,
			file_name,
			file_path,
			file_size,
			mime_type,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, requestID, filepath.Base(file.Filename), storedName, publicPath, file.Size, strings.TrimSpace(file.Header.Get("Content-Type")))
	return err
}

func removeProjectRequestSupportingDocument(publicPath string) error {
	publicPath = strings.TrimSpace(publicPath)
	if publicPath == "" {
		return nil
	}

	relative := strings.TrimPrefix(publicPath, "/")
	relative = filepath.FromSlash(relative)
	return os.Remove(relative)
}

func buildRequestNoFilename(requestNo, originalFilename string) string {
	requestNo = strings.TrimSpace(requestNo)
	if requestNo == "" {
		requestNo = fmt.Sprintf("request-%d", time.Now().Unix())
	}

	var builder strings.Builder
	builder.Grow(len(requestNo))
	for _, ch := range requestNo {
		switch {
		case ch >= 'a' && ch <= 'z':
			builder.WriteRune(ch)
		case ch >= 'A' && ch <= 'Z':
			builder.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			builder.WriteRune(ch)
		case ch == '-' || ch == '_':
			builder.WriteRune(ch)
		default:
			builder.WriteRune('-')
		}
	}

	safeBase := strings.Trim(builder.String(), "-_")
	if safeBase == "" {
		safeBase = fmt.Sprintf("request-%d", time.Now().Unix())
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(originalFilename)))
	if ext == "" {
		ext = ".bin"
	}

	return safeBase + ext
}
