package api

import (
	"net/http"
	"ops_service/internal/auth"
	"ops_service/internal/model"
	"ops_service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService  *service.UserService
	tokenService *auth.TokenService
}

func NewAuthHandler(userService *service.UserService, tokenService *auth.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenService.GenerateToken(userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

func (h *AuthHandler) DummyLogin(c *gin.Context) {
	role := c.Query("role")
	if role != string(model.RoleClient) && role != string(model.RoleModerator) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role, must be 'client' or 'moderator'"})
		return
	}

	token, err := h.tokenService.GenerateDummyToken(model.UserRole(role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
