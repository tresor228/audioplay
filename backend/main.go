package main

import (
	"yourpackage/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configuration des routes
	api.SetupRoutes(router)

	// Démarrer le serveur
	router.Run(":8080")
}
