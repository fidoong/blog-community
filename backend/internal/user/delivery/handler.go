package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/pkg/validator"
	"github.com/blog/blog-community/internal/user/application"
	"github.com/blog/blog-community/internal/user/domain"
)

// UserHandler handles HTTP requests for user service.
type UserHandler struct {
	useCase    application.UseCase
	jwtSecret  string
	tokenStore auth.TokenStore
}

func NewUserHandler(useCase application.UseCase, jwtSecret string, store auth.TokenStore) *UserHandler {
	return &UserHandler{useCase: useCase, jwtSecret: jwtSecret, tokenStore: store}
}

func NewUserHandlerFromConfig(useCase application.UseCase, cfg *configs.Config, store auth.TokenStore) *UserHandler {
	return NewUserHandler(useCase, cfg.JWTSecret, store)
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email,max=128"`
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
	ExpiresIn    int        `json:"expiresIn"`
	User         userOutput `json:"user"`
}

type userOutput struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
	Role      string `json:"role"`
}

type updateProfileRequest struct {
	Username  string `json:"username" validate:"omitempty,max=64"`
	AvatarURL string `json:"avatarUrl" validate:"omitempty,max=512,url"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	u, err := h.useCase.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	h.respondWithTokens(c, u)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	u, err := h.useCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	h.respondWithTokens(c, u)
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	userID, err := h.tokenStore.GetUserIDByRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(errors.ErrUnauthorized)
		return
	}

	// Delete old refresh token
	_ = h.tokenStore.DeleteRefreshToken(c.Request.Context(), req.RefreshToken)

	u, err := h.useCase.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	h.respondWithTokens(c, u)
}

func (h *UserHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	if err := h.tokenStore.DeleteRefreshToken(c.Request.Context(), req.RefreshToken); err != nil {
		// Ignore error, token might already be expired
	}

	c.JSON(http.StatusOK, gin.H{"code": "OK", "message": "success"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	u, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toUserOutput(u))
}

func (h *UserHandler) GetMe(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	u, err := h.useCase.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toUserOutput(u))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	u, err := h.useCase.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	if req.Username != "" {
		u.Username = req.Username
	}
	if req.AvatarURL != "" {
		u.AvatarURL = req.AvatarURL
	}

	if err := h.useCase.Update(c.Request.Context(), u); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toUserOutput(u))
}

func (h *UserHandler) respondWithTokens(c *gin.Context, u *domain.User) {
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

	if err := h.tokenStore.SaveRefreshToken(c.Request.Context(), u.ID, refreshToken, 7*24*time.Hour); err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	response.Success(c.Writer, authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		User:         toUserOutput(u),
	})
}

func toUserOutput(u *domain.User) userOutput {
	return userOutput{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
	}
}
