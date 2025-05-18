package model

import (
	"time"
)

type Video struct {
	Id             string    `db:"id" json:"id"`
	Title          string    `db:"title" json:"title"`
	FullTitle      string    `db:"full_title" json:"fullTitle"`
	Description    string    `db:"description" json:"description"`
	YoutubeId      string    `db:"youtube_id" json:"youtubeId"`
	Duration       int32     `db:"duration" json:"duration"`
	DurationString string    `db:"duration_string" json:"durationString"`
	Thumbnail      string    `db:"thumbnail" json:"thumbnail"`
	Tags           []string  `db:"tags" json:"tags"`
	Categories     []string  `db:"categories" json:"categories"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
