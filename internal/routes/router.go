package routes

import (
	"net/http"
	// "os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rahulguha/promptly/internal/api"
)

// RegisterProfileRoutes sets up the routes for profile management
func RegisterProfileRoutes(r *gin.RouterGroup, handler *ProfileHandler) {
	profiles := r.Group("/profiles")
	{
		profiles.GET("", handler.GetProfiles)
		profiles.GET("/:id", handler.GetProfile)
		profiles.POST("", handler.CreateProfile)
		profiles.PUT("/:id", handler.UpdateProfile)
		profiles.DELETE("/:id", handler.DeleteProfile)
	}
}

// RegisterRoutes sets up all the routes for the application
func RegisterRoutes(r *gin.Engine, handler *Handler) {
	// Configure session middleware
	store := cookie.NewStore([]byte(handler.Cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("promptly-session", store))

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5175"}, // Correct frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allow credentials
		MaxAge:           12 * time.Hour,
	}))

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Initialize the API handler with the config
		apiHandler := api.NewAPIHandler(handler.Cfg)

		// Profile routes
		profileHandler := &ProfileHandler{Store: handler.ProfileStore}
		RegisterProfileRoutes(v1, profileHandler)

		// Persona routes
		personas := v1.Group("/personas")
		{
			personas.GET("", handler.GetPersonas)
			personas.GET("/:id", handler.GetPersona)
			personas.POST("", handler.CreatePersona)
			personas.PUT("/:id", handler.UpdatePersona)
			personas.DELETE("/:id", handler.DeletePersona)
		}

		// Template routes
		templates := v1.Group("/templates")
		{
			templates.GET("", handler.GetTemplates)
			templates.GET("/:id", handler.GetTemplate)
			templates.POST("", handler.CreateTemplate)
			templates.PUT("/:id", handler.UpdateTemplate)
			templates.POST("/:id/version", handler.CreateTemplateVersion)
			templates.DELETE("/:id", handler.DeleteTemplate)
		}

		// Prompt routes
		prompts := v1.Group("/prompts")
		{
			prompts.GET("", handler.GetPrompts)
			prompts.GET("/:id", handler.GetPrompt)
			prompts.POST("", handler.CreatePrompt)
			prompts.PUT("/:id", handler.UpdatePrompt)
			prompts.DELETE("/:id", handler.DeletePrompt)
		}

		// Generate prompt from template
		v1.POST("/generate-prompt", handler.GeneratePrompt)

		// Auth routes
		auth := v1.Group("/api/auth")
		{
			auth.GET("/login", apiHandler.Login)
			auth.GET("/callback", apiHandler.Callback)
			auth.GET("/me", apiHandler.GetMe)
			auth.GET("/logout", apiHandler.Logout)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "promptly",
		})
	})
}

// NewRouter creates a new Gin router and registers the routes
func NewRouter(handler *Handler) *gin.Engine {
	r := gin.Default()
	RegisterRoutes(r, handler)
	return r
}