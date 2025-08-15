package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/rahulguha/promptly/models"
    "github.com/rahulguha/promptly/storage"
)

type Handler struct {
    Store storage.StorageProvider
}

func (h *Handler) CreatePrompt(c *gin.Context) {
    var prompt models.Prompt
    if err := c.ShouldBindJSON(&prompt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.Store.SavePrompt(prompt); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, prompt)
}

func (h *Handler) ListPrompts(c *gin.Context) {
    prompts, err := h.Store.GetPrompts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, prompts)
}

func (h *Handler) SearchPrompt(c *gin.Context) {
    userRole := c.Param("userRole")
    llmRole := c.Param("llmRole")
    results, err := h.Store.SearchPrompt(userRole, llmRole)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, results)
}

func (h *Handler) GetRoles(c *gin.Context) {
    userRole := c.Param("userRole")
    resp, err := h.Store.GetRoles(userRole)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}
