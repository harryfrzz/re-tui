package data

import (
	"embed"
	"encoding/json"

	"reddit-tui/internal/models"
)

//go:embed sample_posts.json
var sampleData embed.FS

func LoadSamplePosts() ([]models.Post, error) {
	data, err := sampleData.ReadFile("sample_posts.json")
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
