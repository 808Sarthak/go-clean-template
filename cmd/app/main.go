package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/app"
	"github.com/joho/godotenv"
)

// init runs before main and loads environment variables from .env
func init() {
	// Load .env file - try to find it smartly

	// 1. Try PROJECT_ROOT environment variable first
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		envPath := filepath.Join(projectRoot, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err == nil {
				log.Printf("Loaded .env from: %s", envPath)
				return
			}
		}
	}

	// 2. Walk up from current working directory to find .env
	if wd, err := os.Getwd(); err == nil {
		for {
			envPath := filepath.Join(wd, ".env")
			if _, err := os.Stat(envPath); err == nil {
				if err := godotenv.Load(envPath); err == nil {
					log.Printf("Loaded .env from: %s", envPath)
					return
				}
			}

			// Move to parent directory
			parent := filepath.Dir(wd)
			if parent == wd { // Reached root directory
				break
			}
			wd = parent
		}
	}

	// 3. If nothing found, try current directory as last resort
	_ = godotenv.Load(".env")
}

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
