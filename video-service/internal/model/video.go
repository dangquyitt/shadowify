package model

type Video struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Thumbnail  string `json:"thumbnail"`
	Source     string `json:"source"`
	ExternalId string `json:"externalId"`
}
