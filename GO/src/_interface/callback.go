package _interface

// CamIdCallback function type for retrieving CamIds based on UserId
type CamIdCallback interface {
	GetAllCamIdsByUser(userID string) ([]string, error)
}
