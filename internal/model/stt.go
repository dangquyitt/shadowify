package model

type TranscribeInput struct {
	AudioBase64 string `json:"audio_base64"`
}

type TranscribeOutput struct {
	Text string `json:"text"`
}
