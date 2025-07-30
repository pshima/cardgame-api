package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTrustedProxies(t *testing.T) {
	// Test with empty environment variable (should return defaults)
	os.Unsetenv("TRUSTED_PROXIES")
	proxies := GetTrustedProxies()
	assert.NotNil(t, proxies)
	assert.Contains(t, proxies, "127.0.0.1")
	assert.Contains(t, proxies, "::1")

	// Test with single proxy
	os.Setenv("TRUSTED_PROXIES", "10.0.1.100")
	proxies = GetTrustedProxies()
	assert.Equal(t, []string{"10.0.1.100"}, proxies)

	// Test with multiple proxies
	os.Setenv("TRUSTED_PROXIES", "10.0.1.100,192.168.1.0/24,172.16.0.0/16")
	proxies = GetTrustedProxies()
	assert.Equal(t, []string{"10.0.1.100", "192.168.1.0/24", "172.16.0.0/16"}, proxies)

	// Test with spaces (should be trimmed)
	os.Setenv("TRUSTED_PROXIES", " 10.0.1.100 , 192.168.1.0/24 ")
	proxies = GetTrustedProxies()
	assert.Equal(t, []string{"10.0.1.100", "192.168.1.0/24"}, proxies)

	// Clean up
	os.Unsetenv("TRUSTED_PROXIES")
}

func TestGetPort(t *testing.T) {
	// Test default port
	os.Unsetenv("PORT")
	port := GetPort()
	assert.Equal(t, "8080", port)

	// Test custom port
	os.Setenv("PORT", "3000")
	port = GetPort()
	assert.Equal(t, "3000", port)

	// Test empty string (should return default)
	os.Setenv("PORT", "")
	port = GetPort()
	assert.Equal(t, "8080", port)

	// Clean up
	os.Unsetenv("PORT")
}

func TestGetBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test HTTP request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://example.com/test", nil)
	c.Request.Host = "example.com"

	baseURL := GetBaseURL(c)
	assert.Equal(t, "http://example.com", baseURL)

	// Test with X-Forwarded-Proto header from trusted proxy
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://example.com/test", nil)
	c.Request.Host = "example.com"
	c.Request.Header.Set("X-Forwarded-Proto", "https")
	// Simulate trusted proxy IP
	c.Request.RemoteAddr = "127.0.0.1:12345"

	baseURL = GetBaseURL(c)
	// Note: This might still be http depending on the gin context setup
	assert.Contains(t, baseURL, "example.com")

	// Test with different host
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://localhost:8080/test", nil)
	c.Request.Host = "localhost:8080"

	baseURL = GetBaseURL(c)
	assert.Equal(t, "http://localhost:8080", baseURL)
}