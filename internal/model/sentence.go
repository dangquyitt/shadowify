package model

type Sentence struct {
	Base
	SegmentId string `db:"segment_id" json:"segment_id"`
}
