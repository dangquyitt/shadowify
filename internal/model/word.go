package model

import "shadowify/internal/pagination"

type Word struct {
	Base
	MeaningVI string `db:"meaning_vi" json:"meaning_vi"`
	MeaningEN string `db:"meaning_en" json:"meaning_en"`
	UserId    string `db:"user_id" json:"user_id"`
	SegmentId string `db:"segment_id" json:"segment_id"`
}

type WordCreateRequest struct {
	MeaningEN string `json:"meaning_en" binding:"required"`
	SegmentId string `json:"segment_id"`
}

type WordFilter struct {
	pagination.Pagination

	UserId string
	Q      *string `json:"q" query:"q"`
}
