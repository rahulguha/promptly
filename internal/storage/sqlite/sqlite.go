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

const Schema = `
-- Profiles table - stores user-defined personas
CREATE TABLE IF NOT EXISTS profiles (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	attributes TEXT, -- JSON blob for structured attributes
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Personas table - stores user and LLM role definitions
CREATE TABLE IF NOT EXISTS personas (
	id TEXT PRIMARY KEY,
	user_role_display TEXT NOT NULL,
	llm_role_display TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	profile_id TEXT,
	FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

-- Prompt templates table - stores reusable prompt templates with variables
CREATE TABLE IF NOT EXISTS prompt_templates (
	id TEXT PRIMARY KEY,
	name TEXT,
	persona_id TEXT NOT NULL,
	version INTEGER NOT NULL DEFAULT 1,
	meta_role TEXT,
	task TEXT,
	answer_guideline TEXT,
	template TEXT NOT NULL,
	variables TEXT NOT NULL, -- JSON array of variable names
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	profile_id TEXT,
	FOREIGN KEY (persona_id) REFERENCES personas(id) ON DELETE CASCADE,
	FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

-- Prompts table - stores generated prompts from templates
CREATE TABLE IF NOT EXISTS prompts (
	id TEXT PRIMARY KEY,
	name TEXT,
	template_id TEXT NOT NULL,
	template_version INTEGER NOT NULL DEFAULT 1,
	variable_values TEXT NOT NULL, -- JSON object with variable values
	content TEXT NOT NULL, -- Final generated prompt content
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	profile_id TEXT,
	FOREIGN KEY (template_id) REFERENCES prompt_templates(id) ON DELETE CASCADE,
	FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_personas_user_role ON personas(user_role_display);
CREATE INDEX IF NOT EXISTS idx_personas_llm_role ON personas(llm_role_display);
CREATE INDEX IF NOT EXISTS idx_templates_persona ON prompt_templates(persona_id);
CREATE INDEX IF NOT EXISTS idx_prompts_template ON prompts(template_id);
`

type SQLiteStorage struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// NewSQLiteStorage creates a new SQLite storage instance by opening a new DB connection
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

	// Initialize schema if it doesn't exist
	if err := InitializeSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SQLiteStorage{
		db:   db,
		path: dbPath,
	}, nil
}

// NewSQLiteStorageWithDB creates a new SQLite storage instance from an existing DB connection
func NewSQLiteStorageWithDB(db *sql.DB) *SQLiteStorage {
	return &SQLiteStorage{
		db: db,
	}
}

// Close closes the database connection.
// Note: This should only be called if the storage instance was created with NewSQLiteStorage.
func (s *SQLiteStorage) Close() error {
	if s.path != "" { // Only close DBs opened by NewSQLiteStorage
		return s.db.Close()
	}
	return nil
}

// InitializeSchema creates the database schema on a given DB connection
func InitializeSchema(db *sql.DB) error {
	_, err := db.Exec(Schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

// Persona operations
func (s *SQLiteStorage) CreatePersona(persona *models.Persona) (*models.Persona, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	persona.ID = uuid.New()

	query := `INSERT INTO personas (id, user_role_display, llm_role_display, profile_id) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, persona.ID.String(), persona.UserRoleDisplay, persona.LLMRoleDisplay, persona.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to create persona: %w", err)
	}

	return persona, nil
}

func (s *SQLiteStorage) GetAllPersonas(profileID string) ([]*models.Persona, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, user_role_display, llm_role_display, profile_id FROM personas`
	args := []interface{}{}

	if profileID != "" {
		query += " WHERE profile_id = ? OR profile_id = '00000000-0000-0000-0000-000000000000'"
		args = append(args, profileID)
	}

	query += " ORDER BY created_at"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query personas: %w", err)
	}
	defer rows.Close()

	var personas []*models.Persona
	for rows.Next() {
		var persona models.Persona
		var idStr string
		var dbProfileID sql.NullString // Use sql.NullString for nullable profile_id
		err := rows.Scan(&idStr, &persona.UserRoleDisplay, &persona.LLMRoleDisplay, &dbProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan persona: %w", err)
		}
		persona.ID = uuid.MustParse(idStr)
		if dbProfileID.Valid {
			persona.ProfileID = dbProfileID.String
		}
		personas = append(personas, &persona)
	}

	return personas, nil
}

func (s *SQLiteStorage) GetPersonaByID(id uuid.UUID) (*models.Persona, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var persona models.Persona
	var idStr string
	var profileID sql.NullString
	query := `SELECT id, user_role_display, llm_role_display, profile_id FROM personas WHERE id = ?`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &persona.UserRoleDisplay, &persona.LLMRoleDisplay, &profileID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("persona not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get persona: %w", err)
	}

	persona.ID = uuid.MustParse(idStr)
	if profileID.Valid {
		persona.ProfileID = profileID.String
	}
	return &persona, nil
}

func (s *SQLiteStorage) UpdatePersona(persona *models.Persona) (*models.Persona, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE personas SET user_role_display = ?, llm_role_display = ?, profile_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := s.db.Exec(query, persona.UserRoleDisplay, persona.LLMRoleDisplay, persona.ProfileID, persona.ID.String())
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

	query := `INSERT INTO prompt_templates (id, name, persona_id, version, meta_role, task, answer_guideline, template, variables, profile_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, template.ID.String(), template.Name, template.PersonaID.String(), template.Version, template.MetaRole, template.Task, template.AnswerGuideline, template.Template, string(variablesJSON), template.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

func (s *SQLiteStorage) GetAllTemplates(profileID string) ([]*models.PromptTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT pt.id, pt.name, pt.persona_id, pt.version, pt.meta_role, pt.task, pt.answer_guideline, pt.template, pt.variables, pt.profile_id 
			  FROM prompt_templates pt`
	args := []interface{}{}

	if profileID != "" {
		query += " LEFT JOIN personas p ON pt.persona_id = p.id WHERE pt.profile_id = ? OR p.profile_id = ?"
		args = append(args, profileID, profileID)
	}

	query += " ORDER BY pt.created_at"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.PromptTemplate
	for rows.Next() {
		var template models.PromptTemplate
		var idStr, personaIDStr, variablesJSON string
		var dbProfileID sql.NullString
		err := rows.Scan(&idStr, &template.Name, &personaIDStr, &template.Version, &template.MetaRole, &template.Task, &template.AnswerGuideline, &template.Template, &variablesJSON, &dbProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		template.ID = uuid.MustParse(idStr)
		template.PersonaID = uuid.MustParse(personaIDStr)
		if dbProfileID.Valid {
			template.ProfileID = dbProfileID.String
		}

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
	var profileID sql.NullString
	query := `SELECT id, name, persona_id, version, meta_role, task, answer_guideline, template, variables, profile_id FROM prompt_templates WHERE id = ? ORDER BY version DESC LIMIT 1`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &template.Name, &personaIDStr, &template.Version, &template.MetaRole, &template.Task, &template.AnswerGuideline, &template.Template, &variablesJSON, &profileID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	template.ID = uuid.MustParse(idStr)
	template.PersonaID = uuid.MustParse(personaIDStr)
	if profileID.Valid {
		template.ProfileID = profileID.String
	}

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

	query := `UPDATE prompt_templates SET name = ?, persona_id = ?, meta_role = ?, task = ?, answer_guideline = ?, template = ?, variables = ?, profile_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND version = ?`
	result, err := s.db.Exec(query, template.Name, template.PersonaID.String(), template.MetaRole, template.Task, template.AnswerGuideline, template.Template, string(variablesJSON), template.ProfileID, template.ID.String(), template.Version)
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

	query := `INSERT INTO prompt_templates (id, name, persona_id, version, meta_role, task, answer_guideline, template, variables) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, template.ID.String(), template.Name, template.PersonaID.String(), template.Version, template.MetaRole, template.Task, template.AnswerGuideline, template.Template, string(variablesJSON))
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

	query := `INSERT INTO prompts (id, name, template_id, template_version, variable_values, content, profile_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	// Log the SQL statement
	fmt.Println("--- SQL Statement ---")
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Args: %v\n", []interface{}{prompt.ID.String(), prompt.Name, prompt.TemplateID.String(), prompt.TemplateVersion, string(valuesJSON), prompt.Content, prompt.ProfileID})
	fmt.Println("---------------------")

	_, err = s.db.Exec(query, prompt.ID.String(), prompt.Name, prompt.TemplateID.String(), prompt.TemplateVersion, string(valuesJSON), prompt.Content, prompt.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	return prompt, nil
}


func (s *SQLiteStorage) GetAll(profileID string) ([]*models.Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, name, template_id, template_version, variable_values, content, profile_id FROM prompts`
	args := []interface{}{}

	if profileID != "" {
		query += " WHERE profile_id = ?"
		args = append(args, profileID)
	}

	query += " ORDER BY created_at"
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query prompts: %w", err)
	}
	defer rows.Close()

	var prompts []*models.Prompt
	for rows.Next() {
		var prompt models.Prompt
		var idStr, templateIDStr, valuesJSON string
		var dbProfileID sql.NullString
		err := rows.Scan(&idStr, &prompt.Name, &templateIDStr, &prompt.TemplateVersion, &valuesJSON, &prompt.Content, &dbProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prompt: %w", err)
		}

		prompt.ID = uuid.MustParse(idStr)
		prompt.TemplateID = uuid.MustParse(templateIDStr)
		if dbProfileID.Valid {
			prompt.ProfileID = dbProfileID.String
		}

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
	var profileID sql.NullString
	query := `SELECT id, name, template_id, template_version, variable_values, content, profile_id FROM prompts WHERE id = ?`
	err := s.db.QueryRow(query, id.String()).Scan(&idStr, &prompt.Name, &templateIDStr, &prompt.TemplateVersion, &valuesJSON, &prompt.Content, &profileID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("prompt not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	prompt.ID = uuid.MustParse(idStr)
	prompt.TemplateID = uuid.MustParse(templateIDStr)
	if profileID.Valid {
		prompt.ProfileID = profileID.String
	}

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

	query := `UPDATE prompts SET name = ?, template_id = ?, template_version = ?, variable_values = ?, content = ?, profile_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := s.db.Exec(query, prompt.Name, prompt.TemplateID.String(), prompt.TemplateVersion, string(valuesJSON), prompt.Content, prompt.ProfileID, prompt.ID.String())
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

	query := `SELECT id, name, persona_id, version, meta_role, task, answer_guideline, template, variables FROM prompt_templates WHERE persona_id = ? ORDER BY created_at`
	rows, err := s.db.Query(query, personaID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query templates by persona: %w", err)
	}
	defer rows.Close()

	var templates []*models.PromptTemplate
	for rows.Next() {
		var template models.PromptTemplate
		var idStr, personaIDStr, variablesJSON string
		err := rows.Scan(&idStr, &template.Name, &personaIDStr, &template.Version, &template.MetaRole, &template.Task, &template.AnswerGuideline, &template.Template, &variablesJSON)
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

// Profile operations

func (s *SQLiteStorage) CreateProfile(profile *models.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	profile.ID = uuid.New().String()

	attributesJSON, err := json.Marshal(profile.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	query := `INSERT INTO profiles (id, name, description, attributes) VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(query, profile.ID, profile.Name, profile.Description, string(attributesJSON))
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

func (s *SQLiteStorage) GetAllProfiles() ([]*models.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, name, description, attributes, created_at, updated_at FROM profiles ORDER BY created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*models.Profile
	for rows.Next() {
		var profile models.Profile
		var attributesJSON string
		err := rows.Scan(&profile.ID, &profile.Name, &profile.Description, &attributesJSON, &profile.CreatedAt, &profile.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		if err := json.Unmarshal([]byte(attributesJSON), &profile.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func (s *SQLiteStorage) GetProfileByID(id string) (*models.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var profile models.Profile
	var attributesJSON string
	query := `SELECT id, name, description, attributes, created_at, updated_at FROM profiles WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&profile.ID, &profile.Name, &profile.Description, &attributesJSON, &profile.CreatedAt, &profile.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if err := json.Unmarshal([]byte(attributesJSON), &profile.Attributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	return &profile, nil
}

func (s *SQLiteStorage) UpdateProfile(profile *models.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	attributesJSON, err := json.Marshal(profile.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	query := `UPDATE profiles SET name = ?, description = ?, attributes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := s.db.Exec(query, profile.Name, profile.Description, string(attributesJSON), profile.ID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}

func (s *SQLiteStorage) DeleteProfile(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM profiles WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}
