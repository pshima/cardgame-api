package config

import (
	"os"
	"strings"

	"github.com/peteshima/cardgame-api/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger creates a structured JSON logger with configurable log levels.
// It uses zap for high-performance logging with proper time formatting and caller information.
func InitLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(getLogLevel())
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}

// getLogLevel reads the LOG_LEVEL environment variable and returns the corresponding zap log level.
// It defaults to INFO level if the environment variable is not set or invalid.
func getLogLevel() zapcore.Level {
	level := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch level {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitMetrics sets up OpenTelemetry metrics with Prometheus exporter for monitoring.
// It creates a meter provider and initializes all application metrics for observability.
func InitMetrics(logger *zap.Logger) (metric.Meter, *middleware.MetricsRegistry) {
	exporter, err := prometheus.New()
	if err != nil {
		logger.Fatal("Failed to create Prometheus exporter", zap.Error(err))
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	otel.SetMeterProvider(meterProvider)
	meter := meterProvider.Meter("cardgame-api")

	metricsRegistry, err := middleware.NewMetricsRegistry(meter)
	if err != nil {
		logger.Fatal("Failed to create metrics registry", zap.Error(err))
	}

	return meter, metricsRegistry
}