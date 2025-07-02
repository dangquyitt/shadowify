package model

type TranscribeInput struct {
	AudioBase64 string `json:"audio_base64"`
}

type TranscribeOutput struct {
	Text string `json:"text"`
}

type EvaluateInput struct {
	AudioBase64 string `json:"audio_base64"`
}

type EvaluateOutput struct {
	Cefr      string `json:"cefr"`
	MeaningEN string `json:"meaning_en"`
	MeaningVI string `json:"meaning_vi"`
}
