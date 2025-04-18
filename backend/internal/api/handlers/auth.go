package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	// TODO: Add user service dependency
	jwtSecret []byte
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		jwtSecret: []byte(jwtSecret),
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the register request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Check if user already exists
	// TODO: Hash password
	// TODO: Create user in database

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user from database
	// TODO: Compare passwords
	// TODO: Generate JWT token

	// Mock implementation for now
	token, refreshToken, expiresAt, err := h.generateTokens(1, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// TODO: Validate refresh token
	// TODO: Generate new tokens

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// generateTokens generates JWT tokens
func (h *AuthHandler) generateTokens(userID uint, email string) (token string, refreshToken string, expiresAt time.Time, err error) {
	// Set token expiration
	expiresAt = time.Now().Add(24 * time.Hour)

	// Create the token claims
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   expiresAt.Unix(),
	}

	// Create token with claims
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString(h.jwtSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Create refresh token (simplified for now)
	refreshToken = "refresh-token-placeholder"

	return token, refreshToken, expiresAt, nil
}

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
