package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// AudioService gère les opérations liées à l'audio
type AudioService struct {
	// Vous pouvez ajouter ici des dépendances comme un client DB, cache, etc.
}

// AudioResult représente un résultat audio
type AudioResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration  int    `json:"duration"` // en secondes
	StreamURL string `json:"stream_url"`
	Thumbnail string `json:"thumbnail"`
	Source    string `json:"source"` // "soundcloud", "jamendo", etc.
}

// SearchAudioContent recherche des contenus audio
func SearchAudioContent(query string) ([]AudioResult, error) {
	// Vous pouvez implémenter plusieurs sources ici
	// Pour cet exemple, nous utiliserons SoundCloud et Jamendo

	var results []AudioResult

	// 1. Essayer SoundCloud
	scResults, err := searchSoundCloud(query)
	if err != nil {
		log.Printf("SoundCloud search error: %v", err)
	} else {
		results = append(results, scResults...)
	}

	// 2. Essayer Jamendo
	jamendoResults, err := searchJamendo(query)
	if err != nil {
		log.Printf("Jamendo search error: %v", err)
	} else {
		results = append(results, jamendoResults...)
	}

	if len(results) == 0 {
		return nil, errors.New("no audio content found")
	}

	return results, nil
}

// searchSoundCloud recherche sur SoundCloud (exemple avec leur API)
func searchSoundCloud(query string) ([]AudioResult, error) {
	// NOTE: Vous aurez besoin d'un client ID SoundCloud
	const clientID = "VOTRE_CLIENT_ID_SOUNDCLOUD"
	const apiURL = "https://api.soundcloud.com/tracks"

	params := url.Values{}
	params.Add("q", query)
	params.Add("client_id", clientID)
	params.Add("limit", "10")
	params.Add("streamable", "true")

	resp, err := http.Get(fmt.Sprintf("%s?%s", apiURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("soundcloud API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("soundcloud API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read soundcloud response: %v", err)
	}

	var scTracks []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Duration int    `json:"duration"`
		User     struct {
			Username string `json:"username"`
		} `json:"user"`
		ArtworkURL string `json:"artwork_url"`
		StreamURL  string `json:"stream_url"`
	}

	if err := json.Unmarshal(body, &scTracks); err != nil {
		return nil, fmt.Errorf("failed to parse soundcloud response: %v", err)
	}

	var results []AudioResult
	for _, track := range scTracks {
		if track.StreamURL == "" {
			continue
		}

		streamURL := fmt.Sprintf("%s?client_id=%s", track.StreamURL, clientID)

		results = append(results, AudioResult{
			ID:        fmt.Sprintf("sc-%d", track.ID),
			Title:     track.Title,
			Artist:    track.User.Username,
			Duration:  track.Duration / 1000, // conversion ms -> s
			StreamURL: streamURL,
			Thumbnail: track.ArtworkURL,
			Source:    "soundcloud",
		})
	}

	return results, nil
}

// searchJamendo recherche sur Jamendo (exemple avec leur API)
func searchJamendo(query string) ([]AudioResult, error) {
	const clientID = "VOTRE_CLIENT_ID_JAMENDO"
	const apiURL = "https://api.jamendo.com/v3.0/tracks/"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("format", "json")
	params.Add("limit", "10")
	params.Add("search", query)
	params.Add("include", "musicinfo")

	resp, err := http.Get(fmt.Sprintf("%s?%s", apiURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("jamendo API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jamendo API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read jamendo response: %v", err)
	}

	var jamendoResponse struct {
		Headers map[string]interface{} `json:"headers"`
		Results []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Duration int    `json:"duration"`
			Artist   string `json:"artist_name"`
			Audio    string `json:"audio"`
			Image    string `json:"image"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &jamendoResponse); err != nil {
		return nil, fmt.Errorf("failed to parse jamendo response: %v", err)
	}

	var results []AudioResult
	for _, track := range jamendoResponse.Results {
		// Jamendo fournit directement l'URL MP3
		if !strings.HasSuffix(track.Audio, ".mp3") {
			continue
		}

		results = append(results, AudioResult{
			ID:        fmt.Sprintf("jam-%s", track.ID),
			Title:     track.Name,
			Artist:    track.Artist,
			Duration:  track.Duration,
			StreamURL: track.Audio,
			Thumbnail: track.Image,
			Source:    "jamendo",
		})
	}

	return results, nil
}

// GetAudioDetails récupère les détails d'un audio spécifique
func GetAudioDetails(audioID string) (*AudioResult, error) {
	// Implémentation basique - vous devrez l'adapter selon vos besoins
	parts := strings.Split(audioID, "-")
	if len(parts) < 2 {
		return nil, errors.New("invalid audio ID format")
	}

	source := parts[0]
	id := parts[1]

	switch source {
	case "sc":
		return getSoundCloudTrackDetails(id)
	case "jam":
		return getJamendoTrackDetails(id)
	default:
		return nil, errors.New("unknown audio source")
	}
}

func getSoundCloudTrackDetails(trackID string) (*AudioResult, error) {
	// Implémentez la récupération des détails depuis SoundCloud
	return nil, errors.New("not implemented")
}

func getJamendoTrackDetails(trackID string) (*AudioResult, error) {
	// Implémentez la récupération des détails depuis Jamendo
	return nil, errors.New("not implemented")
}
