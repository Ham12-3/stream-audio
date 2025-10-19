package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	WebRTC   WebRTCConfig
	NATS     NATSConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type WebRTCConfig struct {
	ICEServers []string
	UDPPortMin int
	UDPPortMax int
}

type NATSConfig struct {
	URL     string
	Subject string
}

type ServicesConfig struct {
	ASRURL string
	TTSURL string
	LLMURL string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		WebRTC: WebRTCConfig{
			ICEServers: []string{
				getEnv("STUN_SERVER", "stun:stun.l.google.com:19302"),
			},
			UDPPortMin: getEnvInt("UDP_PORT_MIN", 10000),
			UDPPortMax: getEnvInt("UDP_PORT_MAX", 20000),
		},
		NATS: NATSConfig{
			URL:     getEnv("NATS_URL", "nats://localhost:4222"),
			Subject: getEnv("NATS_SUBJECT", "voice."),
		},
		Services: ServicesConfig{
			ASRURL: getEnv("ASR_URL", "localhost:50051"),
			TTSURL: getEnv("TTS_URL", "localhost:50052"),
			LLMURL: getEnv("LLM_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
