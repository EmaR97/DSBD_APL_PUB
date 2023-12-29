package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CustomLogFormatter(param gin.LogFormatterParams) string {
	colorCode := getStatusColorCode(param.StatusCode)
	log := fmt.Sprintf(
		"%s %s \"%s %s %s |%s%d\033[0m| %s\"",
		formatTimestamp(), param.ClientIP, param.Request.Proto, param.Method, param.Path, colorCode, param.StatusCode,
		param.Latency,
	)
	if param.ErrorMessage != "" {
		log = fmt.Sprintf("%s ERROR:\"%s\"\n", log, param.ErrorMessage)
	}
	return log + "\n"
}

func getStatusColorCode(statusCode int) string {
	switch {
	case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
		return "\033[32m" // Green for 2xx status codes
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		return "\033[33m" // Yellow for 4xx status codes
	default:
		return "\033[31m" // Red for other status codes
	}
}

func formatTimestamp() string {
	return fmt.Sprintf("[%s]", time.Now().Format("2006/01/02 - 15:04:05.000"))
}
