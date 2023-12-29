package handler

import (
	"CamMonitoring/src/entity"
	"CamMonitoring/src/repository"
	"CamMonitoring/src/utility"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// VideoFeedHandler struct
type VideoFeedHandler struct {
	cameraRepository *repository.CameraRepository
}

// NewVideoFeedHandler creates a new VideoFeedHandler with the provided CameraRepository
func NewVideoFeedHandler(cameraRepository *repository.CameraRepository) *VideoFeedHandler {
	return &VideoFeedHandler{
		cameraRepository: cameraRepository,
	}
}

// Cam handles video feed rendering
func (h *VideoFeedHandler) Cam(c *gin.Context) {
	var lastSeen int64 = 0

	// Check if the "lastSeen" query parameter is provided
	if s := c.Query("lastSeen"); s != "" {
		a, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			utility.ErrorLog().Printf("Invalid lastSeen parameter: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid lastSeen parameter"})
			return
		}
		lastSeen = a
	}

	// Render the video feed HTML page
	c.HTML(
		http.StatusOK, "videoFeed.html", struct {
			Id       string
			LastSeen int64
		}{Id: c.Param("id"), LastSeen: lastSeen},
	)
}

// Index handles rendering the index page with user's cameras
func (h *VideoFeedHandler) Index(c *gin.Context) {
	// Retrieve user ID from the cookie
	userId, err := c.Cookie("user")
	if err != nil {
		utility.ErrorLog().Printf("User not authenticated: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Retrieve cameras associated with the user
	cameras, err := h.cameraRepository.GetAllByUser(userId)
	if err != nil {
		utility.ErrorLog().Printf("Error retrieving user cameras: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user cameras"})
		return
	}

	// Render the index page with camera data
	data := struct{ Cameras []entity.Camera }{Cameras: cameras}
	c.HTML(http.StatusOK, "index.html", data)

}
