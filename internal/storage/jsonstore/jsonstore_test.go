package jsonstore

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
)

func TestNewFileStorage(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")

	storage, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatalf("Failed to create FileStorage: %v", err)
	}

	if storage.filePath != filePath {
		t.Errorf("Expected filePath %s, got %s", filePath, storage.filePath)
	}

	// Check if file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
}

func TestFileStorage_Create(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	prompt := &models.Prompt{
		TemplateID: uuid.New(),
		Values:     map[string]string{"name": "world"},
		Content:    "Hello world",
	}

	created, err := storage.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	if created.Content != prompt.Content {
		t.Errorf("Expected content %s, got %s", prompt.Content, created.Content)
	}
}

func TestFileStorage_GetAll(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	// Create test prompts
	prompt1 := &models.Prompt{TemplateID: uuid.New(), Content: "Test 1"}
	prompt2 := &models.Prompt{TemplateID: uuid.New(), Content: "Test 2"}

	storage.Create(prompt1)
	storage.Create(prompt2)

	prompts, err := storage.GetAll("")
	if err != nil {
		t.Fatalf("Failed to get all prompts: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(prompts))
	}
}

func TestFileStorage_GetByID(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	prompt := &models.Prompt{TemplateID: uuid.New(), Content: "Test"}
	created, _ := storage.Create(prompt)

	found, err := storage.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get prompt by ID: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, found.ID)
	}

	// Test non-existent ID
	nonExistentID := uuid.New()
	_, err = storage.GetByID(nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent ID")
	}
}

func TestFileStorage_Update(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	prompt := &models.Prompt{TemplateID: uuid.New(), Content: "Original"}
	created, _ := storage.Create(prompt)

	// Update the prompt
	created.Content = "Updated"
	updated, err := storage.Update(created)
	if err != nil {
		t.Fatalf("Failed to update prompt: %v", err)
	}

	if updated.Content != "Updated" {
		t.Errorf("Expected content 'Updated', got %s", updated.Content)
	}

	// Verify it was persisted
	found, _ := storage.GetByID(created.ID)
	if found.Content != "Updated" {
		t.Error("Update was not persisted")
	}
}

func TestFileStorage_Delete(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	prompt := &models.Prompt{TemplateID: uuid.New(), Content: "Test"}
	created, _ := storage.Create(prompt)

	err := storage.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete prompt: %v", err)
	}

	// Verify it was deleted
	_, err = storage.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted prompt")
	}
}

func TestFileStorage_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent creates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			prompt := &models.Prompt{
				TemplateID: uuid.New(),
				Content:    fmt.Sprintf("Test %d", i),
			}
			_, err := storage.Create(prompt)
			if err != nil {
				t.Errorf("Concurrent create failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all prompts were created
	prompts, err := storage.GetAll("")
	if err != nil {
		t.Fatalf("Failed to get all prompts: %v", err)
	}

	if len(prompts) != numGoroutines {
		t.Errorf("Expected %d prompts, got %d", numGoroutines, len(prompts))
	}
}

func TestFileStorage_PersonaOperations(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	// Test persona CRUD operations
	persona := &models.Persona{
		UserRoleDisplay: "Software Developer",
		LLMRoleDisplay:  "Senior Code Reviewer",
	}

	created, err := storage.CreatePersona(persona)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	// Test get by ID
	found, err := storage.GetPersonaByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get persona: %v", err)
	}

	if found.UserRoleDisplay != persona.UserRoleDisplay {
		t.Errorf("Expected user role display %s, got %s", persona.UserRoleDisplay, found.UserRoleDisplay)
	}
}