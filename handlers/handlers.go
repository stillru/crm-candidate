package handlers

import (
	"crm-candidate/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"content": "index.html",
		"message": "Добро пожаловать в CRM для кандидата!",
	})
}

func AddCompanyHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"content": "add_company.html",
		"message": "Добавить компанию",
	})
}

func AddRecruiterHandler(c *gin.Context) {
	companies, err := db.GetCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "base.html", gin.H{
		"content":   "add_recruiter.html",
		"Companies": companies,
	})
}

func StatusHandler(c *gin.Context) {
	interactions, err := db.GetInteractions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "base.html", gin.H{
		"content":      "status.html",
		"message":      "Статус взаимодействий",
		"Interactions": interactions,
	})
}

func SaveCompanyHandler(c *gin.Context) {
	// Логика сохранения компании
	c.Redirect(http.StatusFound, "/add-company")
}
