package sqlite

import (
	// "os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
)

func TestSQLiteStorage_Personas(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Initialize storage
	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create schema manually for test
	schema := `
	CREATE TABLE personas (
		id TEXT PRIMARY KEY,
		user_role_display TEXT NOT NULL,
		llm_role_display TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = storage.db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Test Create
	persona := &models.Persona{
		UserRoleDisplay: "Software Developer",
		LLMRoleDisplay:  "Code Reviewer",
	}

	created, err := storage.CreatePersona(persona)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	if created.UserRoleDisplay != persona.UserRoleDisplay {
		t.Errorf("Expected user role display %s, got %s", persona.UserRoleDisplay, created.UserRoleDisplay)
	}

	// Test GetAll
	personas, err := storage.GetAllPersonas()
	if err != nil {
		t.Fatalf("Failed to get all personas: %v", err)
	}

	if len(personas) != 1 {
		t.Errorf("Expected 1 persona, got %d", len(personas))
	}

	// Test GetByID
	found, err := storage.GetPersonaByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get persona by ID: %v", err)
	}

	if found.UserRoleDisplay != created.UserRoleDisplay {
		t.Error("Retrieved persona doesn't match created persona")
	}

	// Test Update
	found.LLMRoleDisplay = "Senior Code Reviewer"
	updated, err := storage.UpdatePersona(found)
	if err != nil {
		t.Fatalf("Failed to update persona: %v", err)
	}

	if updated.LLMRoleDisplay != "Senior Code Reviewer" {
		t.Error("Persona was not updated")
	}

	// Test Delete
	err = storage.DeletePersona(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete persona: %v", err)
	}

	// Verify deletion
	_, err = storage.GetPersonaByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
}

func TestSQLiteStorage_Templates(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create schema
	schema := `
	CREATE TABLE personas (
		id TEXT PRIMARY KEY,
		user_role_display TEXT NOT NULL,
		llm_role_display TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE prompt_templates (
		id TEXT PRIMARY KEY,
		persona_id TEXT NOT NULL,
		template TEXT NOT NULL,
		variables TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (persona_id) REFERENCES personas(id) ON DELETE CASCADE
	);`
	_, err = storage.db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Create a persona first
	persona := &models.Persona{
		UserRoleDisplay: "Developer",
		LLMRoleDisplay:  "Reviewer",
	}
	createdPersona, err := storage.CreatePersona(persona)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Test template creation
	template := &models.PromptTemplate{
		PersonaID: createdPersona.ID,
		Template:  "Review this {{language}} code for {{focus}}",
		Variables: []string{"language", "focus"},
	}

	created, err := storage.CreateTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	// Test GetTemplatesByPersonaID
	templates, err := storage.GetTemplatesByPersonaID(createdPersona.ID)
	if err != nil {
		t.Fatalf("Failed to get templates by persona: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	if len(templates[0].Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(templates[0].Variables))
	}
}

func TestSQLiteStorage_Prompts(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create schema
	schema := `
	CREATE TABLE prompts (
		id TEXT PRIMARY KEY,
		template_id TEXT NOT NULL,
		"values" TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = storage.db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Test prompt creation
	prompt := &models.Prompt{
		TemplateID: uuid.New(),
		Values:     map[string]string{"language": "Go", "focus": "performance"},
		Content:    "Review this Go code for performance",
	}

	created, err := storage.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	// Test GetAll
	prompts, err := storage.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all prompts: %v", err)
	}

	if len(prompts) != 1 {
		t.Errorf("Expected 1 prompt, got %d", len(prompts))
	}

	// Verify values were stored correctly
	if prompts[0].Values["language"] != "Go" {
		t.Error("Values not stored correctly")
	}
}