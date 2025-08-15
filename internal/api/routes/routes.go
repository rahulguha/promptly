
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/rahulguha/promptly/api"
)

func RegisterRoutes(r *gin.Engine, h *api.Handler) {
    r.POST("/prompt", h.CreatePrompt)
    r.GET("/prompts", h.ListPrompts)
    r.GET("/prompt/search/:userRole/:llmRole", h.SearchPrompt)
    r.GET("/roles/:userRole", h.GetRoles)
}
