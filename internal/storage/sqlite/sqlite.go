package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	storage := &SQLiteStorage{
		db:   db,
		path: dbPath,
	}

	return storage, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// Persona operations
func (s *SQLiteStorage) CreatePersona(persona *models.Persona) (*models.Persona, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	persona.ID = uuid.New()
	
	query := `INSERT INTO personas (id, user_role_display, llm_role_display) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, persona.ID.String(), persona.UserRoleDisplay, persona.LLMRoleDisplay)
	if err != nil {
		return nil, fmt.Errorf("failed to create persona: %w", err)
	}

	return persona, nil
}

func (s *SQLiteStorage) GetAllPersonas() ([]*models.Persona, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, user_role_display, llm_role_display FROM personas ORDER BY created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query personas: %w", err)
	}
	defer rows.Close()

	var personas []*models.Persona
	for rows.Next() {
		var persona models.Persona
		var idStr string
		err := rows.Scan(&idStr, &persona.UserRoleDisplay, &persona.LLMRoleDisplay)
		if err != nil {
			return nil, fmt.Errorf("failed to scan persona: %w", err)
		}
		persona.ID = uuid.MustParse(idStr)
		personas = append(personas, &persona)
	}

	return personas, nil
}

func (s *SQLiteStorage) GetPersonaByID(id uuid.UUID) (*models.Persona, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var persona models.Persona
	var idStr string
	query := `SELECT id, user_role_display, llm_role_display FROM personas WHERE id = ?`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &persona.UserRoleDisplay, &persona.LLMRoleDisplay)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("persona not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get persona: %w", err)
	}

	persona.ID = uuid.MustParse(idStr)
	return &persona, nil
}

func (s *SQLiteStorage) UpdatePersona(persona *models.Persona) (*models.Persona, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE personas SET user_role_display = ?, llm_role_display = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := s.db.Exec(query, persona.UserRoleDisplay, persona.LLMRoleDisplay, persona.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update persona: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("persona not found")
	}

	return persona, nil
}

func (s *SQLiteStorage) DeletePersona(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM personas WHERE id = ?`
	result, err := s.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete persona: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("persona not found")
	}

	return nil
}

// Template operations
func (s *SQLiteStorage) CreateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	template.ID = uuid.New()
	template.Version = 1 // New templates start at version 1
	
	variablesJSON, err := json.Marshal(template.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `INSERT INTO prompt_templates (id, persona_id, version, template, variables) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, template.ID.String(), template.PersonaID.String(), template.Version, template.Template, string(variablesJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

func (s *SQLiteStorage) GetAllTemplates() ([]*models.PromptTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, persona_id, version, template, variables FROM prompt_templates ORDER BY created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.PromptTemplate
	for rows.Next() {
		var template models.PromptTemplate
		var idStr, personaIDStr, variablesJSON string
		err := rows.Scan(&idStr, &personaIDStr, &template.Version, &template.Template, &variablesJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		template.ID = uuid.MustParse(idStr)
		template.PersonaID = uuid.MustParse(personaIDStr)

		if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}

		templates = append(templates, &template)
	}

	return templates, nil
}

func (s *SQLiteStorage) GetTemplateByID(id uuid.UUID) (*models.PromptTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var template models.PromptTemplate
	var idStr, personaIDStr, variablesJSON string
	query := `SELECT id, persona_id, version, template, variables FROM prompt_templates WHERE id = ? ORDER BY version DESC LIMIT 1`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &personaIDStr, &template.Version, &template.Template, &variablesJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	template.ID = uuid.MustParse(idStr)
	template.PersonaID = uuid.MustParse(personaIDStr)

	if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
	}

	return &template, nil
}

func (s *SQLiteStorage) UpdateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update the current version in place
	variablesJSON, err := json.Marshal(template.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `UPDATE prompt_templates SET persona_id = ?, template = ?, variables = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND version = ?`
	result, err := s.db.Exec(query, template.PersonaID.String(), template.Template, string(variablesJSON), template.ID.String(), template.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("template version not found")
	}

	return template, nil
}

func (s *SQLiteStorage) CreateTemplateVersion(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the current max version for this template ID
	var maxVersion int
	versionQuery := `SELECT MAX(version) FROM prompt_templates WHERE id = ?`
	err := s.db.QueryRow(versionQuery, template.ID.String()).Scan(&maxVersion)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	// Create new version
	template.Version = maxVersion + 1
	
	variablesJSON, err := json.Marshal(template.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `INSERT INTO prompt_templates (id, persona_id, version, template, variables) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, template.ID.String(), template.PersonaID.String(), template.Version, template.Template, string(variablesJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create new template version: %w", err)
	}

	return template, nil
}

func (s *SQLiteStorage) DeleteTemplate(id uuid.UUID, version int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM prompt_templates WHERE id = ? AND version = ?`
	result, err := s.db.Exec(query, id.String(), version)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("template version not found")
	}

	return nil
}

// Prompt operations
func (s *SQLiteStorage) Create(prompt *models.Prompt) (*models.Prompt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prompt.ID = uuid.New()
	
	valuesJSON, err := json.Marshal(prompt.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal values: %w", err)
	}

	query := `INSERT INTO prompts (id, template_id, template_version, variable_values, content) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, prompt.ID.String(), prompt.TemplateID.String(), prompt.TemplateVersion, string(valuesJSON), prompt.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	return prompt, nil
}

func (s *SQLiteStorage) GetAll() ([]*models.Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, template_id, template_version, variable_values, content FROM prompts ORDER BY created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query prompts: %w", err)
	}
	defer rows.Close()

	var prompts []*models.Prompt
	for rows.Next() {
		var prompt models.Prompt
		var idStr, templateIDStr, valuesJSON string
		err := rows.Scan(&idStr, &templateIDStr, &prompt.TemplateVersion, &valuesJSON, &prompt.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prompt: %w", err)
		}

		prompt.ID = uuid.MustParse(idStr)
		prompt.TemplateID = uuid.MustParse(templateIDStr)

		if err := json.Unmarshal([]byte(valuesJSON), &prompt.Values); err != nil {
			return nil, fmt.Errorf("failed to unmarshal values: %w", err)
		}

		prompts = append(prompts, &prompt)
	}

	return prompts, nil
}

func (s *SQLiteStorage) GetByID(id uuid.UUID) (*models.Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var prompt models.Prompt
	var idStr, templateIDStr, valuesJSON string
	query := `SELECT id, template_id, template_version, variable_values, content FROM prompts WHERE id = ?`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &templateIDStr, &prompt.TemplateVersion, &valuesJSON, &prompt.Content)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("prompt not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	prompt.ID = uuid.MustParse(idStr)
	prompt.TemplateID = uuid.MustParse(templateIDStr)

	if err := json.Unmarshal([]byte(valuesJSON), &prompt.Values); err != nil {
		return nil, fmt.Errorf("failed to unmarshal values: %w", err)
	}

	return &prompt, nil
}

func (s *SQLiteStorage) Update(prompt *models.Prompt) (*models.Prompt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valuesJSON, err := json.Marshal(prompt.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal values: %w", err)
	}

	query := `UPDATE prompts SET template_id = ?, template_version = ?, variable_values = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := s.db.Exec(query, prompt.TemplateID.String(), prompt.TemplateVersion, string(valuesJSON), prompt.Content, prompt.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("prompt not found")
	}

	return prompt, nil
}

func (s *SQLiteStorage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM prompts WHERE id = ?`
	result, err := s.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("prompt not found")
	}

	return nil
}

func (s *SQLiteStorage) GetTemplatesByPersonaID(personaID uuid.UUID) ([]*models.PromptTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, persona_id, version, template, variables FROM prompt_templates WHERE persona_id = ? ORDER BY created_at`
	rows, err := s.db.Query(query, personaID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query templates by persona: %w", err)
	}
	defer rows.Close()

	var templates []*models.PromptTemplate
	for rows.Next() {
		var template models.PromptTemplate
		var idStr, personaIDStr, variablesJSON string
		err := rows.Scan(&idStr, &personaIDStr, &template.Version, &template.Template, &variablesJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		template.ID = uuid.MustParse(idStr)
		template.PersonaID = uuid.MustParse(personaIDStr)

		if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}

		templates = append(templates, &template)
	}

	return templates, nil
}