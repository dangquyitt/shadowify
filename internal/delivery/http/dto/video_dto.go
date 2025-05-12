package dto

type CreateVideoRequest struct {
	YoutubeId string `json:"youtubeId" validate:"required"`
}

type CreateVideoResponse struct {
	Id string `json:"id"`
}

type VideoResponse struct {
	Id             string   `json:"id"`
	YoutubeId      string   `json:"youtubeId"`
	Title          string   `json:"title"`
	Duration       int64    `json:"duration"`
	DurationString string   `json:"durationString"`
	Thumbnail      string   `json:"thumbnail"`
	Tags           []string `json:"tags"`
}
