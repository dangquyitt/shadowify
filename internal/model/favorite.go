package model

type Favorite struct {
	Base
	UserId  string `json:"user_id"`
	VideoId string `json:"video_id"`
}
