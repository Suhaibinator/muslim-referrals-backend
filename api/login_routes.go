package api

import (
	"github.com/Suhaibinator/muslim-referrals-backend/config"
	"net/http"

	"encoding/base64"
	"encoding/json"
)

func (hs *HttpServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code from the request
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	token, err := hs.service.GetTokenFromCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokenBytes, marshalErr := json.Marshal(token)
	if marshalErr != nil {
		http.Error(w, "Failed to marshal token: "+marshalErr.Error(), http.StatusInternalServerError)
		return
	}
	base64Token := base64.StdEncoding.EncodeToString(tokenBytes)
	_, newUser, errorRetrievingUserFromToken := hs.service.GetUserIdFromTokenDigest(r.Context(), base64Token)

	if errorRetrievingUserFromToken != nil {
		http.Error(w, "Failed to retrieve user from token: "+errorRetrievingUserFromToken.Error(), http.StatusInternalServerError)
		return
	}

	var redirectPath string
	if newUser {
		redirectPath = config.NEW_USER_SIGNUP_PATH
	} else {
		redirectPath = config.DEFAULT_LOGIN_PATH
	}

	// Set the auth cookie securely and redirect the user
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    base64Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, redirectPath, http.StatusFound)
}
