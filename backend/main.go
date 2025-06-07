package main

import (
	"audioplay/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configuration des routes
	api.SetupRoutes(router)

	// Démarrer le serveur
	router.Run(":8080")
}
