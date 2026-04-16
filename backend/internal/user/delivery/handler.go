package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/pkg/validator"
	"github.com/blog/blog-community/internal/user/application"
)

// UserHandler handles HTTP requests for user service.
type UserHandler struct {
	useCase   application.UseCase
	jwtSecret string
}

func NewUserHandler(useCase application.UseCase, jwtSecret string) *UserHandler {
	return &UserHandler{useCase: useCase, jwtSecret: jwtSecret}
}

func NewUserHandlerFromConfig(useCase application.UseCase, cfg *configs.Config) *UserHandler {
	return NewUserHandler(useCase, cfg.JWTSecret)
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

	token, err := auth.GenerateAccessToken(u.ID, u.Email, u.Username, u.Role, h.jwtSecret)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	response.Success(c.Writer, authResponse{
		AccessToken: token,
		ExpiresIn:   7200,
		User: userOutput{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			Role:      u.Role,
		},
	})
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

	token, err := auth.GenerateAccessToken(u.ID, u.Email, u.Username, u.Role, h.jwtSecret)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	response.Success(c.Writer, authResponse{
		AccessToken: token,
		ExpiresIn:   7200,
		User: userOutput{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			Role:      u.Role,
		},
	})
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

	response.Success(c.Writer, userOutput{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
	})
}
