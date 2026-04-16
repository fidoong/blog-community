package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/internal/user/application"
)

// OAuthHandler handles OAuth2 flows.
type OAuthHandler struct {
	useCase      application.UseCase
	jwtSecret    string
	frontendURL  string
	githubOAuth  *oauth2.Config
	googleOAuth  *oauth2.Config
}

func NewOAuthHandler(useCase application.UseCase, jwtSecret, frontendURL, apiBaseURL, ghClientID, ghClientSecret, gClientID, gClientSecret string) *OAuthHandler {
	return newOAuthHandler(useCase, jwtSecret, frontendURL, apiBaseURL, ghClientID, ghClientSecret, gClientID, gClientSecret)
}

func NewOAuthHandlerFromConfig(useCase application.UseCase, cfg *configs.Config) *OAuthHandler {
	return newOAuthHandler(useCase, cfg.JWTSecret, cfg.FrontendURL, cfg.APIBaseURL, cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GoogleClientID, cfg.GoogleClientSecret)
}

func newOAuthHandler(useCase application.UseCase, jwtSecret, frontendURL, apiBaseURL, ghClientID, ghClientSecret, gClientID, gClientSecret string) *OAuthHandler {
	return &OAuthHandler{
		useCase:     useCase,
		jwtSecret:   jwtSecret,
		frontendURL: frontendURL,
		githubOAuth: &oauth2.Config{
			ClientID:     ghClientID,
			ClientSecret: ghClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  apiBaseURL + "/auth/oauth/github/callback",
			Scopes:       []string{"read:user", "user:email"},
		},
		googleOAuth: &oauth2.Config{
			ClientID:     gClientID,
			ClientSecret: gClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  apiBaseURL + "/auth/oauth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
}

type oauthURLResponse struct {
	AuthURL string `json:"authUrl"`
}

func (h *OAuthHandler) GetGitHubAuthURL(c *gin.Context) {
	// In production, generate and store state in Redis
	url := h.githubOAuth.AuthCodeURL("state")
	response.Success(c.Writer, oauthURLResponse{AuthURL: url})
}

func (h *OAuthHandler) GetGoogleAuthURL(c *gin.Context) {
	url := h.googleOAuth.AuthCodeURL("state")
	response.Success(c.Writer, oauthURLResponse{AuthURL: url})
}

func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Error(errors.ErrInvalidInput)
		return
	}

	token, err := h.githubOAuth.Exchange(context.Background(), code)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	// Fetch GitHub user info
	client := h.githubOAuth.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	// If email is empty, fetch emails
	email := ghUser.Email
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary && e.Verified {
						email = e.Email
						break
					}
				}
			}
		}
	}

	if email == "" {
		email = fmt.Sprintf("%d@github.local", ghUser.ID)
	}

	u, err := h.useCase.OAuthLoginOrRegister(c.Request.Context(), "github", fmt.Sprintf("%d", ghUser.ID), email, ghUser.Login, ghUser.AvatarURL)
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := auth.GenerateAccessToken(u.ID, u.Email, u.Username, u.Role, h.jwtSecret)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	c.Redirect(http.StatusFound, h.frontendURL+"/login?token="+accessToken)
}
