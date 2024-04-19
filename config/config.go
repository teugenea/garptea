package config

import (
	"fmt"
	"os"

	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type ConfigKey string

const (
	PORT                  ConfigKey = "PORT"
	HOST                  ConfigKey = "HOST"
	OIDC_PROVIDER_URL     ConfigKey = "OIDC_PROVIDER_URL"
	JWKS_URL              ConfigKey = "JWKS_URL"
	AUTH_URL              ConfigKey = "AUTH_URL"
	OIDC_ACCESS_TOKEN_URL ConfigKey = "OIDC_ACCESS_TOKEN_URL"
	OIDC_CLIENT_ID        ConfigKey = "OIDC_CLIENT_ID"
	OIDC_CLIENT_SECRET    ConfigKey = "OIDC_CLIENT_SECRET"
	OIDC_APP_NAME         ConfigKey = "OIDC_APP_NAME"
	OIDC_ORGANIZATION     ConfigKey = "OIDC_ORGANIZATION"
	HOOK_SECRET           ConfigKey = "HOOK_SECRET"
	PUBLIC_CERT_FILE      ConfigKey = "PUBLIC_CERT_FILE"
	TLS_ENABLED           ConfigKey = "TLS_ENABLED"
	TLS_CERT_FILE         ConfigKey = "TLS_CERT_FILE"
	TLS_KEY_FILE          ConfigKey = "TLS_KEY_FILE"
)

var (
	mandatoryKeys = [...]ConfigKey{
		JWKS_URL,
		OIDC_PROVIDER_URL,
		AUTH_URL,
		OIDC_ACCESS_TOKEN_URL,
		OIDC_CLIENT_ID,
		OIDC_CLIENT_SECRET,
		HOST,
	}
)

func GetEnvVar(key ConfigKey) string {
	return os.Getenv(string(key))
}

func GetStringOrDefault(key ConfigKey, defaultValue string) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetStringOrEmpty(key ConfigKey) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		return ""
	}
	return value
}

func GetBoolOrDefault(key ConfigKey, defaultValue bool) bool {
	strValue := GetEnvVar(key)
	if len(strValue) == 0 {
		return defaultValue
	}
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot parse bool config value (%s), return default (%t)", key, defaultValue))
		return defaultValue
	}
	return value
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		err2 := godotenv.Load("../.env")
		if err2 != nil {
			log.Info("Cannot load .env file. System env variables will be used")
		}
	}
	log.SetLevel(log.LevelInfo)

	for _, mandatoryKey := range mandatoryKeys {
		if len(GetEnvVar(mandatoryKey)) == 0 {
			panic(fmt.Sprintf("Cannot load config with key '%s'", mandatoryKey))
		}
	}
}
