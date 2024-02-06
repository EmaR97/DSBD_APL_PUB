package handler

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"CamMonitoring/src/_interface"
	"CamMonitoring/src/entity"
	"CamMonitoring/src/utility"
	"github.com/gin-gonic/gin"
)

// GenerateImageUrl_ is a function type for generating image URLs.
type GenerateImageUrl_ func(string, time.Duration) (string, error)

// CameraExernalRequestHandler is a handler for camera-related operations.
type CameraExernalRequestHandler struct {
	_interface.Handler[entity.Camera]
	GenerateImageUrl_
}

func handleError(context *gin.Context, statusCode int, message string, err error) {
	utility.ErrorLog().Printf("%s: %v", message, err)
	context.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

// Next handles requests for the next camera frame.
func (h *CameraExernalRequestHandler) Next(context *gin.Context) {
	// Parse the lastSeen parameter from the request URL
	lastSeenIndex, err := strconv.Atoi(context.Param("lastSeen"))
	if err != nil {
		handleError(context, http.StatusBadRequest, "Invalid lastSeen parameter", err)
		return
	}

	// Get the camera ID from the request URL
	id := context.Param("id")
	utility.InfoLog().Printf("Id: %s", id)
	// Retrieve frames for the specified camera ID
	receivedFrames, err := h.Repository.GetByID(id)
	if err != nil {
		handleError(context, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	// Find the successive frame based on the lastSeenIndex
	nextFrame, found := receivedFrames.FindSuccessiveElement(int64(lastSeenIndex))
	if found {
		// Generate an image URL for the next frame
		imageUrl, err := h.GenerateImageUrl_(fmt.Sprintf("%s/%d.jpg", id, nextFrame), time.Minute*60)
		if err != nil {
			handleError(context, http.StatusInternalServerError, "Error sending next frame", err)
			return
		}

		// Prepare the response payload
		response := map[string]interface{}{
			"nextFrame": nextFrame,
			"imageUrl":  imageUrl,
		}

		context.JSON(http.StatusOK, response)
	} else {
		handleError(context, http.StatusServiceUnavailable, "No next frame available", nil)
	}
}

// Create handles the creation of a new camera.
func (h *CameraExernalRequestHandler) Create(context *gin.Context) {
	// Check if the request method is POST
	if context.Request.Method != http.MethodPost {
		handleError(context, http.StatusMethodNotAllowed, "Invalid request method", nil)
		return
	}

	utility.InfoLog().Println("Handling addCamera request")

	// Retrieve the user ID from the cookie
	userId, err := context.Cookie("user")
	if err != nil {
		handleError(context, http.StatusMethodNotAllowed, "Error retrieving user ID from cookie", err)
		return
	}

	// Generate a new camera with a unique ID
	newCamera := *entity.NewCamera(uuid.New().String(), userId)

	// Bind JSON data from the request to the new camera
	if err := context.ShouldBindJSON(&newCamera); err != nil {
		handleError(context, http.StatusBadRequest, "Error binding camera data", err)
		return
	}

	// Add the new camera to the repository
	if err := h.Repository.Add(newCamera); err != nil {
		handleError(context, http.StatusInternalServerError, "Error adding new camera", err)
		return
	}

	// Send a JSON response indicating successful camera addition
	context.JSON(http.StatusOK, gin.H{"message": "Camera added successfully", "cam_id": newCamera.Id})
}

// CameraLogin handles login requests for cameras.
func (h *CameraExernalRequestHandler) CameraLogin(c *gin.Context) {
	// Retrieve camera ID and password from the request
	cameraId := c.PostForm("cam_id")
	password := c.PostForm("password")

	// Get camera details from the repository
	camera, err := h.Repository.GetByID(cameraId)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error retrieving camera details", err)
		return
	}
	// Create a URL for the LoginPost endpoint (replace with the actual URL)
	loginURL := "http://localhost:8080/access/login"
	data := url.Values{
		"username": {camera.UserId},
		"password": {password},
		"cam_id":   {cameraId},
	}

	// Make a POST request to the LoginPost endpoint
	resp, err := http.PostForm(loginURL, data)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error making POST request to LoginPost endpoint", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utility.ErrorLog().Printf("Error closing bodyreader: %v", err)

		}
	}(resp.Body)

	// Set the "token" cookie based on the response
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			c.SetCookie("token", cookie.Value, 365*24*60*60, "/", "localhost", false, true)
		}
	}

	// Read the response body and send it to the user
	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error reading response body", err)
		return
	}

	c.String(resp.StatusCode, bodyBuffer.String())
}
