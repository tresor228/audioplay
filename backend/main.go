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

	// Injecter les services dans les handlers
	api.InitServices(audioService, youtubeService)

	router := gin.Default()
	api.SetupRoutes(router)

	log.Println("Server running on :8080")
	router.Run(":8080")
}
