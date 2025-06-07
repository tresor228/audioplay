package services

import (
	"audioplay/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AudioService struct {
	cfg *config.Config
}

func NewAudioService(cfg *config.Config) *AudioService {
	return &AudioService{cfg: cfg}
}

type AudioResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration  int    `json:"duration"` // en secondes
	StreamURL string `json:"stream_url"`
	Thumbnail string `json:"thumbnail"`
	Source    string `json:"source"` // "soundcloud", "jamendo", etc.
}

// Search recherche des contenus audio sur toutes les plateformes
func (as *AudioService) Search(query string) ([]AudioResult, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	var results []AudioResult
	var errs []error

	// Canal pour collecter les résultats
	resultChan := make(chan []AudioResult)
	errChan := make(chan error)

	// Lancer les recherches en parallèle
	go func() {
		scResults, err := as.searchSoundCloud(query)
		if err != nil {
			errChan <- fmt.Errorf("SoundCloud: %v", err)
			return
		}
		resultChan <- scResults
	}()

	go func() {
		jamResults, err := as.searchJamendo(query)
		if err != nil {
			errChan <- fmt.Errorf("Jamendo: %v", err)
			return
		}
		resultChan <- jamResults
	}()

	// Collecter les résultats avec timeout
	timeout := time.After(5 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultChan:
			results = append(results, res...)
		case err := <-errChan:
			errs = append(errs, err)
		case <-timeout:
			errs = append(errs, errors.New("timeout on audio search"))
			break
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found (%v)", errs)
	}

	return results, nil
}

// searchSoundCloud recherche sur SoundCloud
func (as *AudioService) searchSoundCloud(query string) ([]AudioResult, error) {
	if as.cfg.SoundCloudClientID == "" {
		return nil, errors.New("SoundCloud client ID not configured")
	}

	apiURL := "https://api.soundcloud.com/tracks"
	params := url.Values{
		"q":          []string{query},
		"client_id":  []string{as.cfg.SoundCloudClientID},
		"limit":      []string{"5"},
		"streamable": []string{"true"},
	}

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var scTracks []struct {
		ID       int
		Title    string
		Duration int
		User     struct {
			Username string
		}
		ArtworkURL string
		StreamURL  string
	}

	if err := json.Unmarshal(body, &scTracks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	var results []AudioResult
	for _, track := range scTracks {
		if track.StreamURL == "" {
			continue
		}

		streamURL := fmt.Sprintf("%s?client_id=%s", track.StreamURL, as.cfg.SoundCloudClientID)
		thumbnail := track.ArtworkURL
		if thumbnail == "" {
			thumbnail = "https://picsum.photos/300/180?random=" + fmt.Sprint(track.ID)
		}

		results = append(results, AudioResult{
			ID:        fmt.Sprintf("sc-%d", track.ID),
			Title:     track.Title,
			Artist:    track.User.Username,
			Duration:  track.Duration / 1000, // conversion ms -> s
			StreamURL: streamURL,
			Thumbnail: thumbnail,
			Source:    "soundcloud",
		})
	}

	return results, nil
}

// searchJamendo recherche sur Jamendo
func (as *AudioService) searchJamendo(query string) ([]AudioResult, error) {
	if as.cfg.JamendoClientID == "" {
		return nil, errors.New("Jamendo client ID not configured")
	}

	apiURL := "https://api.jamendo.com/v3.0/tracks/"
	params := url.Values{
		"client_id": []string{as.cfg.JamendoClientID},
		"format":    []string{"json"},
		"limit":     []string{"5"},
		"search":    []string{query},
		"include":   []string{"musicinfo"},
	}

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var jamendoResponse struct {
		Headers map[string]interface{}
		Results []struct {
			ID       string
			Name     string
			Duration int
			Artist   string
			Audio    string
			Image    string
		}
	}

	if err := json.Unmarshal(body, &jamendoResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	var results []AudioResult
	for _, track := range jamendoResponse.Results {
		if !strings.HasSuffix(track.Audio, ".mp3") {
			continue
		}

		thumbnail := track.Image
		if thumbnail == "" {
			thumbnail = "https://picsum.photos/300/180?random=" + track.ID
		}

		results = append(results, AudioResult{
			ID:        fmt.Sprintf("jam-%s", track.ID),
			Title:     track.Name,
			Artist:    track.Artist,
			Duration:  track.Duration,
			StreamURL: track.Audio,
			Thumbnail: thumbnail,
			Source:    "jamendo",
		})
	}

	return results, nil
}

// GetAudioDetails récupère les détails d'un audio spécifique
func (as *AudioService) GetAudioDetails(audioID string) (*AudioResult, error) {
	parts := strings.Split(audioID, "-")
	if len(parts) < 2 {
		return nil, errors.New("invalid audio ID format")
	}

	source := parts[0]
	id := parts[1]

	switch source {
	case "sc":
		return as.getSoundCloudTrackDetails(id)
	case "jam":
		return as.getJamendoTrackDetails(id)
	default:
		return nil, errors.New("unknown audio source")
	}
}

func (as *AudioService) getSoundCloudTrackDetails(trackID string) (*AudioResult, error) {
	if as.cfg.SoundCloudClientID == "" {
		return nil, errors.New("SoundCloud client ID not configured")
	}

	apiURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s?client_id=%s", trackID, as.cfg.SoundCloudClientID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var track struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Duration int    `json:"duration"`
		User     struct {
			Username string `json:"username"`
		} `json:"user"`
		ArtworkURL string `json:"artwork_url"`
		StreamURL  string `json:"stream_url"`
	}

	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if track.StreamURL == "" {
		return nil, errors.New("track not streamable")
	}

	streamURL := fmt.Sprintf("%s?client_id=%s", track.StreamURL, as.cfg.SoundCloudClientID)
	thumbnail := track.ArtworkURL
	if thumbnail == "" {
		thumbnail = "https://picsum.photos/300/180?random=" + trackID
	}

	return &AudioResult{
		ID:        fmt.Sprintf("sc-%d", track.ID),
		Title:     track.Title,
		Artist:    track.User.Username,
		Duration:  track.Duration / 1000,
		StreamURL: streamURL,
		Thumbnail: thumbnail,
		Source:    "soundcloud",
	}, nil
}

func (as *AudioService) getJamendoTrackDetails(trackID string) (*AudioResult, error) {
	if as.cfg.JamendoClientID == "" {
		return nil, errors.New("Jamendo client ID not configured")
	}

	apiURL := fmt.Sprintf("https://api.jamendo.com/v3.0/tracks/?client_id=%s&format=json&id=%s",
		as.cfg.JamendoClientID, trackID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var response struct {
		Results []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Duration int    `json:"duration"`
			Artist   string `json:"artist_name"`
			Audio    string `json:"audio"`
			Image    string `json:"image"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Results) == 0 {
		return nil, errors.New("track not found")
	}

	track := response.Results[0]
	thumbnail := track.Image
	if thumbnail == "" {
		thumbnail = "https://picsum.photos/300/180?random=" + trackID
	}

	return &AudioResult{
		ID:        fmt.Sprintf("jam-%s", track.ID),
		Title:     track.Name,
		Artist:    track.Artist,
		Duration:  track.Duration,
		StreamURL: track.Audio,
		Thumbnail: thumbnail,
		Source:    "jamendo",
	}, nil
}
