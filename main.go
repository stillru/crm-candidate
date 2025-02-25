package main

import (
	"crm-candidate/db"
	"crm-candidate/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Создание роутера
	r := chi.NewRouter()

	// Промежуточное ПО
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Статические файлы
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршруты
	r.Get("/", handlers.IndexHandler)
	r.Get("/add-company", handlers.AddCompanyHandler)
	r.Get("/add-recruiter", handlers.AddRecruiterHandler)
	r.Get("/status", handlers.StatusHandler)
	r.Post("/save-company", handlers.SaveCompanyHandler)

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8090")
	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
