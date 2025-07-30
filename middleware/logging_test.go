package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogging(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router with logging middleware
	router := gin.New()
	router.Use(LogMiddleware(logger, nil))
	
	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Make a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that logging occurred
	logs := recorded.All()
	assert.Len(t, logs, 1)
	assert.Equal(t, "HTTP Request", logs[0].Message)

	// Check log fields
	fields := logs[0].Context
	assert.Len(t, fields, 6) // method, path, status, latency, user_agent, remote_ip

	// Check specific field values
	for _, field := range fields {
		switch field.Key {
		case "method":
			assert.Equal(t, "GET", field.String)
		case "path":
			assert.Equal(t, "/test", field.String)
		case "status":
			assert.Equal(t, int64(200), field.Integer)
		case "user_agent":
			assert.Equal(t, "test-agent", field.String)
		}
	}
}

func TestRequestLoggingWithError(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router with logging middleware
	router := gin.New()
	router.Use(LogMiddleware(logger, nil))
	
	// Add a test route that returns an error
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "test error"})
	})

	// Make a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Check that logging occurred
	logs := recorded.All()
	assert.Len(t, logs, 1)
	assert.Equal(t, "HTTP Request", logs[0].Message)

	// Check that status code is logged correctly
	found := false
	for _, field := range logs[0].Context {
		if field.Key == "status" {
			assert.Equal(t, int64(500), field.Integer)
			found = true
		}
	}
	assert.True(t, found, "Status field should be present")
}

func TestRequestLoggingDifferentMethods(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router with logging middleware
	router := gin.New()
	router.Use(LogMiddleware(logger, nil))
	
	// Add test routes for different methods
	router.POST("/post", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"created": true})
	})
	router.PUT("/put", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"updated": true})
	})
	router.DELETE("/delete", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, gin.H{})
	})

	// Test POST
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test PUT
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/put", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test DELETE
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/delete", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check that all requests were logged
	logs := recorded.All()
	assert.Len(t, logs, 3)

	// Check methods are logged correctly
	methods := []string{}
	for _, log := range logs {
		for _, field := range log.Context {
			if field.Key == "method" {
				methods = append(methods, field.String)
			}
		}
	}
	assert.Contains(t, methods, "POST")
	assert.Contains(t, methods, "PUT")
	assert.Contains(t, methods, "DELETE")
}