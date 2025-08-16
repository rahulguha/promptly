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
		UserRole:  "user",
		LLMRole:   "assistant", 
		Template:  "Hello {{name}}",
		Variables: []string{"name"},
		Values:    map[string]string{"name": "world"},
	}

	created, err := storage.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	if created.Template != prompt.Template {
		t.Errorf("Expected template %s, got %s", prompt.Template, created.Template)
	}
}

func TestFileStorage_GetAll(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	// Create test prompts
	prompt1 := &models.Prompt{UserRole: "user", LLMRole: "assistant", Template: "Test 1"}
	prompt2 := &models.Prompt{UserRole: "user", LLMRole: "assistant", Template: "Test 2"}

	storage.Create(prompt1)
	storage.Create(prompt2)

	prompts, err := storage.GetAll()
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

	prompt := &models.Prompt{UserRole: "user", LLMRole: "assistant", Template: "Test"}
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

	prompt := &models.Prompt{UserRole: "user", LLMRole: "assistant", Template: "Original"}
	created, _ := storage.Create(prompt)

	// Update the prompt
	created.Template = "Updated"
	updated, err := storage.Update(created)
	if err != nil {
		t.Fatalf("Failed to update prompt: %v", err)
	}

	if updated.Template != "Updated" {
		t.Errorf("Expected template 'Updated', got %s", updated.Template)
	}

	// Verify it was persisted
	found, _ := storage.GetByID(created.ID)
	if found.Template != "Updated" {
		t.Error("Update was not persisted")
	}
}

func TestFileStorage_Delete(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	prompt := &models.Prompt{UserRole: "user", LLMRole: "assistant", Template: "Test"}
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
				UserRole: "user",
				LLMRole:  "assistant",
				Template: fmt.Sprintf("Test %d", i),
			}
			_, err := storage.Create(prompt)
			if err != nil {
				t.Errorf("Concurrent create failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all prompts were created
	prompts, err := storage.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all prompts: %v", err)
	}

	if len(prompts) != numGoroutines {
		t.Errorf("Expected %d prompts, got %d", numGoroutines, len(prompts))
	}
}

func TestFileStorage_SearchPrompt(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_prompts.json")
	storage, _ := NewFileStorage(filePath)

	// Create test prompts with different roles
	prompts := []*models.Prompt{
		{UserRole: "developer", LLMRole: "coder", Template: "Code review"},
		{UserRole: "developer", LLMRole: "assistant", Template: "General help"},
		{UserRole: "writer", LLMRole: "editor", Template: "Edit text"},
	}

	for _, p := range prompts {
		storage.Create(p)
	}

	// Test search by user role
	results, err := storage.SearchPrompt("developer", "")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results for developer role, got %d", len(results))
	}

	// Test search by both roles
	results, err = storage.SearchPrompt("developer", "coder")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for developer+coder, got %d", len(results))
	}
}