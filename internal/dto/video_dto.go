package dto

type CreateVideoRequest struct {
	YoutubeRawInput string `json:"youtubeRawInput"`
}

type UpdateVideoRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type VideoResponse struct {
	ID             string   `json:"id"`
	UserID         string   `json:"userId"`
	YoutubeID      string   `json:"youtubeId"`
	YoutubeURL     string   `json:"youtubeUrl"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Duration       int      `json:"duration"`
	DurationString string   `json:"durationString"`
	Thumbnail      string   `json:"thumbnail"`
	ViewCount      int64    `json:"viewCount"`
	Tags           []string `json:"tags"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}
