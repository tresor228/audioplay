package api

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	// Groupe pour les routes API
	apiGroup := router.Group("/api/v1")

	// Routes audio
	apiGroup.GET("/search/audio", SearchAudio)
	apiGroup.POST("/favorites/audio", AddAudioFavorite)
	apiGroup.GET("/favorites/audio", ListAudioFavorites)

	// Routes vidéo
	apiGroup.GET("/search/video", SearchVideo)
	apiGroup.POST("/favorites/video", AddVideoFavorite)
	apiGroup.GET("/favorites/video", ListVideoFavorites)

	// Servir les fichiers statiques pour le frontend
	router.Static("/static", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
}
