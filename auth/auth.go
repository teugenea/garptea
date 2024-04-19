package auth

import (
	"context"
	"errors"
	"net/url"
	"time"

	"fmt"

	"garptea/config"

	"os"

	"github.com/allegro/bigcache/v3"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"golang.org/x/oauth2"
)

var (
	casdoorClient *casdoorsdk.Client
	userCache, _  = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
)

type CodeRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func GetTokenByAccessCode(code string, state string) (*oauth2.Token, error) {
	return getClient().GetOAuthToken(code, state)
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

func ParseJwtToken(token string) (*casdoorsdk.Claims, error) {
	claims, parseErr := getClient().ParseJwtToken(token)
	if parseErr != nil {
		return nil, parseErr
	}
	if validateErr := ValidateUserAndSession(token, claims); validateErr != nil {
		return nil, validateErr
	}
	return claims, nil
}

func ValidateUserAndSession(token string, claims *casdoorsdk.Claims) error {
	if _, cacheErr := userCache.Get(claims.User.Id); cacheErr == nil {
		return nil
	}
	intrToken, err := getClient().IntrospectToken(token, "access_token")
	if err != nil {
		return err
	}
	if !intrToken.Active {
		return errors.New("token is not active")
	}
	user, err := getClient().GetUser(claims.User.Name)
	if err != nil {
		return err
	}
	if user.IsDeleted || user.IsForbidden {
		return errors.New("user is forbidden or deleted")
	}
	userCache.Set(claims.User.Id, []byte(claims.User.Name))
	return nil
}

func ResetUserValidation(id string) {
	userCache.Delete(id)
}

func getClient() *casdoorsdk.Client {
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
		config.GetStringOrEmpty(config.OIDC_ORGANIZATION),
		config.GetStringOrEmpty(config.OIDC_APP_NAME),
	)
}
