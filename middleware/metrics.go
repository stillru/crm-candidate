package middleware

import (
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// responseWriterWrapper оборачивает http.ResponseWriter для сохранения статуса ответа
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader переопределяет метод WriteHeader для сохранения статуса
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// MetricsMiddleware собирает метрики HTTP-запросов
func MetricsMiddleware(meter metric.Meter) func(http.Handler) http.Handler {
	// Создаем счетчики и гистограммы
	requestCounter, _ := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	requestDuration, _ := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Обертываем ResponseWriter для сохранения статуса
			wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			// Передаем запрос следующему обработчику
			next.ServeHTTP(wrappedWriter, r)

			// Собираем метрики
			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(wrappedWriter.statusCode)

			requestCounter.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
					attribute.String("status", statusCode),
				),
			)
			requestDuration.Record(r.Context(), duration,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
					attribute.String("status", statusCode),
				),
			)
		})
	}
}
