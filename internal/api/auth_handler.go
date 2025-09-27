package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	cognito_jwt "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rahulguha/promptly/internal/config"
	"golang.org/x/oauth2"
)

// APIHandler holds dependencies for API handlers.
type APIHandler struct {
	Cfg *config.Config
}

// JWTClaims represents the claims in our JWT token
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

// NewAPIHandler creates a new APIHandler.
func NewAPIHandler(cfg *config.Config) *APIHandler {
	return &APIHandler{Cfg: cfg}
}

// generateJWT creates a signed JWT token with user information
func (h *APIHandler) generateJWT(userID, email, name, picture string) (string, error) {
	claims := JWTClaims{
		UserID:  userID,
		Email:   email,
		Name:    name,
		Picture: picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hour expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "promptly-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.Cfg.SessionSecret))
}

// Auth handlers

// Login handles GET /auth/login
func (h *APIHandler) Login(c *gin.Context) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	state := base64.RawURLEncoding.EncodeToString(b)

	session := sessions.Default(c)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	cognitoDomain := h.Cfg.CognitoDomain
	clientID := h.Cfg.CognitoClientID
	redirectURI := h.Cfg.CognitoRedirectURI
	scopes := "openid profile email"

	authURL := fmt.Sprintf("https://%s/login?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		cognitoDomain, clientID, url.QueryEscape(redirectURI), state, url.QueryEscape(scopes))

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles GET /auth/callback
func (h *APIHandler) Callback(c *gin.Context) {
	session := sessions.Default(c)

	expectedState := session.Get("state")
	if c.Query("state") != expectedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	cognitoDomain := h.Cfg.CognitoDomain
	clientID := h.Cfg.CognitoClientID
	clientSecret := h.Cfg.CognitoClientSecret
	redirectURI := h.Cfg.CognitoRedirectURI

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/oauth2/authorize", cognitoDomain),
			TokenURL: fmt.Sprintf("https://%s/oauth2/token", cognitoDomain),
		},
	}

	token, err := conf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token", "details": err.Error()})
		return
	}

	idTokenRaw, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id_token not found in response"})
		return
	}

	// Parse the ID token to extract claims.
	// NOTE: We are skipping signature validation here because we just received the token
	// directly from Cognito over a secure channel. For a production environment, you
	// should implement full validation of the token's signature and claims.
	idToken, err := cognito_jwt.ParseString(idTokenRaw, cognito_jwt.WithVerify(false))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse id_token", "details": err.Error()})
		return
	}

	userID := idToken.Subject()
	var email, name, picture string
	if e, ok := idToken.Get("email"); ok {
		email = e.(string)
	}
	if n, ok := idToken.Get("name"); ok {
		name = n.(string)
	}
	if p, ok := idToken.Get("picture"); ok {
		picture = p.(string)
	}

	// Generate JWT token
	jwtToken, err := h.generateJWT(userID, email, name, picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token", "details": err.Error()})
		return
	}

	// Get current timestamp in Unix milliseconds
	timestamp := time.Now().UnixMilli()

	// Call the user tracking API asynchronously
	go func() {
		trackURL := fmt.Sprintf("%s/v1/track/users", h.Cfg.FrontendURL) // Assuming FrontendURL is the base URL for the API
		payload := map[string]interface{}{
			"user_id":   userID,
			"email":     email,
			"name":      name,
			"timestamp": timestamp,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshalling user tracking payload: %v\n", err)
			return
		}

		resp, err := http.Post(trackURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("Error calling user tracking API: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("User tracking API returned non-OK status: %d\n", resp.StatusCode)
		}
	}()

	// Redirect to frontend with JWT token as URL parameter
	redirectURL := fmt.Sprintf("%s/auth/success?token=%s", h.Cfg.FrontendURL, jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GetMe handles GET /auth/me
func (h *APIHandler) GetMe(c *gin.Context) {
	// Extract JWT token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check Bearer prefix
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	tokenString := authHeader[len(bearerPrefix):]

	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.Cfg.SessionSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"name":    claims.Name,
		"picture": claims.Picture,
	})
}

// Logout handles GET /auth/logout
func (h *APIHandler) Logout(c *gin.Context) {
	// For JWT-based auth, logout is handled client-side by discarding the token
	// Server-side logout would require token blacklisting, which we're not implementing here
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}