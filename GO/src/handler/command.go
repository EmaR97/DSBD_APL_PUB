package handler

import (
	"CamMonitoring/src/service/comunication"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CommandHandler struct
type CommandHandler struct {
	MqttClient *comunication.MqttClient
}

// NewCommandHandler creates a new CommandHandler with the provided MqttClient
func NewCommandHandler(mqttClient *comunication.MqttClient) *CommandHandler {
	return &CommandHandler{
		MqttClient: mqttClient,
	}
}

// Post handles the HTTP command request
func (h *CommandHandler) Post(c *gin.Context) {
	// Check if the request is a POST request
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	//log.Println("Handling Command request")

	// Bind the payload from the request
	payload := c.PostForm("payload")

	if payload == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Payload cannot be empty"})
		return
	}
	//log.Printf("Command request payload: %s\n", payload)

	// Extract the ID from the URL parameters
	id := c.Param("id")

	// Publish the message to MQTT
	topic := strings.Join([]string{"command", id}, "/")
	err := h.MqttClient.PublishMessage(topic, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error publishing message"})
		return
	}
	//log.Println("Command published successfully")

	// Respond to the client
	c.String(http.StatusOK, "Command published successfully")
}
