package config

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMetrics() (*sdkmetric.MeterProvider, *prometheus.Exporter, error) {
	// Создаем экспортер Prometheus
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	// Создаем провайдер метрик
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	// Устанавливаем глобальный провайдер метрик
	otel.SetMeterProvider(provider)

	// Регистрируем системные метрики
	registerSystemMetrics(provider)

	return provider, exporter, nil
}

func registerSystemMetrics(provider *sdkmetric.MeterProvider) {
	meter := provider.Meter("crmcandidate_system")

	// Метрика для потребляемой памяти
	memoryUsageGauge, err := meter.Int64ObservableGauge(
		"system_memory_usage",
		metric.WithUnit("bytes"),
		metric.WithDescription("Total memory usage in bytes"),
	)
	if err != nil {
		panic(err)
	}

	// Регистрируем callback для обновления метрики
	_, err = meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			o.ObserveInt64(memoryUsageGauge, int64(m.Sys))
			return nil
		},
		memoryUsageGauge,
	)
	if err != nil {
		panic(err)
	}
}
