package handlers

import (
    "crm-candidate/db"
    "github.com/gin-gonic/gin"
    "net/http"
)

func IndexHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "base.html", gin.H{
        "content": "Добро пожаловать в CRM для кандидата!",
    })
}

func AddCompanyHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "add_company.html", nil)
}

func AddRecruiterHandler(c *gin.Context) {
    companies, err := db.GetCompanies()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.HTML(http.StatusOK, "add_recruiter.html", gin.H{
        "Companies": companies,
    })
}

func StatusHandler(c *gin.Context) {
    interactions, err := db.GetInteractions()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.HTML(http.StatusOK, "status.html", gin.H{
        "Interactions": interactions,
    })
}

func SaveCompanyHandler(c *gin.Context) {
    // Логика сохранения компании
    c.Redirect(http.StatusFound, "/add-company")
}
