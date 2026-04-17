package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/errors"
)

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Error(errors.ErrInvalidInput)
		return
	}

	token, err := h.googleOAuth.Exchange(context.Background(), code)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	// Google userinfo endpoint
	client := h.googleOAuth.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}
	defer resp.Body.Close()

	var gUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	username := gUser.Name
	if username == "" {
		username = fmt.Sprintf("google_%s", gUser.ID)
	}

	u, err := h.useCase.OAuthLoginOrRegister(c.Request.Context(), "google", gUser.ID, gUser.Email, username, gUser.Picture)
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := auth.GenerateAccessToken(u.ID, u.Email, u.Username, u.Role, h.jwtSecret)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}
	_ = h.tokenStore.SaveRefreshToken(c.Request.Context(), u.ID, refreshToken, 7*24*time.Hour)

	c.Redirect(http.StatusFound, h.frontendURL+"/login?token="+accessToken+"&refreshToken="+refreshToken)
}
