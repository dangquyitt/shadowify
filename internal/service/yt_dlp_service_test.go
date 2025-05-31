package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMetadata(t *testing.T) {
	youutubeId := "nawe0Nl93IA"
	service := NewYTDLPService()
	metadata, filePath, err := service.DownloadAndExtract(context.Background(), youutubeId)
	assert.NoError(t, err)
	// log.Printf("Metadata: %+v", metadata)
	t.Logf("File Path: %s", filePath)
	assert.NotNil(t, metadata)
	assert.Equal(t, youutubeId, metadata.Id)
}
