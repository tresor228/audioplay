package api

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")

	// Routes audio
	apiGroup.GET("/search/audio", SearchAudio)
	apiGroup.POST("/favorites/audio", AddAudioFavorite)
	apiGroup.GET("/favorites/audio", ListAudioFavorites)
	apiGroup.GET("/audio/:id", GetAudioDetails)

	// Routes vidéo
	apiGroup.GET("/search/video", SearchVideo)
	apiGroup.POST("/favorites/video", AddVideoFavorite)
	apiGroup.GET("/favorites/video", ListVideoFavorites)

	// Nouvelle route pour la lecture
	apiGroup.POST("/play", PlayMedia)

	// Servir les fichiers statiques pour le frontend
	router.Static("/static", "../frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("../frontend/index.html")
	})
}
