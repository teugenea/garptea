package auth

import (
	"net/url"

	"fmt"

	"garptea/config"

	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

const (
	ROLE_USER string = "garptea/user_group"
)

var (
	casdoorClient *casdoorsdk.Client
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

func GetTokenByAccessCode(code string, state string) string {
	resp, _ := getClietn().GetOAuthToken(code, state)
	return resp.AccessToken
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

func ParseJwtToken(token string) (*casdoorsdk.Claims, error) {
	claims, _ := getClietn().ParseJwtToken(token)
	user, _ := getClietn().GetUserByUserId(claims.User.Id)
	session, _ := getClietn().GetSession(claims.User.Name, "garptea-app")
	if session == nil {

	}
	if user == nil {

	}
	if user.IsForbidden {

	}
	return claims, nil
}

func getClietn() *casdoorsdk.Client {
	if casdoorClient == nil {
		casdoorClient = createClient()
	}
	return casdoorClient
}

func createClient() *casdoorsdk.Client {
	cert, err := os.ReadFile(config.GetStringOrEmpty(config.PUBLIC_CERT_FILE))
	if err != nil {
		panic(err)
	}
	return casdoorsdk.NewClient(
		config.GetStringOrEmpty(config.OIDC_PROVIDER_URL),
		config.GetEnvVar(config.OIDC_CLIENT_ID),
		config.GetEnvVar(config.OIDC_CLIENT_SECRET),
		string(cert),
		"garptea",
		"garptea-app",
	)
}
