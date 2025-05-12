package model

type VideoSegment struct {
	Id        string `json:"id"`
	VideoId   string `json:"videoId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Text      string `json:"text"`
}
