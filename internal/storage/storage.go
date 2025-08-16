package storage

import "github.com/rahulguha/promptly/internal/models"
import "github.com/google/uuid"

// StorageProvider defines the interface for prompt storage operations
type StorageProvider interface {
	SavePrompt(p models.Prompt) error
	GetPrompts() ([]models.Prompt, error)
	SearchPrompt(userRole, llmRole string) ([]models.Prompt, error)
	GetRoles(userRole string) (models.RoleResponse, error)
}

// Storage defines the interface for CRUD operations (used by HTTP handlers)
type Storage interface {
	// Persona operations
	GetAllPersonas() ([]*models.Persona, error)
	GetPersonaByID(id uuid.UUID) (*models.Persona, error)
	CreatePersona(persona *models.Persona) (*models.Persona, error)
	UpdatePersona(persona *models.Persona) (*models.Persona, error)
	DeletePersona(id uuid.UUID) error
	
	// Template operations
	GetAllTemplates() ([]*models.PromptTemplate, error)
	GetTemplateByID(id uuid.UUID) (*models.PromptTemplate, error)
	GetTemplatesByPersonaID(personaID uuid.UUID) ([]*models.PromptTemplate, error)
	CreateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error)
	UpdateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error)
	DeleteTemplate(id uuid.UUID) error
	
	// Prompt operations
	GetAll() ([]*models.Prompt, error)
	GetByID(id uuid.UUID) (*models.Prompt, error)
	Create(prompt *models.Prompt) (*models.Prompt, error)
	Update(prompt *models.Prompt) (*models.Prompt, error)
	Delete(id uuid.UUID) error
}