package services

import (
	"math/rand"
	"time"
)

type PlaylistService struct {
	audioService   *AudioService
	youtubeService *YouTubeService
}

func NewPlaylistService(audioService *AudioService, youtubeService *YouTubeService) *PlaylistService {
	return &PlaylistService{
		audioService:   audioService,
		youtubeService: youtubeService,
	}
}

func (ps *PlaylistService) GenerateAutoPlaylist(currentItemID string, source string) ([]PlaylistItem, error) {
	var items []PlaylistItem

	// Générer une playlist basée sur l'item actuel
	switch source {
	case "soundcloud", "jamendo":
		// Pour les audios, on génère une playlist aléatoire à partir du même service
		results, err := ps.audioService.Search("") // Recherche vide pour obtenir des résultats aléatoires
		if err != nil {
			return nil, err
		}

		for _, item := range results {
			if item.ID != currentItemID { // Exclure l'item actuel
				items = append(items, PlaylistItem{
					ID:      item.ID,
					Title:   item.Title,
					Type:    "audio",
					Source:  item.Source,
					Content: item.StreamURL,
				})
			}
		}

	case "youtube":
		// Pour YouTube, on pourrait utiliser les suggestions de YouTube
		results, err := ps.youtubeService.Search("") // Recherche vide pour obtenir des résultats aléatoires
		if err != nil {
			return nil, err
		}

		for _, item := range results {
			if item.ID != currentItemID { // Exclure l'item actuel
				items = append(items, PlaylistItem{
					ID:      item.ID,
					Title:   item.Title,
					Type:    "video",
					Source:  item.Source,
					Content: item.ID, // Pour YouTube, le contenu est l'ID de la vidéo
				})
			}
		}
	}

	// Mélanger la playlist
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

	// Limiter à 10 éléments
	if len(items) > 10 {
		items = items[:10]
	}

	return items, nil
}

type PlaylistItem struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`    // "audio" ou "video"
	Source  string `json:"source"`  // "soundcloud", "jamendo", "youtube"
	Content string `json:"content"` // URL ou ID selon le type
}
