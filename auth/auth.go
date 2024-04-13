package auth

import (
	"encoding/json"
	"net/url"

	"fmt"

	"garptea/config"

	"github.com/gofiber/fiber/v2"
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
	oidcUrl := config.GetStringOrDefault(config.OIDC_PROVIDER_URL, "")
	oidcAccessTokenUrl := config.GetStringOrDefault(config.OIDC_ACCESS_TOKEN_URL, "")
	req := fiber.Post(oidcUrl + oidcAccessTokenUrl)
	jsonBody, _ := json.Marshal(body)
	req.Body(jsonBody)
	_, rawResp, _ := req.Bytes()
	token := TokenResponse{}
	json.Unmarshal(rawResp, &token)
	return token.AccessToken
}

func GetLoginUrl() string {
	oidcUrl := config.GetStringOrDefault(config.OIDC_PROVIDER_URL, "")
	oidcAuthUrl := config.GetStringOrDefault(config.AUTH_URL, "")
	oidcClientId := config.GetStringOrDefault(config.OIDC_CLIENT_ID, "")
	return fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid&state=loggedin",
		oidcUrl+oidcAuthUrl,
		oidcClientId,
		url.QueryEscape("http://psyduck.home:3000/token"),
	)
}

func GetJwksUrl() string {
	oidcUrl := config.GetStringOrDefault(config.OIDC_PROVIDER_URL, "")
	oidcJwksUrl := config.GetStringOrDefault(config.JWKS_URL, "")
	return oidcUrl + oidcJwksUrl
}
