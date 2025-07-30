package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MetricsRegistry holds all OpenTelemetry metrics for monitoring application performance and business metrics.
// It provides counters, histograms, and gauges for tracking HTTP requests, game activity, and errors.
type MetricsRegistry struct {
	HttpRequestsTotal     metric.Int64Counter
	HttpRequestDuration   metric.Float64Histogram
	HttpRequestsInFlight  metric.Int64UpDownCounter
	ActiveGames          metric.Int64UpDownCounter
	ActiveCustomDecks    metric.Int64UpDownCounter
	CardsDealt           metric.Int64Counter
	GamesCreated         metric.Int64Counter
	ApiErrors            metric.Int64Counter
}

// LogMiddleware creates a Gin middleware that logs all HTTP requests with detailed information.
// It captures request details, measures latency, and records metrics for monitoring and debugging.
func LogMiddleware(logger *zap.Logger, metricsRegistry *MetricsRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()
		clientIP := c.ClientIP()
		method := c.Request.Method

		gameID := ""
		if gameIDParam := c.Param("gameId"); gameIDParam != "" {
			gameID = gameIDParam
		}

		metricsRegistry.HttpRequestsInFlight.Add(c.Request.Context(), 1)
		defer metricsRegistry.HttpRequestsInFlight.Add(c.Request.Context(), -1)

		logger.Info("Request started",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("user_agent", userAgent),
			zap.String("client_ip", clientIP),
			zap.String("game_id", gameID),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []attribute.KeyValue{
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.Int("status_code", status),
			attribute.String("client_ip", clientIP),
		}

		if gameID != "" {
			attrs = append(attrs, attribute.String("game_id", gameID))
		}

		metricsRegistry.HttpRequestsTotal.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		metricsRegistry.HttpRequestDuration.Record(c.Request.Context(), latency.Seconds(), metric.WithAttributes(attrs...))

		logLevel := zapcore.InfoLevel
		logMsg := "Request completed"

		if status >= 400 {
			logLevel = zapcore.WarnLevel
			if status >= 500 {
				logLevel = zapcore.ErrorLevel
				metricsRegistry.ApiErrors.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
			}
			logMsg = "Request failed"
		}

		logger.Log(logLevel, logMsg,
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("user_agent", userAgent),
			zap.String("client_ip", clientIP),
			zap.String("game_id", gameID),
			zap.Int("status_code", status),
			zap.Duration("latency", latency),
			zap.String("latency_human", latency.String()),
		)
	}
}