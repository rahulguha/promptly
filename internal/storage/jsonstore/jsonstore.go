package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
	// "github.com/rahulguha/promptly/internal/storage"
)

type FileStorage struct {
	filePath         string
	templatesPath    string
	personasPath     string
	mutex            sync.RWMutex
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fs := &FileStorage{
		filePath:      filePath,
		templatesPath: filepath.Join(filepath.Dir(filePath), "prompt_template.json"),
		personasPath:  filepath.Join(filepath.Dir(filePath), "persona.json"),
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Create empty prompts file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := fs.save([]models.Prompt{}); err != nil {
			return nil, fmt.Errorf("failed to create initial file: %w", err)
		}
	}
	
	// Create empty templates file if it doesn't exist
	if _, err := os.Stat(fs.templatesPath); os.IsNotExist(err) {
		if err := fs.saveTemplates([]models.PromptTemplate{}); err != nil {
			return nil, fmt.Errorf("failed to create initial templates file: %w", err)
		}
	}
	
	// Create empty personas file if it doesn't exist
	if _, err := os.Stat(fs.personasPath); os.IsNotExist(err) {
		if err := fs.savePersonas([]models.Persona{}); err != nil {
			return nil, fmt.Errorf("failed to create initial personas file: %w", err)
		}
	}
	
	return fs, nil
}

func (fs *FileStorage) load() ([]models.Prompt, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, err
	}
	
	var prompts []models.Prompt
	if len(data) == 0 {
		return prompts, nil
	}
	
	if err := json.Unmarshal(data, &prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

func (fs *FileStorage) save(prompts []models.Prompt) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := json.MarshalIndent(prompts, "", "  ")
	if err != nil {
		return err
	}
	
	if err := os.MkdirAll(filepath.Dir(fs.filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	return os.WriteFile(fs.filePath, data, 0644)
}


// Storage interface methods (for HTTP CRUD operations)
func (fs *FileStorage) GetAll() ([]*models.Prompt, error) {
	prompts, err := fs.load()
	if err != nil {
		return nil, err
	}
	
	// Convert []models.Prompt to []*models.Prompt
	result := make([]*models.Prompt, len(prompts))
	for i := range prompts {
		result[i] = &prompts[i]
	}
	return result, nil
}

func (fs *FileStorage) GetByID(id uuid.UUID) (*models.Prompt, error) {
	prompts, err := fs.load()
	if err != nil {
		return nil, err
	}
	
	for _, prompt := range prompts {
		if prompt.ID == id {
			return &prompt, nil
		}
	}
	return nil, fmt.Errorf("prompt with ID %s not found", id)
}

func (fs *FileStorage) Create(prompt *models.Prompt) (*models.Prompt, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	// Load without additional locking since we already have the write lock
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, err
	}
	
	var prompts []models.Prompt
	if len(data) > 0 {
		if err := json.Unmarshal(data, &prompts); err != nil {
			return nil, err
		}
	}
	
	// Generate new ID if not set
	if prompt.ID == uuid.Nil {
		prompt.ID = uuid.New()
	}
	
	prompts = append(prompts, *prompt)
	
	// Save without additional locking
	jsonData, err := json.MarshalIndent(prompts, "", "  ")
	if err != nil {
		return nil, err
	}
	
	if err := os.WriteFile(fs.filePath, jsonData, 0644); err != nil {
		return nil, err
	}
	
	return prompt, nil
}

func (fs *FileStorage) Update(prompt *models.Prompt) (*models.Prompt, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, err
	}
	
	var prompts []models.Prompt
	if len(data) > 0 {
		if err := json.Unmarshal(data, &prompts); err != nil {
			return nil, err
		}
	}
	
	for i, p := range prompts {
		if p.ID == prompt.ID {
			prompts[i] = *prompt
			
			jsonData, err := json.MarshalIndent(prompts, "", "  ")
			if err != nil {
				return nil, err
			}
			
			if err := os.WriteFile(fs.filePath, jsonData, 0644); err != nil {
				return nil, err
			}
			
			return prompt, nil
		}
	}
	
	return nil, fmt.Errorf("prompt with ID %s not found", prompt.ID)
}

func (fs *FileStorage) Delete(id uuid.UUID) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err
	}
	
	var prompts []models.Prompt
	if len(data) > 0 {
		if err := json.Unmarshal(data, &prompts); err != nil {
			return err
		}
	}
	
	for i, prompt := range prompts {
		if prompt.ID == id {
			// Remove the prompt from slice
			prompts = append(prompts[:i], prompts[i+1:]...)
			
			jsonData, err := json.MarshalIndent(prompts, "", "  ")
			if err != nil {
				return err
			}
			
			return os.WriteFile(fs.filePath, jsonData, 0644)
		}
	}
	
	return fmt.Errorf("prompt with ID %s not found", id)
}

// Template storage methods

func (fs *FileStorage) loadTemplates() ([]models.PromptTemplate, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	
	data, err := os.ReadFile(fs.templatesPath)
	if err != nil {
		return nil, err
	}
	
	var templates []models.PromptTemplate
	if len(data) == 0 {
		return templates, nil
	}
	
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (fs *FileStorage) saveTemplates(templates []models.PromptTemplate) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(fs.templatesPath, data, 0644)
}

func (fs *FileStorage) GetAllTemplates() ([]*models.PromptTemplate, error) {
	templates, err := fs.loadTemplates()
	if err != nil {
		return nil, err
	}
	
	result := make([]*models.PromptTemplate, len(templates))
	for i := range templates {
		result[i] = &templates[i]
	}
	return result, nil
}

func (fs *FileStorage) GetTemplateByID(id uuid.UUID) (*models.PromptTemplate, error) {
	templates, err := fs.loadTemplates()
	if err != nil {
		return nil, err
	}
	
	for _, template := range templates {
		if template.ID == id {
			return &template, nil
		}
	}
	return nil, fmt.Errorf("template with ID %s not found", id)
}

func (fs *FileStorage) CreateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.templatesPath)
	if err != nil {
		return nil, err
	}
	
	var templates []models.PromptTemplate
	if len(data) > 0 {
		if err := json.Unmarshal(data, &templates); err != nil {
			return nil, err
		}
	}
	
	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}
	
	templates = append(templates, *template)
	
	jsonData, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return nil, err
	}
	
	if err := os.WriteFile(fs.templatesPath, jsonData, 0644); err != nil {
		return nil, err
	}
	
	return template, nil
}

func (fs *FileStorage) UpdateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.templatesPath)
	if err != nil {
		return nil, err
	}
	
	var templates []models.PromptTemplate
	if len(data) > 0 {
		if err := json.Unmarshal(data, &templates); err != nil {
			return nil, err
		}
	}
	
	for i, t := range templates {
		if t.ID == template.ID {
			templates[i] = *template
			
			jsonData, err := json.MarshalIndent(templates, "", "  ")
			if err != nil {
				return nil, err
			}
			
			if err := os.WriteFile(fs.templatesPath, jsonData, 0644); err != nil {
				return nil, err
			}
			
			return template, nil
		}
	}
	
	return nil, fmt.Errorf("template with ID %s not found", template.ID)
}

func (fs *FileStorage) DeleteTemplate(id uuid.UUID) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.templatesPath)
	if err != nil {
		return err
	}
	
	var templates []models.PromptTemplate
	if len(data) > 0 {
		if err := json.Unmarshal(data, &templates); err != nil {
			return err
		}
	}
	
	for i, template := range templates {
		if template.ID == id {
			templates = append(templates[:i], templates[i+1:]...)
			
			jsonData, err := json.MarshalIndent(templates, "", "  ")
			if err != nil {
				return err
			}
			
			return os.WriteFile(fs.templatesPath, jsonData, 0644)
		}
	}
	
	return fmt.Errorf("template with ID %s not found", id)
}

func (fs *FileStorage) GetTemplatesByPersonaID(personaID uuid.UUID) ([]*models.PromptTemplate, error) {
	templates, err := fs.loadTemplates()
	if err != nil {
		return nil, err
	}
	
	var result []*models.PromptTemplate
	for _, template := range templates {
		if template.PersonaID == personaID {
			result = append(result, &template)
		}
	}
	return result, nil
}

// Persona storage methods

func (fs *FileStorage) loadPersonas() ([]models.Persona, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	
	data, err := os.ReadFile(fs.personasPath)
	if err != nil {
		return nil, err
	}
	
	var personas []models.Persona
	if len(data) == 0 {
		return personas, nil
	}
	
	if err := json.Unmarshal(data, &personas); err != nil {
		return nil, err
	}
	return personas, nil
}

func (fs *FileStorage) savePersonas(personas []models.Persona) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := json.MarshalIndent(personas, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(fs.personasPath, data, 0644)
}

func (fs *FileStorage) GetAllPersonas() ([]*models.Persona, error) {
	personas, err := fs.loadPersonas()
	if err != nil {
		return nil, err
	}
	
	result := make([]*models.Persona, len(personas))
	for i := range personas {
		result[i] = &personas[i]
	}
	return result, nil
}

func (fs *FileStorage) GetPersonaByID(id uuid.UUID) (*models.Persona, error) {
	personas, err := fs.loadPersonas()
	if err != nil {
		return nil, err
	}
	
	for _, persona := range personas {
		if persona.ID == id {
			return &persona, nil
		}
	}
	return nil, fmt.Errorf("persona with ID %s not found", id)
}

func (fs *FileStorage) CreatePersona(persona *models.Persona) (*models.Persona, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.personasPath)
	if err != nil {
		return nil, err
	}
	
	var personas []models.Persona
	if len(data) > 0 {
		if err := json.Unmarshal(data, &personas); err != nil {
			return nil, err
		}
	}
	
	if persona.ID == uuid.Nil {
		persona.ID = uuid.New()
	}
	
	personas = append(personas, *persona)
	
	jsonData, err := json.MarshalIndent(personas, "", "  ")
	if err != nil {
		return nil, err
	}
	
	if err := os.WriteFile(fs.personasPath, jsonData, 0644); err != nil {
		return nil, err
	}
	
	return persona, nil
}

func (fs *FileStorage) UpdatePersona(persona *models.Persona) (*models.Persona, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.personasPath)
	if err != nil {
		return nil, err
	}
	
	var personas []models.Persona
	if len(data) > 0 {
		if err := json.Unmarshal(data, &personas); err != nil {
			return nil, err
		}
	}
	
	for i, p := range personas {
		if p.ID == persona.ID {
			personas[i] = *persona
			
			jsonData, err := json.MarshalIndent(personas, "", "  ")
			if err != nil {
				return nil, err
			}
			
			if err := os.WriteFile(fs.personasPath, jsonData, 0644); err != nil {
				return nil, err
			}
			
			return persona, nil
		}
	}
	
	return nil, fmt.Errorf("persona with ID %s not found", persona.ID)
}

func (fs *FileStorage) DeletePersona(id uuid.UUID) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	
	data, err := os.ReadFile(fs.personasPath)
	if err != nil {
		return err
	}
	
	var personas []models.Persona
	if len(data) > 0 {
		if err := json.Unmarshal(data, &personas); err != nil {
			return err
		}
	}
	
	for i, persona := range personas {
		if persona.ID == id {
			personas = append(personas[:i], personas[i+1:]...)
			
			jsonData, err := json.MarshalIndent(personas, "", "  ")
			if err != nil {
				return err
			}
			
			return os.WriteFile(fs.personasPath, jsonData, 0644)
		}
	}
	
	return fmt.Errorf("persona with ID %s not found", id)
}

// Close closes the storage (no-op for JSON storage)
func (fs *FileStorage) Close() error {
	return nil
}