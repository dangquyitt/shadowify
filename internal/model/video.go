package model

import (
	"shadowify/internal/database"
	"shadowify/internal/pagination"
)

type VideoFilter struct {
	Q *string `json:"q" query:"q"`
	pagination.Pagination
}

type FavoriteVideoFilter struct {
	pagination.Pagination
}

type Video struct {
	Base
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
