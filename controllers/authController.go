package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"gobase-app/config"
	helpers "gobase-app/helper"
	"gobase-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const userModelType = "App\\Models\\User"

func renderLoginPage(c *gin.Context, status int, errorMessage string) {
	kanbanStatuses, err := getProjectKanbanStatuses()
	if err != nil {
		kanbanStatuses = []models.ProjectStatusOption{}
	}

	kanbanProjects, err := getProjectKanbanProjects()
	if err != nil {
		kanbanProjects = []models.Project{}
	}

	data := gin.H{
		"Title":          "Login",
		"KanbanStatuses": kanbanStatuses,
		"KanbanProjects": kanbanProjects,
	}

	if strings.TrimSpace(errorMessage) != "" {
		data["Error"] = errorMessage
	}

	c.HTML(status, "login.html", data)
}

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil {
		c.Redirect(302, "/dashboard")
		return
	}
	renderLoginPage(c, http.StatusOK, "")
}

func LoginPost(c *gin.Context) {
	employeeID := strings.TrimSpace(c.PostForm("employee_id"))
	password := c.PostForm("password")
	if employeeID == "" || password == "" {
		renderLoginPage(c, http.StatusOK, "Employee ID dan password wajib diisi")
		return
	}

	var (
		userID int
		dbName string
		dbMail string
		dbPass sql.NullString
		dbRole sql.NullString
	)
	err := config.DB.QueryRow(`
		SELECT 
			u.id,
			u.name,
			u.email,
			u.password,
			COALESCE(GROUP_CONCAT(DISTINCT r.name ORDER BY r.name SEPARATOR ', '), '') AS role
		FROM users u
		LEFT JOIN model_has_roles mhr ON mhr.model_id = u.id AND mhr.model_type = ?
		LEFT JOIN roles r ON r.id = mhr.role_id
		WHERE u.employee_id = ? AND u.deleted_at IS NULL
		GROUP BY u.id, u.name, u.email, u.password
	`, userModelType, employeeID).
		Scan(&userID, &dbName, &dbMail, &dbPass, &dbRole)

	if err == sql.ErrNoRows {
		renderLoginPage(c, http.StatusOK, "Employee ID tidak ditemukan")
		return
	} else if err != nil {
		renderLoginPage(c, http.StatusInternalServerError, "Terjadi kesalahan saat mengambil data user")
		return
	}

	if !dbPass.Valid || bcrypt.CompareHashAndPassword([]byte(dbPass.String), []byte(password)) != nil {
		renderLoginPage(c, http.StatusOK, "Password salah")
		return
	}

	userInitials := helpers.Initials(dbName)
	session := sessions.Default(c)
	session.Set("user", models.SessionUser{
		UserID:          userID,
		Name:            dbName,
		Email:           dbMail,
		Initials:        userInitials,
		Role:            dbRole.String,
		IsAuthenticated: true,
	})
	session.Set("user_id", userID)
	if err := session.Save(); err != nil {
		renderLoginPage(c, http.StatusInternalServerError, "Gagal menyimpan sesi: "+err.Error())
		return
	}

	c.Redirect(302, "/dashboard")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, "/")
}

func CreateUser(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	employeeID := strings.TrimSpace(c.PostForm("employee_id"))
	password := c.PostForm("password")
	if name == "" || email == "" || employeeID == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, employee_id, dan password wajib diisi"})
		return
	}

	var existingUser string
	err := config.DB.QueryRow("SELECT email FROM users WHERE email = ? AND deleted_at IS NULL", email).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	err = config.DB.QueryRow("SELECT employee_id FROM users WHERE employee_id = ? AND deleted_at IS NULL", employeeID).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO users (name, email, employee_id, password, type, created_at, updated_at) VALUES (?, ?, ?, ?, 'db', NOW(), NOW())",
		name,
		email,
		employeeID,
		string(hashedPassword),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
