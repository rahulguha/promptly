package api

import (
	"fmt" // Added fmt import
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rahulguha/promptly/internal/tracking"
)

// UserTrackingHandler holds dependencies for user tracking operations.
type UserTrackingHandler struct {
	Tracker tracking.Tracker
}

// NewUserTrackingHandler creates a new UserTrackingHandler.
func NewUserTrackingHandler(tracker tracking.Tracker) *UserTrackingHandler {
	return &UserTrackingHandler{Tracker: tracker}
}

// TrackUserRequest represents the request body for tracking a user.
type TrackUserRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

// TrackUser handles the API endpoint for tracking a user.
func (h *UserTrackingHandler) TrackUser(c *gin.Context) {
	var req TrackUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := h.Tracker.UserExists(req.Email) // Check existence by email (GSI)
	if err != nil {
		fmt.Printf("Error checking user existence for %s: %v\n", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence", "details": err.Error()})
		return
	}

	if !exists {
		// Create new user record
		err := h.Tracker.CreateUserRecord(req.UserID, req.Email, req.Name, req.Timestamp)
		if err != nil {
			fmt.Printf("Error creating user record for %s: %v\n", req.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User record created"})
		return
	} else {
		// Update existing user record
		err := h.Tracker.UpdateUserRecord(req.UserID, req.Email, req.Name, req.Timestamp)
		if err != nil {
			fmt.Printf("Error updating user record for %s: %v\n", req.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user record", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User record updated"})
		return
	}
}

// TrackActivityRequest represents the request body for tracking an activity.
type TrackActivityRequest struct {
	UserID        string                 `json:"user_id" binding:"required"`
	Email         string                 `json:"email" binding:"required"`
	Timestamp     int64                  `json:"timestamp" binding:"required"`
	ActivityType  string                 `json:"activity_type" binding:"required"`
	ActivityResult string                 `json:"activity_result" binding:"required"`
	ActivityDetails map[string]interface{} `json:"activity_details"`
}

// TrackActivity handles the API endpoint for tracking an activity.
func (h *UserTrackingHandler) TrackActivity(c *gin.Context) {
	var req TrackActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Tracker.CreateActivityLog(
		req.UserID,
		req.Email,
		req.Timestamp,
		req.ActivityType,
		req.ActivityResult,
		req.ActivityDetails,
	)

	if err != nil {
		fmt.Printf("Error creating activity log for %s: %v\n", req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity log", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity log created"})
}