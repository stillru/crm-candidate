package main

import (
    "crm-candidate/db"
    "crm-candidate/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    if err := db.InitDB(); err != nil {
        panic(err)
    }

    r := gin.Default()

    r.LoadHTMLGlob("templates/*")

    r.Static("/static", "./static")

    r.GET("/", handlers.IndexHandler)
    r.GET("/add-company", handlers.AddCompanyHandler)
    r.GET("/add-recruiter", handlers.AddRecruiterHandler)
    r.GET("/status", handlers.StatusHandler)

    r.POST("/save-company", handlers.SaveCompanyHandler)

    r.Run(":8090")
}
