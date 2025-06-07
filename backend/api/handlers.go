package api

import (
	"audioplay/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	audioService   *services.AudioService
	youtubeService *services.YouTubeService
)

func InitServices(audio *services.AudioService, youtube *services.YouTubeService) {
	audioService = audio
	youtubeService = youtube
}

func SearchAudio(c *gin.Context) {
	query := c.Query("q")
	results, err := audioService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func SearchVideo(c *gin.Context) {
	query := c.Query("q")
	results, err := youtubeService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// Ajouter cette nouvelle fonction handler
func GetAudioDetails(c *gin.Context) {
	audioID := c.Param("id")
	result, err := audioService.GetAudioDetails(audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Dans api/handlers.go
func AddAudioFavorite(c *gin.Context) {
	// Implémentation pour ajouter un audio aux favoris
	c.JSON(http.StatusOK, gin.H{"message": "Audio added to favorites"})
}

func ListAudioFavorites(c *gin.Context) {
	// Implémentation pour lister les audios favoris
	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}

func AddVideoFavorite(c *gin.Context) {
	// Implémentation pour ajouter une vidéo aux favoris
	c.JSON(http.StatusOK, gin.H{"message": "Video added to favorites"})
}

func ListVideoFavorites(c *gin.Context) {
	// Implémentation pour lister les vidéos favorites
	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}
