package services

import (
	"audioplay/config"
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeService struct {
	service *youtube.Service
}

func NewYouTubeService(cfg *config.Config) (*YouTubeService, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YouTubeAPIKey))
	if err != nil {
		return nil, err
	}

	return &YouTubeService{service: service}, nil
}

func (ys *YouTubeService) Search(query string) ([]VideoResult, error) {
	call := ys.service.Search.List([]string{"id", "snippet"}).
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
