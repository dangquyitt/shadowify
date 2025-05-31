package model

type Segment struct {
	Id       string  `db:"id" json:"id"`
	VideoId  string  `db:"video_id" json:"video_id"`
	StartSec float32 `db:"start_sec" json:"start_sec"`
	EndSec   float32 `db:"end_sec" json:"end_sec"`
	Content  string  `db:"content" json:"content"`
}
