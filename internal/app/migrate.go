//go:build migrate

package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	// Load .env file - try to find it smartly

	// 1. Try PROJECT_ROOT environment variable first
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		envPath := filepath.Join(projectRoot, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err == nil {
				log.Printf("Loaded .env from: %s", envPath)
				// Continue to check DB credentials
			}
		}
	}

	// 2. Walk up from current working directory to find .env
	if _, err := os.LookupEnv("DB_NAME"); err == false { // DB_NAME not set yet
		if wd, err := os.Getwd(); err == nil {
			for {
				envPath := filepath.Join(wd, ".env")
				if _, err := os.Stat(envPath); err == nil {
					if err := godotenv.Load(envPath); err == nil {
						log.Printf("Loaded .env from: %s", envPath)
						break
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
	}

	// 3. If nothing found, try current directory as last resort
	_ = godotenv.Load(".env")

	user, ok := os.LookupEnv("DB_USER")
	if !ok || len(user) == 0 {
		log.Fatalf("migrate: environment variable not declared: DB_USER")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatalf("migrate: environment variable not declared: DB_PASSWORD")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok || len(name) == 0 {
		log.Fatalf("migrate: environment variable not declared: DB_NAME")
	}

	databaseURL := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&multiStatements=true",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		name,
	)

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: mysql is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: mysql connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
