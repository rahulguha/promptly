package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rahulguha/promptly/internal/models"
	"github.com/rahulguha/promptly/internal/storage"
)

// ProfileHandler contains the dependencies for profile HTTP handlers
type ProfileHandler struct {
	Store storage.ProfileStorage
}

// GetProfiles handles GET /profiles
func (h *ProfileHandler) GetProfiles(c *gin.Context) {
	profiles, err := h.Store.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// GetProfile handles GET /profiles/:id
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	id := c.Param("id")
	profile, err := h.Store.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// CreateProfile handles POST /profiles
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.CreateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// UpdateProfile handles PUT /profiles/:id
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = id
	err := h.Store.UpdateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DeleteProfile handles DELETE /profiles/:id
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	id := c.Param("id")
	err := h.Store.DeleteProfile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
