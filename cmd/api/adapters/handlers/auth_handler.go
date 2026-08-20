package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	service core.AuthService
	log     logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service core.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		log:     log,
	}
}

// SignUp godoc
// @Summary      Sign Up Private Network Operator
// @Description  Creates a new Private Network operator account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body  object{username=string,password=string}  true  "Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string "Invalid request body or credentials"
// @Failure      409  {object}  map[string]string "Username already exists"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /private-network/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Error("Failed to parse signup request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	err := h.service.SignUp(c.Request.Context(), body.Username, body.Password)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Login godoc
// @Summary      Login Private Network Operator
// @Description  Authenticates a Private Network operator and returns a JWT token in a cookie
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body  body  object{username=string,password=string}  true  "Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string "Invalid request body or credentials"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /private-network/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Error("Failed to parse login request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	tokenString, err := h.service.Login(c.Request.Context(), body.Username, body.Password)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	// Set JWT token in cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*core.TokenExpiryDays, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{})
}
