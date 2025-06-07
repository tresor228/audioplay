package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	YouTubeAPIKey      string
	SoundCloudClientID string
	JamendoClientID    string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		YouTubeAPIKey:      os.Getenv("YOUTUBE_API_KEY"),
		SoundCloudClientID: os.Getenv("SOUNDCLOUD_CLIENT_ID"),
		JamendoClientID:    os.Getenv("JAMENDO_CLIENT_ID"),
	}
}
