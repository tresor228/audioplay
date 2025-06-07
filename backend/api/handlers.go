package api

import (
	"net/http"
	"yourpackage/services"

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

// ... autres handlers pour les favoris
