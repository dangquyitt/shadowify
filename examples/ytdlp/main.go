package main

import (
	"context"
	"encoding/json"
	"os/exec"
)

type YtDlpResult struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}

func main() {
	// Thay URL này bằng video bạn muốn test
	url := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	cmd := exec.CommandContext(context.TODO(), "yt-dlp", "--skip-download", "-j", url)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var result YtDlpResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		panic(err)
	}

	println("Title:", result.Title)
	println("Thumbnail:", result.Thumbnail)
}
