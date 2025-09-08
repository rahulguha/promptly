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

	// Schema is now automatically created by NewSQLiteStorage()

	// Create a profile to satisfy the foreign key constraint
	profile := &models.Profile{
		Name:        "Test Profile",
		Description: "A profile for testing",
	}
	err = storage.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Test Create
	persona := &models.Persona{
		UserRoleDisplay: "Software Developer",
		LLMRoleDisplay:  "Code Reviewer",
		ProfileID:       profile.ID,
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
	personas, err := storage.GetAllPersonas("")
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

	// Schema is now automatically created by NewSQLiteStorage()

	// Create a profile to satisfy the foreign key constraint
	profile := &models.Profile{
		Name:        "Test Profile",
		Description: "A profile for testing",
	}
	err = storage.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Create a persona first
	persona := &models.Persona{
		UserRoleDisplay: "Developer",
		LLMRoleDisplay:  "Reviewer",
		ProfileID:       profile.ID,
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
		ProfileID: profile.ID,
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

	// Schema is now automatically created by NewSQLiteStorage()

	// Create a profile to satisfy the foreign key constraint
	profile := &models.Profile{
		Name:        "Test Profile",
		Description: "A profile for testing",
	}
	err = storage.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Create a persona and template first to satisfy foreign key constraints
	persona := &models.Persona{
		UserRoleDisplay: "Developer",
		LLMRoleDisplay:  "Reviewer",
		ProfileID:       profile.ID,
	}
	createdPersona, err := storage.CreatePersona(persona)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	template := &models.PromptTemplate{
		PersonaID: createdPersona.ID,
		Template:  "Review this {{language}} code for {{focus}}",
		Variables: []string{"language", "focus"},
		ProfileID: profile.ID,
	}
	createdTemplate, err := storage.CreateTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Test prompt creation
	prompt := &models.Prompt{
		TemplateID:      createdTemplate.ID,
		TemplateVersion: createdTemplate.Version,
		Values:          map[string]string{"language": "Go", "focus": "performance"},
		Content:         "Review this Go code for performance",
		ProfileID:       profile.ID,
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