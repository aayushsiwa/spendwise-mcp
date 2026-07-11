package config

import "os"

type Config struct {
	ServerName     string
	ServerVersion  string
	BackendBaseURL string
}

func Load() Config {
	return Config{
		ServerName:     getEnv("MCP_SERVER_NAME", "SpendWise MCP"),
		ServerVersion:  getEnv("MCP_SERVER_VERSION", "0.1.0"),
		BackendBaseURL: getEnv("SPENDWISE_BACKEND_BASE_URL", "http://localhost:8090/api/v1"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
