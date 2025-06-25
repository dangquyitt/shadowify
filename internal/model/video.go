package model

import (
	"shadowify/internal/database"
	"shadowify/internal/pagination"
)

type VideoType string

const (
	VideoPopular  VideoType = "popular"
	VideoRecent   VideoType = "recent"
	VideoFavorite VideoType = "favorite"
)

type VideoFilter struct {
	pagination.Pagination

	Q        *string   `json:"q" query:"q"`
	Type     VideoType `json:"type" query:"type"` // "popular", "recent", "favorite"
	Category *string   `json:"category" query:"category"`
}

type FavoriteVideoFilter struct {
	Q *string `json:"q" query:"q"`
	pagination.Pagination
}

type Video struct {
	Base
	ViewCount      int64                       `db:"view_count" json:"view_count"`
	Cefr           string                      `db:"cefr" json:"cefr"`
	LanguageId     string                      `db:"language_id" json:"language_id"`
	Title          string                      `db:"title" json:"title"`
	FullTitle      string                      `db:"full_title" json:"full_title"`
	Description    string                      `db:"description" json:"description"`
	YoutubeId      string                      `db:"youtube_id" json:"youtube_id"`
	Duration       int32                       `db:"duration" json:"duration"`
	DurationString string                      `db:"duration_string" json:"duration_string"`
	Thumbnail      string                      `db:"thumbnail" json:"thumbnail"`
	Tags           database.JSONType[[]string] `db:"tags" json:"tags"`
	Categories     database.JSONType[[]string] `db:"categories" json:"categories"`
}

type VideoDetail struct {
	Video
	IsFavorite bool `json:"is_favorite"`
}
