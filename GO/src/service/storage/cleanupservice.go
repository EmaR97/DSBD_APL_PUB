package storage

import (
	"fmt"
	"time"

	"CamMonitoring/src/_interface"
	"CamMonitoring/src/entity"
	"CamMonitoring/src/utility"
)

// CleanupService is a service for managing periodic cleanup of camera images.
type CleanupService struct {
	cameraRepository _interface.Repository[entity.Camera]
	interval         time.Duration
	maxAge           time.Duration
	minioClient      *MinioClient
	stopChan         chan struct{}
}

// NewCleanupService creates a new instance of CleanupService.
func NewCleanupService(
	cameraRepository _interface.Repository[entity.Camera], interval time.Duration, maxAge time.Duration,
	minioClient *MinioClient,
) *CleanupService {
	return &CleanupService{
		cameraRepository: cameraRepository,
		interval:         interval,
		maxAge:           maxAge,
		minioClient:      minioClient,
		stopChan:         make(chan struct{}),
	}
}

// Start starts the cleanup service, running the cleanup periodically.
func (cs *CleanupService) Start() {
	ticker := time.NewTicker(cs.interval)
	cs.cleanupAllCameras()

	for {
		select {
		case <-ticker.C:
			cs.cleanupAllCameras()
		case <-cs.stopChan:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the cleanup service.
func (cs *CleanupService) Stop() {
	close(cs.stopChan)
}

// cleanupAllCameras applies the cleanup operation to all cameras in the repository.
func (cs *CleanupService) cleanupAllCameras() {
	// Retrieve all cameras from the repository
	cameras, err := cs.cameraRepository.GetAll()
	if err != nil {
		utility.ErrorLog().Printf("Error retrieving cameras for cleanup: %v", err)
		return
	}

	// Apply cleanup to each camera
	for _, camera := range cameras {
		cs.cleanupCamera(&camera)
	}
}

// cleanupCamera applies the cleanup operation to a specific camera.
func (cs *CleanupService) cleanupCamera(camera *entity.Camera) {
	cs.CleanupOldImages(camera)

	// Update the camera entity in the repository
	err := cs.cameraRepository.Update(*camera)
	if err != nil {
		utility.ErrorLog().Printf("Error updating camera %s: %v", camera.ID, err)
		// Handle the error if needed
	}
}

// CleanupOldImages deletes frames older than the specified maxAge from the camera.
func (cs *CleanupService) CleanupOldImages(camera *entity.Camera) {
	// Get the current time
	currentTime := time.Now().UnixNano()

	// Identify frames older than maxAge
	var framesToKeep []entity.Frame
	for _, frame := range camera.Frames {
		frameTimestamp := frame.Timestamp
		frameAge := time.Duration(currentTime-frameTimestamp) * time.Nanosecond
		maxAge := camera.MaxAge
		if maxAge == 0 {
			maxAge = cs.maxAge
		}
		if frameAge <= maxAge {
			framesToKeep = append(framesToKeep, frame)
		} else {
			// Frame is older than maxAge, delete associated object from Minio
			objectName := fmt.Sprintf("%s/%d.jpg", camera.Id, frame.Timestamp)
			err := cs.minioClient.DeleteObject(objectName)
			if err != nil {
				utility.ErrorLog().Printf("Error deleting object %s: %v", objectName, err)
				// Handle the error if needed
			}
		}
	}

	// Update the camera's frame array with the frames to keep
	camera.Frames = framesToKeep
}
