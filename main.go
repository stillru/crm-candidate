package main

import (
	"context"
	"crmcandidate/config"
	"crmcandidate/db"
	"crmcandidate/middleware"
	"crmcandidate/routes"
	"crmcandidate/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

func main() {
	// Инициализация логгера
	logger := config.InitLogger()

	// Инициализация метрик
	metricsProvider, _, err := config.InitMetrics()
	if err != nil {
		logger.Fatalf("Failed to initialize metrics: %v", err)
	}
	otel.SetMeterProvider(metricsProvider)

	// Инициализация базы данных
	dbConn, err := db.NewDB()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Инициализация сервисов
	candidateService := services.NewCandidateService(dbConn)

	// Создаем роутер
	router := chi.NewRouter()

	// Добавляем middleware для логирования и метрик
	meter := otel.Meter("crmcandidate")
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.MetricsMiddleware(meter))

	// Настройка маршрутов
	routes.SetupRoutes(router, candidateService)
	router.Handle("/metrics", promhttp.Handler())

	// Запуск сервера
	server := &http.Server{
		Addr:    ":8005",
		Handler: router,
	}

	go func() {
		logger.Println("Server is running on :8005")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Printf("Failed to shutdown server: %v", err)
	}
	logger.Println("Server stopped")
}
