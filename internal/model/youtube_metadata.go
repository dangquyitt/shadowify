package model

type YoutubeMetadata struct {
	Id             string   `json:"id"`
	Description    string   `json:"description"`
	Title          string   `json:"title"`
	FullTitle      string   `json:"fulltitle"`
	Duration       int32    `json:"duration"`
	DurationString string   `json:"duration_string"`
	Thumbnail      string   `json:"thumbnail"`
	Tags           []string `json:"tags"`
	Categories     []string `json:"categories"`
}
