package middleware

import (
	"CamMonitoring/src/utility"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
)

func AuthorizedIPMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := getClientIP(c.Request)

		// Check if the client's IP is in the list of allowed IPs
		if !isIPAllowed(clientIP, allowedIPs) {
			utility.ErrorLog().Printf("Blocked request from: %s", clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Continue with the next middleware or handler
		c.Next()
	}
}

// getClientIP gets the client's IP address from the request, considering X-Real-IP header and remote address.
func getClientIP(req *http.Request) string {
	// The X-Real-IP header is set by some reverse proxies to represent the real client IP
	if xRealIP := req.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// If X-Real-IP is not set, fall back to the remote address in the format "IP:port"
	// Use strings.Split to extract the IP part
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		utility.ErrorLog().Printf("Error extracting IP from remote address: %v", err)
		return ""
	}

	return ip
}

// isIPAllowed checks if the given IP is in the list of allowed IPs.
func isIPAllowed(ip string, allowedIPs []string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		utility.ErrorLog().Printf("Error parsing client IP: %s", ip)
		return false
	}

	for _, allowedIP := range allowedIPs {
		if net.ParseIP(strings.TrimSpace(allowedIP)).Equal(clientIP) {
			return true
		}
	}
	return false
}
