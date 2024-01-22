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
	// Increment the WaitGroup counter to indicate the start of the frame storage process
	fh.wg.Add(1)
	defer fh.wg.Done() // Decrement the WaitGroup counter when the function exits

	// Infinite loop for continuously handling frames
	for {
		// Check for a stop signal from the stop channel
		select {
		case <-fh.stopChan:
			utility.InfoLog().Println("Frame storage process stopped")
			return // Exit the function when a stop signal is received
		default:
			// Attempt to consume a message from the consumer
			messageContent, err := fh.consumer.Consume()
			if err != nil {
				// Log an error and continue to the next iteration if there is an issue consuming the message
				utility.ErrorLog().Println("Error consuming message:", err)
				continue
			}

			// Skip further processing if the message content is empty
			if messageContent == "" {
				continue
			}
			// Process the received image content
			err = fh.handleMessage(messageContent)
			if err != nil {
				// Log an error if there is an issue processing and storing the frame
				utility.ErrorLog().Println("Error processing and storing frame:", err)
				// Consider whether to continue processing or return an error to the caller
			}
		}
	}
}

// Stop stops the frame storage process.
func (fh *FrameHandler) Stop() {
	// Close the stop channel to signal the termination of the frame storage process
	close(fh.stopChan)
	fh.wg.Wait() // Wait for the frame storage process to complete before returning
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
