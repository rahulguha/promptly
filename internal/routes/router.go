package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all the routes for the application
func RegisterRoutes(r *gin.Engine, handler *Handler) {
	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// API v1 routes
	v1 := r.Group("/v1")
	{
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
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "promptly",
		})
	})
}