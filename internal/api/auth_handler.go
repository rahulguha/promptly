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
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rahulguha/promptly/internal/config"
	"golang.org/x/oauth2"
)

// APIHandler holds dependencies for API handlers.
type APIHandler struct {
	Cfg *config.Config
}

// NewAPIHandler creates a new APIHandler.
func NewAPIHandler(cfg *config.Config) *APIHandler {
	return &APIHandler{Cfg: cfg}
}

// Auth handlers

// Login handles GET /auth/login
func (h *APIHandler) Login(c *gin.Context) {
	fmt.Println("Login handler called")
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
	fmt.Println("Callback handler called")
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
	idToken, err := jwt.ParseString(idTokenRaw, jwt.WithVerify(false))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse id_token", "details": err.Error()})
		return
	}
	fmt.Printf("ID Token claims: %+v\n", idToken)
	// Store essential user info in the session
	session.Set("user_id", idToken.Subject())

	userID := idToken.Subject()
	var email string
	if e, ok := idToken.Get("email"); ok {
		email = e.(string)
		session.Set("email", email)
	}
	var name string
	if n, ok := idToken.Get("name"); ok {
		name = n.(string)
		session.Set("name", name)
	}
	if picture, ok := idToken.Get("picture"); ok {
		session.Set("picture", picture)
	}
	session.Set("authenticated", true)

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

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session", "details": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.Cfg.FrontendURL)
}

// GetMe handles GET /auth/me
func (h *APIHandler) GetMe(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")

	if authenticated == nil || !authenticated.(bool) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": session.Get("user_id"),
		"email":   session.Get("email"),
		"name":    session.Get("name"),
		"picture":    session.Get("picture"),
	})
}

// Logout handles GET /auth/logout
func (h *APIHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1}) // Expire the cookie immediately
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}