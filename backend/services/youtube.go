package services

import (
	"audioplay/config"
	"context"
	"errors"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeService struct {
	service *youtube.Service
	cfg     *config.Config
}

func NewYouTubeService(cfg *config.Config) (*YouTubeService, error) {
	if cfg.YouTubeAPIKey == "" {
		return nil, errors.New("YouTube API key not configured")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YouTubeAPIKey))
	if err != nil {
		return nil, err
	}

	return &YouTubeService{service: service, cfg: cfg}, nil
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
		thumbnail := ""
		if item.Snippet.Thumbnails != nil {
			thumbnail = item.Snippet.Thumbnails.Medium.Url // Utilisation de Medium au lieu de Default
		}

		results = append(results, VideoResult{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			Thumbnail:   thumbnail,
			Duration:    "", // YouTube nécessite une requête supplémentaire pour la durée
			Source:      "youtube",
		})
	}

	return results, nil
}

func (ys *YouTubeService) GetVideoDetails(videoID string) (*VideoResult, error) {
	// Implémentez cette méthode pour obtenir plus de détails sur une vidéo
	call := ys.service.Videos.List([]string{"id", "snippet", "contentDetails"}).
		Id(videoID)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, errors.New("video not found")
	}
	item := response.Items[0]
	return &VideoResult{
		ID:          item.Id,
		Title:       item.Snippet.Title,
		Description: item.Snippet.Description,
		Thumbnail:   item.Snippet.Thumbnails.Medium.Url,
		Duration:    item.ContentDetails.Duration,
		Source:      "youtube",
	}, nil
}

type VideoResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Duration    string `json:"duration"`
	Source      string `json:"source"`
}
