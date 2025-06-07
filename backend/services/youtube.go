package services

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const API_KEY = "VOTRE_CLE_API_YOUTUBE"

func SearchYouTubeContent(query string) ([]VideoResult, error) {
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Printf("Error creating YouTube client: %v", err)
		return nil, err
	}

	// Appel à l'API YouTube
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		Type("video").
		MaxResults(10)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var results []VideoResult
	for _, item := range response.Items {
		results = append(results, VideoResult{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			Thumbnail:   item.Snippet.Thumbnails.Default.Url,
		})
	}

	return results, nil
}

type VideoResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
