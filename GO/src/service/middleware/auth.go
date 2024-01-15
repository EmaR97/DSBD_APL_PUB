package middleware

import (
	"CamMonitoring/src/utility"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// AuthMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware(validateLocation string, exemptRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the route is exempted from authentication
		for _, route := range exemptRoutes {
			if route == c.Request.URL.Path {
				c.Next()
				return
			}
		}
		// Check if "token" cookie is present
		token, err := c.Cookie("token")
		if err != nil {
			handleAuthError(c, http.StatusForbidden, "Authentication failed: Token not present")
			return
		}

		// Send a request to the handler to validate the token
		isValid, err := validateTokenWithHandler(token, validateLocation)
		if err != nil || !isValid {
			handleAuthError(c, http.StatusForbidden, "Authentication failed: Token validation error")
			return
		}

		c.Next()
	}
}

// validateTokenWithHandler sends a request to the handler to validate the token
func validateTokenWithHandler(token string, validateLocation string) (bool, error) {
	// Create a request with the token in the request body
	reqBody := bytes.NewBuffer([]byte("token=" + token))
	response, err := http.Post(validateLocation, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		utility.ErrorLog().Printf("Error sending validation request: %v", err)
		return false, fmt.Errorf("error sending validation request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utility.ErrorLog().Printf("Error closing bodyreader: %v", err)

		}
	}(response.Body)

	// Check the response status
	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// handleAuthError responds to the client with an authentication error and aborts the request
func handleAuthError(c *gin.Context, statusCode int, message string) {
	utility.WarningLog().Printf("Auth Error: %s", message)
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
