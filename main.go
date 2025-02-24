package main

import (
	"crm-candidate/db"
	"crm-candidate/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	// Инициализация Gin
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Маршруты
	r.GET("/", handlers.IndexHandler)
	r.GET("/add-company", handlers.AddCompanyHandler)
	r.GET("/add-recruiter", handlers.AddRecruiterHandler)
	r.GET("/status", handlers.StatusHandler)
	r.POST("/save-company", handlers.SaveCompanyHandler)

	// Запуск сервера
	r.Run(":8090")
}
