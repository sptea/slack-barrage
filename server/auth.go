package main

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"os"

	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	v2 "google.golang.org/api/oauth2/v2"
)

type GoogleAuthConfig struct {
	ClientID          string
	ClientSecret      string
	AuthorizeEndpoint string
	TokenEndpoint     string
}

var authConfig GoogleAuthConfig

func InitAuthConfig(logger *log.Logger) {
	authConfig.ClientID = os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	if authConfig.ClientID == "" {
		logger.Panic("Failed to import env: GOOGLE_AUTH_CLIENT_ID")
	}
	authConfig.ClientSecret = os.Getenv("GOOGLE_AUTH_CLIENT_SECRET")
	if authConfig.ClientSecret == "" {
		logger.Panic("Failed to import env: GOOGLE_AUTH_CLIENT_SECRET")
	}
	authConfig.AuthorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	authConfig.TokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"

	// Register UserInfo to use it on session
	gob.Register(&UserInfo{})

}

func GetConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authConfig.AuthorizeEndpoint,
			TokenURL: authConfig.TokenEndpoint,
		},
		Scopes:      []string{"openid", "email", "profile"},
		RedirectURL: "http://localhost:8080/callback",
	}

	return config
}

func AuthMethod(session *sessions.Session) string {
	state, _ := uuid.NewV4()
	session.Values["origState"] = state.String()

	config := GetConfig()
	url := config.AuthCodeURL(state.String())

	return url
}

func callbackMethod(session *sessions.Session, state string, code string, logger *log.Logger) error {
	config := GetConfig()
	context := context.Background()
	origState := session.Values["origState"].(string)

	if origState != state {
		return errors.New("callbackMethod: invalid state")
	}

	tok, err := config.Exchange(context, code)
	if err != nil {
		logger.Println(err)
		return err
	}

	if tok.Valid() == false {
		return errors.New("callbackMethod: invaild token")
	}

	service, _ := v2.New(config.Client(context, tok))

	// accessToken is used only this time
	tokenInfo, _ := service.Tokeninfo().AccessToken(tok.AccessToken).Context(context).Do()

	userInfo := &UserInfo{Email: tokenInfo.Email}

	session.Values["userInfo"] = userInfo

	return nil
}
