package api

import (
	"audioplay/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchAudio(c *gin.Context) {
	query := c.Query("q")

	// Appeler le service audio
	results, err := services.SearchAudioContent(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SearchVideo(c *gin.Context) {
	query := c.Query("q")

	// Appeler le service YouTube
	results, err := services.SearchYouTubeContent(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// TODO: Implémentez ces handlers pour les favoris
func AddAudioFavorite(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func ListAudioFavorites(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func AddVideoFavorite(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func ListVideoFavorites(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}
