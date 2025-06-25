package model

import "shadowify/internal/pagination"

type Sentence struct {
	Base
	UserId    string `db:"user_id" json:"user_id"`
	SegmentId string `db:"segment_id" json:"segment_id"`
	MeaningVI string `db:"meaning_vi" json:"meaning_vi"`
	MeaningEN string `db:"meaning_en" json:"meaning_en"`
}

type SentenceCreateRequest struct {
	SegmentId string `json:"segment_id" binding:"required"`
	MeaningEN string `json:"meaning_en" binding:"required"`
}

type SentenceFilter struct {
	pagination.Pagination

	UserId    string  `json:"user_id" query:"user_id"`
	SegmentId string  `json:"segment_id" query:"segment_id"`
	Q         *string `json:"q" query:"q"`
}
