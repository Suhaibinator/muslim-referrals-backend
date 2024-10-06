package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DatabasePath      string
	GoogleOauthConfig *oauth2.Config
)

func init() {
	DatabasePath = os.Getenv("SQLITE_DB_PATH")
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

const (
	NEW_USER_SIGNUP_PATH = "/app/new-user-signup"
	DEFAULT_LOGIN_PATH   = "/app/"

	Port = "8090"
)
