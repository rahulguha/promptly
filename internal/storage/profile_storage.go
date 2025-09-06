package storage

import "github.com/rahulguha/promptly/internal/models"

// ProfileStorage defines the interface for profile storage operations
type ProfileStorage interface {
	GetAllProfiles() ([]*models.Profile, error)
	GetProfileByID(id string) (*models.Profile, error)
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	DeleteProfile(id string) error
}
