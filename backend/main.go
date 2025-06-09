// backend/main.go
package main

import (
	"audioplay/api"
	"audioplay/config"
	"audioplay/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Initialiser les services
	audioService := services.NewAudioService(cfg)
	youtubeService, err := services.NewYouTubeService(cfg)
	if err != nil {
		log.Fatalf("Failed to create YouTube service: %v", err)
	}

	// Initialiser le service de playlist
	playlistService := services.NewPlaylistService(audioService, youtubeService)

	// Injecter les services dans les handlers
	api.InitServices(audioService, youtubeService, playlistService)

	router := gin.Default()
	api.SetupRoutes(router)

	log.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
