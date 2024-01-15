package entity

import (
	"sort"
	"time"
)

// Camera represents a camera entity with frames and related information.
type Camera struct {
	Id     string        `json:"id" bson:"_id"`
	Frames []Frame       `json:"frames" bson:"frames"`
	UserId string        `json:"user_id" bson:"user_id"` // Reference to the owning User
	MaxAge time.Duration `json:"maxAge" bson:"maxAge"`
}

// Frame represents a single frame captured by the camera.
type Frame struct {
	Timestamp      int64 `json:"ts" bson:"ts"`
	PersonDetected bool  `json:"pd" bson:"pd"`
}

// NewCamera creates a new Camera instance with default values.
func NewCamera(Id string, UserId string) *Camera {
	return &Camera{
		Id:     Id,
		Frames: nil,
		UserId: UserId,
		MaxAge: 0,
	}
}

// ID returns the ID of the camera.
func (c Camera) ID() string {
	return c.Id
}

// FindSuccessiveElement finds the first frame with a timestamp greater than the target time.
func (c *Camera) FindSuccessiveElement(targetTime int64) (int64, bool) {
	index := sort.Search(
		len(c.Frames), func(i int) bool {
			return c.Frames[i].Timestamp > targetTime
		},
	)

	if index < len(c.Frames) {
		return c.Frames[index].Timestamp, true
	}

	return 0, false
}

// AddNewFrame adds a new frame to the camera's frames list, maintaining chronological order.
func (c *Camera) AddNewFrame(timestampNs int64, personDetected bool) error {
	index := sort.Search(
		len(c.Frames), func(i int) bool {
			return c.Frames[i].Timestamp >= timestampNs
		},
	)

	newFrame := Frame{
		Timestamp:      timestampNs,
		PersonDetected: personDetected,
	}

	if index == len(c.Frames) {
		c.Frames = append(c.Frames, newFrame)
	} else {
		c.Frames = append(c.Frames[:index], append([]Frame{newFrame}, c.Frames[index:]...)...)
	}

	return nil
}
