package auth

import (
	"encoding/json"
	"net/url"

	"fmt"

	"garptea/config"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ROLE_USER string = "garptea/user"
)

type CodeRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetTokenByAccessCode(code string) string {
	body := CodeRequest{
		Code:         code,
		ClientId:     config.GetEnvVar(config.OIDC_CLIENT_ID),
		ClientSecret: config.GetEnvVar(config.OIDC_CLIENT_SECRET),
		GrantType:    "authorization_code",
	}
	oidcUrl := config.GetStringOrEmpty(config.OIDC_PROVIDER_URL)
	oidcAccessTokenUrl := config.GetStringOrEmpty(config.OIDC_ACCESS_TOKEN_URL)
	req := fiber.Post(oidcUrl + oidcAccessTokenUrl)
	jsonBody, _ := json.Marshal(body)
	req.Body(jsonBody)
	_, rawResp, _ := req.Bytes()
	token := TokenResponse{}
	json.Unmarshal(rawResp, &token)
	return token.AccessToken
}

func GetLoginUrl() string {
	oidcUrl := config.GetStringOrEmpty(config.OIDC_PROVIDER_URL)
	oidcAuthUrl := config.GetStringOrEmpty(config.AUTH_URL)
	oidcClientId := config.GetStringOrEmpty(config.OIDC_CLIENT_ID)
	return fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid&state=loggedin",
		oidcUrl+oidcAuthUrl,
		oidcClientId,
		url.QueryEscape("https://psyduck.home:3000/token"),
	)
}

func GetJwksUrl() string {
	oidcUrl := config.GetStringOrEmpty(config.OIDC_PROVIDER_URL)
	oidcJwksUrl := config.GetStringOrEmpty(config.JWKS_URL)
	return oidcUrl + oidcJwksUrl
}

func ParseJwtToken(token string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		s, _ := os.ReadFile(config.GetStringOrEmpty(config.PUBLIC_CERT_FILE))
		return jwt.ParseRSAPublicKeyFromPEM([]byte(s))
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}
