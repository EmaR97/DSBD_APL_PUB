package handler

import (
	"CamMonitoring/src/message"
	"CamMonitoring/src/repository"
	"CamMonitoring/src/utility"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

type GenerateImageUrl func(id string, timestamp int64) string

// CameraInternalRequestHandler handles the storage of frame information.
type CameraInternalRequestHandler struct {
	cameraRepository *repository.CameraRepository
	notify           func(jsonString string, key string) error
	sendResponse     func(response, addressCallback string) error
	generateImageUrl GenerateImageUrl
}

// NewFrameHandler creates a new CameraInternalRequestHandler instance.
func NewFrameHandler(
	cameraRepo *repository.CameraRepository, notify func(jsonString string, key string) error,
	generateImageUrl GenerateImageUrl, sendResponse func(response, addressCallback string) error,
) *CameraInternalRequestHandler {
	return &CameraInternalRequestHandler{
		cameraRepository: cameraRepo,
		notify:           notify,
		generateImageUrl: generateImageUrl,
		sendResponse:     sendResponse,
	}
}

// ProcessAndStoreFrame  processes the received message and stores frame information.
func (fh *CameraInternalRequestHandler) ProcessAndStoreFrame(messageContent string) error {
	frameInfo := &message.FrameInfo{}
	// Unmarshal the serialized data into the Notification message
	if err := proto.Unmarshal([]byte(messageContent), frameInfo); err != nil {
		utility.ErrorLog().Println("Error unmarshaling message:", err)
		return err
	}

	utility.InfoLog().Printf("Received frame: CamID=%s, Timestamp=%d\n", frameInfo.CamId, frameInfo.GetTimestamp())

	receivedFrames, err := fh.cameraRepository.GetByID(frameInfo.CamId)
	if err != nil {
		return err
	}

	if err := receivedFrames.AddNewFrame(
		frameInfo.GetTimestamp(), frameInfo.GetPersonDetected(),
	); err != nil {
		return err
	}

	if err := fh.cameraRepository.Update(receivedFrames); err != nil {
		return err
	}

	if frameInfo.GetPersonDetected() {
		link := fh.generateImageUrl(receivedFrames.Id, frameInfo.GetTimestamp())
		notification := message.Notification{
			CamId:     frameInfo.GetCamId(),
			Timestamp: frameInfo.GetTimestamp(),
			Link:      link,
		}
		marshaler := jsonpb.Marshaler{}
		jsonString, err := marshaler.MarshalToString(&notification)
		if err != nil {
			return err
		}
		return fh.notify(jsonString, frameInfo.GetCamId())

	}

	return nil
}
func (fh *CameraInternalRequestHandler) GetCamIds(messageContent string) error {

	request := &message.UserIdRequest{}
	// Unmarshal the serialized data into the Notification message
	if err := proto.Unmarshal([]byte(messageContent), request); err != nil {
		utility.ErrorLog().Println("Error unmarshalling message:", err)
		return err
	}

	// Your logic to retrieve cam_ids associated with user_id using the callback function
	userID := request.UserId
	camIds, err := fh.cameraRepository.GetAllCamIdsByUser(userID)
	if err != nil {
		utility.ErrorLog().Printf("Error retrieving CamIds: %v", err)
		return err
	}
	utility.InfoLog().Printf("UserId: %s, camIds: %s", userID, camIds)

	response := message.CamIdsResponse{CamIds: camIds, RequestId: request.RequestId}
	marshaller := jsonpb.Marshaler{}
	jsonString, err := marshaller.MarshalToString(&response)
	if err != nil {
		return err
	}
	return fh.sendResponse(jsonString, request.ResponseTopic)

}
