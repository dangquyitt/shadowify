package model

import "shadowify/internal/pagination"

type Sentence struct {
	Base
	SegmentId string `db:"segment_id" json:"segment_id"`
}

type SentenceCreateRequest struct {
	SegmentId string `json:"segment_id" binding:"required"`
}

type SentenceFilter struct {
	pagination.Pagination

	UserId    string  `json:"user_id" query:"user_id"`
	SegmentId string  `json:"segment_id" query:"segment_id"`
	Q         *string `json:"q" query:"q"`
}
