package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"keyless-auth/domain"
	"keyless-auth/repository"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleHandler struct {
	oauthConfig *oauth2.Config
	repository  *repository.UserRepository
}

func NewGoogleHandler(repo *repository.UserRepository) *GoogleHandler {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/photoslibrary.readonly",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleHandler{
		oauthConfig: config,
		repository:  repo,
	}
}

func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (h *GoogleHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateState()
	url := h.oauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GoogleHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/auth/failure", http.StatusTemporaryRedirect)
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, "/auth/failure", http.StatusTemporaryRedirect)
		return
	}

	client := h.oauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Redirect(w, r, "/auth/failure", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Redirect(w, r, "/auth/failure", http.StatusTemporaryRedirect)
		return
	}

	user := &domain.User{
		ID:           userInfo.ID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if err := h.repository.SaveUser(user); err != nil {
		http.Redirect(w, r, "/auth/failure", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/auth/success", http.StatusTemporaryRedirect)
}
