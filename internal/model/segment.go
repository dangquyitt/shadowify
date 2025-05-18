package model

type Segment struct {
	Id       string  `db:"id" json:"id"`
	VideoId  string  `db:"video_id" json:"videoId"`
	StartSec float32 `db:"start_sec" json:"startSec"`
	EndSec   float32 `db:"end_sec" json:"endSec"`
	Content  string  `db:"content" json:"content"`
}
