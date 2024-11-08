package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DatabasePath      string
	GoogleOauthConfig *oauth2.Config
)

const (
	NEW_USER_SIGNUP_PATH = "/app/new-user-signup"
	DEFAULT_LOGIN_PATH   = "/app/"

	Port = "80"

	OAuthRedirectPath = "/login"
)

var (
	OauthRedirectHost = "https://muslimreferrals.xyz"
)

func init() {

	redirectHostEnvVar := os.Getenv("OAUTH_REDIRECT_HOST") // http://localhost:3000/login or https://muslimreferrals.xyz/login
	if redirectHostEnvVar != "" {
		OauthRedirectHost = redirectHostEnvVar
	}

	DatabasePath = os.Getenv("SQLITE_DB_PATH")
	log.Println("Google redirect URL: ", os.Getenv("GOOGLE_REDIRECT_URL"))
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  OauthRedirectHost + OAuthRedirectPath,
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
