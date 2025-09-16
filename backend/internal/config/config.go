package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server Configuration
	Environment string
	Port        string

	// Google Cloud Configuration
	GoogleCloudProjectID         string
	GoogleCloudRegion            string
	GoogleApplicationCredentials string

	// Firebase Configuration
	FirebaseProjectID               string
	FirebasePrivateKeyID            string
	FirebasePrivateKey              string
	FirebaseClientEmail             string
	FirebaseClientID                string
	FirebaseAuthURI                 string
	FirebaseTokenURI                string
	FirebaseAuthProviderX509CertURL string
	FirebaseClientX509CertURL       string

	// Gemini AI Configuration
	GeminiAPIKey string

	// External APIs
	GoogleMapsAPIKey string
	WeatherAPIKey    string

	// JWT Configuration
	JWTSecret     string
	JWTExpiration int

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int
}

func Load() *Config {
	return &Config{
		// Server
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),

		// Google Cloud
		GoogleCloudProjectID:         getEnv("GOOGLE_CLOUD_PROJECT_ID", ""),
		GoogleCloudRegion:            getEnv("GOOGLE_CLOUD_REGION", "us-central1"),
		GoogleApplicationCredentials: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),

		// Firebase
		FirebaseProjectID:               getEnv("FIREBASE_PROJECT_ID", ""),
		FirebasePrivateKeyID:            getEnv("FIREBASE_PRIVATE_KEY_ID", ""),
		FirebasePrivateKey:              getEnv("FIREBASE_PRIVATE_KEY", ""),
		FirebaseClientEmail:             getEnv("FIREBASE_CLIENT_EMAIL", ""),
		FirebaseClientID:                getEnv("FIREBASE_CLIENT_ID", ""),
		FirebaseAuthURI:                 getEnv("FIREBASE_AUTH_URI", ""),
		FirebaseTokenURI:                getEnv("FIREBASE_TOKEN_URI", ""),
		FirebaseAuthProviderX509CertURL: getEnv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL", ""),
		FirebaseClientX509CertURL:       getEnv("FIREBASE_CLIENT_X509_CERT_URL", ""),

		// Gemini AI
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),

		// External APIs
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		WeatherAPIKey:    getEnv("WEATHER_API_KEY", ""),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24), // hours

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 3600), // seconds
	}
}

var config *Config

// GetConfig returns the current configuration, loading it if necessary
func GetConfig() *Config {
	if config == nil {
		config = Load()
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
