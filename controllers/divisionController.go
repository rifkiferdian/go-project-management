package controllers

import (
	"gobase-app/config"
	"gobase-app/repositories"
	"gobase-app/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func DivisionIndex(c *gin.Context) {
	divisionSvc := divisionService()
	divisions, err := divisionSvc.GetDivisions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "division.html", gin.H{
		"Title": "Divisions",
		"Page":  "division",
		"Rows":  divisions,
	})
}

func DivisionStore(c *gin.Context) {
	divisionSvc := divisionService()
	if err := divisionSvc.CreateDivision(c.PostForm("name")); err != nil {
		renderDivisionWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/divisions")
}

func DivisionUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		renderDivisionWithError(c, "divisi tidak valid")
		return
	}

	divisionSvc := divisionService()
	if err := divisionSvc.UpdateDivision(id, c.PostForm("name")); err != nil {
		renderDivisionWithError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/divisions")
}

func DivisionDelete(c *gin.Context) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid division id")
		return
	}

	divisionSvc := divisionService()
	if err := divisionSvc.DeleteDivision(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/divisions")
}

func divisionService() *services.DivisionService {
	return &services.DivisionService{
		Repo: &repositories.DivisionRepository{DB: config.DB},
	}
}

func renderDivisionWithError(c *gin.Context, message string) {
	divisionSvc := divisionService()
	divisions, err := divisionSvc.GetDivisions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "division.html", gin.H{
		"Title": "Divisions",
		"Page":  "division",
		"Rows":  divisions,
		"Error": message,
	})
}
