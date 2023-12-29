package entity

type User struct {
	Id          string `json:"id" bson:"_id"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	CamPassword string `json:"cam_password" bson:"cam_password"`
}

func (u User) ID() string {
	return u.Id
}
