package storage

import (
	"CamMonitoring/src/_interface"
	"CamMonitoring/src/entity"
	"CamMonitoring/src/message"
	"CamMonitoring/src/service/comunication"
	"CamMonitoring/src/utility"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
	"sync"
)

type GenerateImageUrl func(id string, timestamp int64) string

// FrameHandler handles the storage of frame information.
type FrameHandler struct {
	cameraRepository  _interface.Repository[entity.Camera]
	consumer          comunication.KafkaConsumer
	producer          comunication.KafkaProducer
	topicNotification string
	generateImageUrl  GenerateImageUrl
	stopChan          chan struct{}
	wg                sync.WaitGroup
}

// NewFrameHandler creates a new FrameHandler instance.
func NewFrameHandler(
	cameraRepo _interface.Repository[entity.Camera], consumer comunication.KafkaConsumer,

	producer comunication.KafkaProducer, topicNotification string, generateImageUrl GenerateImageUrl,
) *FrameHandler {
	return &FrameHandler{
		cameraRepository:  cameraRepo,
		consumer:          consumer,
		producer:          producer,
		topicNotification: topicNotification,
		generateImageUrl:  generateImageUrl,
		stopChan:          make(chan struct{}),
	}
}

// Start initiates the frame storage process.
func (fh *FrameHandler) Start() {
	fh.wg.Add(1)
	defer fh.wg.Done()
	for {
		select {
		case <-fh.stopChan:
			utility.InfoLog().Println("Frame storage process stopped")
			return
		default:
			messageContent, err := fh.consumer.Consume()
			if err != nil {
				utility.ErrorLog().Println("Error consuming message:", err)
				continue
			}
			if messageContent == "" {
				continue
			}
			// Process the received image content
			err = fh.handleMessage(messageContent)
			if err != nil {
				utility.ErrorLog().Println("Error processing and storing frame:", err)
				// Consider whether to continue or return an error to the caller
			}
		}
	}
}

// Stop stops the frame storage process.
func (fh *FrameHandler) Stop() {
	close(fh.stopChan)
	fh.wg.Wait()

}

// handleMessage processes the received message and stores frame information.
func (fh *FrameHandler) handleMessage(messageContent string) error {
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

		if err := fh.producer.Publish(fh.topicNotification, jsonString); err != nil {
			return err
		}
	}

	return nil
}
