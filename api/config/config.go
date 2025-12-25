package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port                    string
	FirebaseProjectID       string
	FirebaseCredentialsPath string
	InfluxDBURL             string
	InfluxDBToken           string
	InfluxDBOrg             string
	InfluxDBBucket          string
	AllowedOrigins          string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists (development only)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:                    getEnv("PORT", "3000"),
		FirebaseProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		InfluxDBURL:             getEnv("INFLUXDB_URL", ""),
		InfluxDBToken:           getEnv("INFLUXDB_TOKEN", ""),
		InfluxDBOrg:             getEnv("INFLUXDB_ORG", ""),
		InfluxDBBucket:          getEnv("INFLUXDB_BUCKET", "ergometer-workouts"),
		AllowedOrigins:          getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
