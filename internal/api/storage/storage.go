package storage

import "github.com/rahulguha/promptly/models"

type StorageProvider interface {
    SavePrompt(p models.Prompt) error
    GetPrompts() ([]models.Prompt, error)
    SearchPrompt(userRole, llmRole string) ([]models.Prompt, error)
    GetRoles(userRole string) (models.RoleResponse, error)
}

