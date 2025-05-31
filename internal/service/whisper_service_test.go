package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhisperService_DetectLanguage(t *testing.T) {
	service := NewWhisperService()
	audioFilePath := "../../tmp/nawe0Nl93IA.wav" // Replace with a valid audio file path for testing

	// Call the DetectLanguage method
	language, err := service.DetectLanguage(context.Background(), audioFilePath)
	if err != nil {
		t.Fatalf("DetectLanguage failed: %v", err)
	}

	// Check if the language is not empty
	assert.NotEmpty(t, language, "Expected a non-empty language result")
	assert.Equal(t, "ja", language, "Expected language to be 'jp'") // Adjust expected value based on your test audio file
}
