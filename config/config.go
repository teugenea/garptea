package config

import (
	"bytes"
	"fmt"
	"os"

	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type ConfiKey string

const (
	OIDC_PROVIDER_URL     ConfiKey = "OIDC_PROVIDER_URL"
	JWKS_URL              ConfiKey = "JWKS_URL"
	AUTH_URL              ConfiKey = "AUTH_URL"
	OIDC_ACCESS_TOKEN_URL ConfiKey = "OIDC_ACCESS_TOKEN_URL"
	OIDC_CLIENT_ID        ConfiKey = "OIDC_CLIENT_ID"
	OIDC_CLIENT_SECRET    ConfiKey = "OIDC_CLIENT_SECRET"
	PUBLIC_CERT_FILE      ConfiKey = "PUBLIC_CERT_FILE"
	PORT                  ConfiKey = "PORT"
	TLS_ENABLED           ConfiKey = "TLS_ENABLED"
	TLS_CERT_FILE         ConfiKey = "TLS_CERT_FILE"
	TLS_KEY_FILE          ConfiKey = "TLS_KEY_FILE"
)

var (
	_mandatoryKeys = [...]ConfiKey{
		JWKS_URL,
		OIDC_PROVIDER_URL,
		AUTH_URL,
		OIDC_ACCESS_TOKEN_URL,
		OIDC_CLIENT_ID,
		OIDC_CLIENT_SECRET,
	}
)

func GetEnvVar(key ConfiKey) string {
	return os.Getenv(string(key))
}

func GetStringOrDefault(key ConfiKey, defaultValue string) string {
	value := os.Getenv(string(key))
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetBoolOrDefault(key ConfiKey, defaultValue bool) bool {
	strValue := GetEnvVar(key)
	if len(strValue) == 0 {
		return defaultValue
	}
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot parse bool config value (%s)", key))
		return defaultValue
	}
	return value
}

func ContatEnvVars(keys ...ConfiKey) string {
	var buffer bytes.Buffer
	for _, key := range keys {
		buffer.WriteString(GetEnvVar(key))
	}
	return buffer.String()
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

	for _, mandatoryKey := range _mandatoryKeys {
		if len(GetEnvVar(mandatoryKey)) == 0 {
			panic(fmt.Sprintf("Cannot load config with key '%s'", mandatoryKey))
		}
	}
}
