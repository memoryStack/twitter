package initializers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func LoadEnv(environment string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	path := filepath.Join(wd, ".env."+environment)
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("failed to load env file %s: %v", path, err)
	}
	log.Printf("loaded env file: %s", path)
}
