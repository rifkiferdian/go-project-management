package controllers

import (
	"net/http"

	"gobase-app/config"
	"gobase-app/repositories"
	"gobase-app/services"

	"github.com/gin-gonic/gin"
)

func ProfilePage(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	successMessage := ""
	if c.Query("success") == "password-updated" {
		successMessage = "Password berhasil diubah"
	}

	userService := &services.UserService{
		Repo: &repositories.UserRepository{DB: config.DB},
	}

	renderProfilePage(c, userService, userID, "", successMessage)
}

func ProfilePasswordUpdate(c *gin.Context) {
	userID := sessionUserID(c)
	if userID <= 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userService := &services.UserService{
		Repo: &repositories.UserRepository{DB: config.DB},
	}

	err := userService.ChangePassword(
		userID,
		c.PostForm("current_password"),
		c.PostForm("new_password"),
		c.PostForm("confirm_password"),
	)
	if err != nil {
		renderProfilePage(c, userService, userID, err.Error(), "")
		return
	}

	c.Redirect(http.StatusSeeOther, "/profile?success=password-updated")
}

func renderProfilePage(c *gin.Context, userService *services.UserService, userID int, errorMessage, successMessage string) {
	profile, err := userService.GetProfile(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "profile.html", gin.H{
		"Title":   "Profile",
		"Page":    "profile",
		"Profile": profile,
		"Error":   errorMessage,
		"Success": successMessage,
	})
}
