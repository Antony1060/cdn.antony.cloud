package env

import (
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"os"
)

type EnvConfig struct {
	Port string
	Token string
}

func New() *EnvConfig {
	// don't care about error
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Warn("Failed to load .env file")
	}
	return &EnvConfig{
		Port: getEnvWithDefault("PORT", "8080"),
		Token: getEnvWithDefault("TOKEN", "secret"),
	}
}

func getEnvWithDefault(variable, def string) string {
	env, exists := os.LookupEnv(variable)
	if exists {
		return env
	}
	return def
}