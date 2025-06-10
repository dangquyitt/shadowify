package model

type SavedType string

const (
	SavedTypeSentence SavedType = "sentence"
	SavedTypeWord     SavedType = "word"
)

type Saved struct {
	Base
	UserId    string    `json:"user_id"`
	VideoId   string    `json:"video_id"`
	Type      SavedType `json:"type"`
	SegmentId string    `json:"segment_id"`
}
