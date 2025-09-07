package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rahulguha/promptly/internal/models"
	"github.com/rahulguha/promptly/internal/storage"
)

// ProfileHandler is a placeholder for profile-related routes.
// The actual storage is retrieved from the context in each handler.
type ProfileHandler struct{}

// GetProfiles handles GET /profiles
func (h *ProfileHandler) GetProfiles(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	profileStore := store.(storage.ProfileStorage)

	profiles, err := profileStore.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// GetProfile handles GET /profiles/:id
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	profileStore := store.(storage.ProfileStorage)

	id := c.Param("id")
	profile, err := profileStore.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// CreateProfile handles POST /profiles
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	profileStore := store.(storage.ProfileStorage)

	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := profileStore.CreateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// UpdateProfile handles PUT /profiles/:id
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	profileStore := store.(storage.ProfileStorage)

	id := c.Param("id")
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = id
	err := profileStore.UpdateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DeleteProfile handles DELETE /profiles/:id
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	profileStore := store.(storage.ProfileStorage)

	id := c.Param("id")
	err := profileStore.DeleteProfile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}