package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gobase-app/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userNotificationItem struct {
	ID               int64  `json:"id"`
	EventType        string `json:"event_type"`
	EntityType       string `json:"entity_type"`
	EntityID         int64  `json:"entity_id"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	ActionURL        string `json:"action_url"`
	Severity         string `json:"severity"`
	IsRead           bool   `json:"is_read"`
	CreatedAtDisplay string `json:"created_at_display"`
}

func NotificationList(c *gin.Context) {
	userID := currentNotificationUserID(c)
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := parseNotificationLimit(c.DefaultQuery("limit", "8"))
	rows, err := getUserNotifications(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unreadCount, err := getUnreadNotificationCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":         rows,
		"unread_count": unreadCount,
	})
}

func NotificationAllPage(c *gin.Context) {
	userID := currentNotificationUserID(c)
	if userID <= 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	rows, err := getUserNotifications(userID, 200)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	unreadCount, err := getUnreadNotificationCount(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "notifications.html", gin.H{
		"Title":       "All Notifications",
		"Page":        "notification",
		"Rows":        rows,
		"UnreadCount": unreadCount,
	})
}

func NotificationMarkRead(c *gin.Context) {
	userID := currentNotificationUserID(c)
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notificationID, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || notificationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification id tidak valid"})
		return
	}

	if _, err := config.DB.Exec(`
		UPDATE user_notifications
		SET is_read = 1, read_at = NOW(), updated_at = NOW()
		WHERE id = ?
			AND user_id = ?
	`, notificationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unreadCount, err := getUnreadNotificationCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":           true,
		"unread_count": unreadCount,
	})
}

func NotificationMarkAllRead(c *gin.Context) {
	userID := currentNotificationUserID(c)
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if _, err := config.DB.Exec(`
		UPDATE user_notifications
		SET is_read = 1, read_at = NOW(), updated_at = NOW()
		WHERE user_id = ?
			AND is_read = 0
	`, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":           true,
		"unread_count": 0,
	})
}

func NotificationStream(c *gin.Context) {
	userID := currentNotificationUserID(c)
	if userID <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("X-Accel-Buffering", "no")

	lastUnread := -1

	send := func(event string, payload interface{}) bool {
		body, err := json.Marshal(payload)
		if err != nil {
			return false
		}
		if _, err := fmt.Fprintf(writer, "event: %s\n", event); err != nil {
			return false
		}
		if _, err := fmt.Fprintf(writer, "data: %s\n\n", string(body)); err != nil {
			return false
		}
		flusher.Flush()
		return true
	}

	unreadCount, err := getUnreadNotificationCount(userID)
	if err == nil {
		lastUnread = unreadCount
		_ = send("notification", gin.H{"unread_count": unreadCount})
	}

	pollTicker := time.NewTicker(5 * time.Second)
	heartbeatTicker := time.NewTicker(25 * time.Second)
	defer pollTicker.Stop()
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-pollTicker.C:
			unreadCount, err := getUnreadNotificationCount(userID)
			if err != nil {
				continue
			}
			if unreadCount != lastUnread {
				lastUnread = unreadCount
				if !send("notification", gin.H{"unread_count": unreadCount}) {
					return
				}
			}
		case <-heartbeatTicker.C:
			if _, err := fmt.Fprint(writer, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func enqueueProjectRequestStepNotifications(tx *sql.Tx, requestID int64, stepOrder int) error {
	var (
		stepID      int64
		requestNo   string
		projectName string
		stepName    string
	)

	err := tx.QueryRow(`
		SELECT
			ss.approval_flow_step_id,
			pr.request_no,
			pr.project_name,
			ss.step_name
		FROM project_requests pr
		JOIN project_request_step_states ss
			ON ss.project_request_id = pr.id
			AND ss.step_order = ?
		WHERE pr.id = ?
		LIMIT 1
	`, stepOrder, requestID).Scan(&stepID, &requestNo, &projectName, &stepName)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	approverIDs, err := getApproverUserIDsByStep(tx, stepID)
	if err != nil {
		return err
	}
	if len(approverIDs) == 0 {
		return nil
	}

	requestNoText := strings.TrimSpace(requestNo)
	projectNameText := strings.TrimSpace(projectName)
	if projectNameText == "" {
		projectNameText = requestNoText
	}

	dedupeKey := fmt.Sprintf("project_request:%d:step:%d:pending_approval", requestID, stepID)
	title := fmt.Sprintf("Persetujuan Dibutuhkan: %s", projectNameText)
	message := fmt.Sprintf(
		"Project request %s (%s) menunggu approval Anda pada step %d - %s.",
		requestNoText,
		projectNameText,
		stepOrder,
		strings.TrimSpace(stepName),
	)
	actionURL := fmt.Sprintf("/project-requests/manage/%d", requestID)

	for _, approverID := range approverIDs {
		if _, err := tx.Exec(`
			INSERT INTO user_notifications (
				user_id,
				event_type,
				entity_type,
				entity_id,
				dedupe_key,
				title,
				message,
				action_url,
				severity,
				is_read,
				created_at,
				updated_at
			) VALUES (?, 'project_request.pending_approval', 'project_request', ?, ?, ?, ?, ?, 'warning', 0, NOW(), NOW())
			ON DUPLICATE KEY UPDATE
				title = VALUES(title),
				message = VALUES(message),
				action_url = VALUES(action_url),
				severity = VALUES(severity),
				is_read = 0,
				read_at = NULL,
				updated_at = NOW()
		`, approverID, requestID, dedupeKey, title, message, actionURL); err != nil {
			return err
		}
	}

	return nil
}

func getApproverUserIDsByStep(tx *sql.Tx, stepID int64) ([]int, error) {
	rows, err := tx.Query(`
		SELECT DISTINCT target.user_id
		FROM (
			SELECT a.approver_user_id AS user_id
			FROM approval_flow_step_approvers a
			WHERE a.approval_flow_step_id = ?
				AND a.is_active = 1
				AND a.approver_type = 'user'
				AND a.approver_user_id IS NOT NULL

			UNION

			SELECT mhr.model_id AS user_id
			FROM approval_flow_step_approvers a
			JOIN model_has_roles mhr
				ON mhr.role_id = a.approver_role_id
				AND mhr.model_type = ?
			WHERE a.approval_flow_step_id = ?
				AND a.is_active = 1
				AND a.approver_type = 'role'
				AND a.approver_role_id IS NOT NULL

			UNION

			SELECT ud.user_id AS user_id
			FROM approval_flow_step_approvers a
			JOIN user_divisions ud
				ON ud.division_id = a.approver_division_id
			WHERE a.approval_flow_step_id = ?
				AND a.is_active = 1
				AND a.approver_type = 'division'
				AND a.approver_division_id IS NOT NULL
		) target
		JOIN users u
			ON u.id = target.user_id
			AND u.deleted_at IS NULL
		ORDER BY target.user_id ASC
	`, stepID, userModelType, stepID, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]int, 0)
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		result = append(result, userID)
	}

	return result, rows.Err()
}

func getUnreadNotificationCount(userID int) (int, error) {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(1)
		FROM user_notifications
		WHERE user_id = ?
			AND is_read = 0
	`, userID).Scan(&count)
	return count, err
}

func getUserNotifications(userID, limit int) ([]userNotificationItem, error) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			COALESCE(event_type, '') AS event_type,
			COALESCE(entity_type, '') AS entity_type,
			COALESCE(entity_id, 0) AS entity_id,
			COALESCE(title, '') AS title,
			COALESCE(message, '') AS message,
			COALESCE(action_url, '') AS action_url,
			COALESCE(severity, 'info') AS severity,
			is_read,
			created_at
		FROM user_notifications
		WHERE user_id = ?
		ORDER BY is_read ASC, created_at DESC
		LIMIT ?
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]userNotificationItem, 0)
	for rows.Next() {
		var (
			item      userNotificationItem
			isReadInt int
			createdAt time.Time
		)

		if err := rows.Scan(
			&item.ID,
			&item.EventType,
			&item.EntityType,
			&item.EntityID,
			&item.Title,
			&item.Message,
			&item.ActionURL,
			&item.Severity,
			&isReadInt,
			&createdAt,
		); err != nil {
			return nil, err
		}

		item.IsRead = isReadInt == 1
		item.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04")
		result = append(result, item)
	}

	return result, rows.Err()
}

func parseNotificationLimit(raw string) int {
	limit, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || limit <= 0 {
		return 8
	}
	if limit > 30 {
		return 30
	}
	return limit
}

func currentNotificationUserID(c *gin.Context) int {
	session := sessions.Default(c)

	if v := session.Get("user_id"); v != nil {
		switch id := v.(type) {
		case int:
			return id
		case int64:
			return int(id)
		case float64:
			return int(id)
		case string:
			parsed, _ := strconv.Atoi(strings.TrimSpace(id))
			return parsed
		}
	}

	if u := session.Get("user"); u != nil {
		switch val := u.(type) {
		case map[string]interface{}:
			return normalizeNotificationSessionID(val["user_id"])
		case gin.H:
			return normalizeNotificationSessionID(val["user_id"])
		}
	}

	return 0
}

func normalizeNotificationSessionID(value interface{}) int {
	switch id := value.(type) {
	case int:
		return id
	case int64:
		return int(id)
	case float64:
		return int(id)
	case string:
		parsed, _ := strconv.Atoi(strings.TrimSpace(id))
		return parsed
	}
	return 0
}
