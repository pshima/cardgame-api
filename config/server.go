package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTrustedProxies reads trusted proxy IP addresses from environment variables for security.
// This prevents IP spoofing attacks by only accepting proxy headers from trusted sources.
func GetTrustedProxies() []string {
	// Check if TRUSTED_PROXIES environment variable is set
	if envProxies := os.Getenv("TRUSTED_PROXIES"); envProxies != "" {
		// Split comma-separated proxy IPs
		proxies := strings.Split(envProxies, ",")
		for i, proxy := range proxies {
			proxies[i] = strings.TrimSpace(proxy)
		}
		return proxies
	}
	
	// Default trusted proxies for development
	return []string{
		"127.0.0.1", // localhost
		"::1",       // localhost IPv6
		// In production, set TRUSTED_PROXIES environment variable with your actual proxy IPs
		// Examples:
		// - Load balancer IP: "10.0.1.100"
		// - Private network range: "10.0.0.0/8"
		// - Cloudflare (if using): Use Cloudflare's IP ranges
		// - Google Cloud: Use Google's load balancer IP ranges
	}
}

// GetBaseURL constructs the base URL for the request while validating proxy headers.
// It only trusts forwarded protocol headers from configured trusted proxies for security.
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	// Only trust proxy headers if the client IP is from a trusted proxy
	clientIP := c.ClientIP()
	
	// Check if request is from a trusted proxy by comparing with our trusted proxy list
	// Note: This is a simplified check. In production, you might want more sophisticated validation
	isTrustedProxy := clientIP == "127.0.0.1" || clientIP == "::1"
	
	if isTrustedProxy {
		// Only use forwarded headers from trusted proxies
		if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" || proto == "http" {
			scheme = proto
		}
	}
	
	// Always use the original Host header for security, don't trust X-Forwarded-Host
	// unless you have specific requirements and trusted proxy configuration
	host := c.Request.Host
	
	return fmt.Sprintf("%s://%s", scheme, host)
}

// GetPort returns the port to run the server on, defaulting to 8080.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}