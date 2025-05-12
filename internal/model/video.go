package model

import "time"

type Video struct {
	Id             string    `json:"id"`
	YoutubeId      string    `json:"youtubeId"`
	Title          string    `json:"title"`
	Duration       int64     `json:"duration"`
	DurationString string    `json:"durationString"`
	Thumbnail      string    `json:"thumbnail"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
