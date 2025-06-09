// backend/api/handlers.go
package api

import (
	"audioplay/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	audioService    *services.AudioService
	youtubeService  *services.YouTubeService
	playlistService *services.PlaylistService
)

func InitServices(audio *services.AudioService, youtube *services.YouTubeService, playlist *services.PlaylistService) {
	audioService = audio
	youtubeService = youtube
	playlistService = playlist
}

// [Les autres handlers existants restent les mêmes]

func PlayMedia(c *gin.Context) {
	var request struct {
		MediaID string `json:"media_id"`
		Type    string `json:"type"` // "audio" ou "video"
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var media interface{}
	var err error
	var source string

	switch request.Type {
	case "audio":
		media, err = audioService.GetAudioDetails(request.MediaID)
		source = media.(*services.AudioResult).Source
	case "video":
		media, err = youtubeService.GetVideoDetails(request.MediaID)
		source = "youtube"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Générer une playlist automatique
	playlist, err := playlistService.GenerateAutoPlaylist(request.MediaID, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"media":    media,
		"playlist": playlist,
	})
}
func SearchAudio(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	results, err := audioService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
func AddAudioFavorite(c *gin.Context) {
	var request struct {
		AudioID string `json:"audio_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := audioService.AddFavorite(request.AudioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
func ListAudioFavorites(c *gin.Context) {
	favorites, err := audioService.ListFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
func GetAudioDetails(c *gin.Context) {
	audioID := c.Param("id")
	if audioID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio ID is required"})
		return
	}

	audio, err := audioService.GetAudioDetails(audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, audio)
}
func SearchVideo(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	results, err := youtubeService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
func AddVideoFavorite(c *gin.Context) {
	var request struct {
		VideoID string `json:"video_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := youtubeService.AddFavorite(request.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
func ListVideoFavorites(c *gin.Context) {
	favorites, err := youtubeService.ListFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
func GetVideoDetails(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video ID is required"})
		return
	}

	video, err := youtubeService.GetVideoDetails(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}
