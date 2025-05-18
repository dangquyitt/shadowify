# YouTube Extractor Service

A gRPC service for extracting metadata and transcripts from YouTube videos using yt-dlp and faster-whisper.

## Features

- Extract YouTube video metadata (title, description, duration, thumbnails, etc.)
- Extract YouTube video transcripts using Faster Whisper speech-to-text

## Setup

### Virtual Environment

To set up the Python virtual environment and install dependencies:

```bash
# Make the setup script executable
chmod +x setup.sh

# Run the setup script
./setup.sh
```

### Configuration

Edit the `config.yaml` file to configure the service:

- `temp_dir`: Directory for temporary files
- `whisper_model_size`: Size of the Whisper model (tiny, base, small, medium, large)
- `whisper_device`: Device to run Whisper on (cpu, cuda)
- `whisper_compute_type`: Compute type for Whisper (float32, float16, int8)

## Running the Service

### Local Development

```bash
# Activate the virtual environment
source venv/bin/activate

# Run the server
python server.py --port 50051 --config config.yaml
```

### Using Docker

```bash
# Build the Docker image
docker build -t shadowify/extractor-service .

# Run the Docker container
docker run -p 50051:50051 shadowify/extractor-service
```

## Using with Go Client

The Go client is implemented in the repository pattern at `/internal/extractor/repository/extractor_repository.go`.

Example usage:

```go
package main

import (
    "context"
    "log"

    "github.com/dangquyitt/shadowify/internal/extractor/repository"
)

func main() {
    // Create a new extractor repository
    extractorRepo, err := repository.NewExtractorRepository("localhost:50051")
    if err != nil {
        log.Fatalf("Failed to create extractor repository: %v", err)
    }
    defer extractorRepo.Close()

    // Extract metadata
    ctx := context.Background()
    video, err := extractorRepo.ExtractYoutubeMetadata(ctx, "VIDEO_ID")
    if err != nil {
        log.Fatalf("Failed to extract metadata: %v", err)
    }
    log.Printf("Video title: %s", video.Title)

    // Extract transcript
    transcript, err := extractorRepo.ExtractYoutubeTranscript(ctx, "VIDEO_ID")
    if err != nil {
        log.Fatalf("Failed to extract transcript: %v", err)
    }
    log.Printf("Transcript language: %s", transcript.LanguageCode)
    log.Printf("Number of segments: %d", len(transcript.Segments))
}
```
