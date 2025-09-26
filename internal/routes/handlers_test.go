package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
	"github.com/rahulguha/promptly/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Set gin to test mode
func init() {
	gin.SetMode(gin.TestMode)
}

// Mock storage for testing
type mockStorage struct {
	mock.Mock
	prompts map[string]*models.Prompt
}

// Prompt operations - matching the Storage interface
func (m *mockStorage) Create(prompt *models.Prompt) (*models.Prompt, error) {
	args := m.Called(prompt)
	if m.prompts == nil {
		m.prompts = make(map[string]*models.Prompt)
	}
	m.prompts[prompt.ID.String()] = prompt
	return args.Get(0).(*models.Prompt), args.Error(1)
}

func (m *mockStorage) GetAll(profileID string) ([]*models.Prompt, error) {
	args := m.Called(profileID)
	return args.Get(0).([]*models.Prompt), args.Error(1)
}

func (m *mockStorage) GetByID(id uuid.UUID) (*models.Prompt, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Prompt), args.Error(1)
}

func (m *mockStorage) Update(prompt *models.Prompt) (*models.Prompt, error) {
	args := m.Called(prompt)
	if m.prompts != nil {
		m.prompts[prompt.ID.String()] = prompt
	}
	return args.Get(0).(*models.Prompt), args.Error(1)
}

func (m *mockStorage) Delete(id uuid.UUID) error {
	args := m.Called(id)
	if m.prompts != nil {
		delete(m.prompts, id.String())
	}
	return args.Error(0)
}

// Persona operations - stub implementations
func (m *mockStorage) GetAllPersonas(profileID string) ([]*models.Persona, error) {
	args := m.Called(profileID)
	return args.Get(0).([]*models.Persona), args.Error(1)
}

func (m *mockStorage) GetPersonaByID(id uuid.UUID) (*models.Persona, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Persona), args.Error(1)
}

func (m *mockStorage) CreatePersona(persona *models.Persona) (*models.Persona, error) {
	args := m.Called(persona)
	return args.Get(0).(*models.Persona), args.Error(1)
}

func (m *mockStorage) UpdatePersona(persona *models.Persona) (*models.Persona, error) {
	args := m.Called(persona)
	return args.Get(0).(*models.Persona), args.Error(1)
}

func (m *mockStorage) DeletePersona(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Template operations - stub implementations
func (m *mockStorage) GetAllTemplates(profileID string) ([]*models.PromptTemplate, error) {
	args := m.Called(profileID)
	return args.Get(0).([]*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) GetTemplateByID(id uuid.UUID) (*models.PromptTemplate, error) {
	args := m.Called(id)
	return args.Get(0).(*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) GetTemplatesByPersonaID(personaID uuid.UUID) ([]*models.PromptTemplate, error) {
	args := m.Called(personaID)
	return args.Get(0).([]*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) CreateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	args := m.Called(template)
	return args.Get(0).(*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) UpdateTemplate(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	args := m.Called(template)
	return args.Get(0).(*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) CreateTemplateVersion(template *models.PromptTemplate) (*models.PromptTemplate, error) {
	args := m.Called(template)
	return args.Get(0).(*models.PromptTemplate), args.Error(1)
}

func (m *mockStorage) DeleteTemplate(id uuid.UUID, version int) error {
	args := m.Called(id, version)
	return args.Error(0)
}

// Cleanup
func (m *mockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Mock evaluator for testing
type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(ctx context.Context, prompt string) (inputGrade float64, suggestedPrompt string, outputGrade float64, err error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).(float64), args.Get(1).(string), args.Get(2).(float64), args.Error(3)
}

// Helper function to create a test handler
func createTestHandler(storage *mockStorage, evaluator *mockEvaluator) *routes.Handler {
	return &routes.Handler{
		LLMEvaluator: evaluator,
	}
}

// Helper function to add mock storage to gin context
func addMockStorageMiddleware(storage *mockStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("store", storage)
		c.Next()
	}
}

func TestCreatePrompt(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.POST("/prompts", handler.CreatePrompt)

	// Test data
	templateID := uuid.New()
	promptID := uuid.New()
	testPrompt := &models.Prompt{
		ID:              promptID,
		Name:            "Test Prompt",
		TemplateID:      templateID,
		TemplateVersion: 1,
		Values:          map[string]string{"key": "value"},
		Content:         "Test content",
		ProfileID:       "test-profile",
	}

	// Mock expectations
	mockStore.On("Create", mock.AnythingOfType("*models.Prompt")).Return(testPrompt, nil)

	// Request body
	body := map[string]interface{}{
		"name":             "Test Prompt",
		"template_id":      templateID.String(),
		"template_version": 1,
		"variable_values":  map[string]string{"key": "value"},
		"content":          "Test content",
		"profile_id":       "test-profile",
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req, _ := http.NewRequest("POST", "/prompts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var responsePrompt models.Prompt
	json.Unmarshal(w.Body.Bytes(), &responsePrompt)
	assert.Equal(t, "Test Prompt", responsePrompt.Name)
	assert.Equal(t, templateID, responsePrompt.TemplateID)
	mockStore.AssertExpectations(t)
}

func TestGetPrompts(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.GET("/prompts", handler.GetPrompts)

	// Test data
	testPrompts := []*models.Prompt{
		{
			ID:              uuid.New(),
			Name:            "Prompt 1",
			TemplateID:      uuid.New(),
			TemplateVersion: 1,
			ProfileID:       "test-profile",
		},
		{
			ID:              uuid.New(),
			Name:            "Prompt 2",
			TemplateID:      uuid.New(),
			TemplateVersion: 1,
			ProfileID:       "test-profile",
		},
	}

	// Mock expectations
	mockStore.On("GetAll", "").Return(testPrompts, nil)

	// Make request
	req, _ := http.NewRequest("GET", "/prompts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var responsePrompts []*models.Prompt
	json.Unmarshal(w.Body.Bytes(), &responsePrompts)
	assert.Len(t, responsePrompts, 2)
	assert.Equal(t, "Prompt 1", responsePrompts[0].Name)
	mockStore.AssertExpectations(t)
}

func TestGetPrompt(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.GET("/prompts/:id", handler.GetPrompt)

	// Test data
	promptID := uuid.New()
	testPrompt := &models.Prompt{
		ID:              promptID,
		Name:            "Test Prompt",
		TemplateID:      uuid.New(),
		TemplateVersion: 1,
		ProfileID:       "test-profile",
	}

	// Mock expectations
	mockStore.On("GetByID", promptID).Return(testPrompt, nil)

	// Make request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/prompts/%s", promptID.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var responsePrompt models.Prompt
	json.Unmarshal(w.Body.Bytes(), &responsePrompt)
	assert.Equal(t, promptID, responsePrompt.ID)
	assert.Equal(t, "Test Prompt", responsePrompt.Name)
	mockStore.AssertExpectations(t)
}

func TestGetPrompt_NotFound(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.GET("/prompts/:id", handler.GetPrompt)

	// Test data
	promptID := uuid.New()

	// Mock expectations - return error for not found
	mockStore.On("GetByID", promptID).Return((*models.Prompt)(nil), fmt.Errorf("prompt not found"))

	// Make request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/prompts/%s", promptID.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockStore.AssertExpectations(t)
}

func TestUpdatePrompt(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.PUT("/prompts/:id", handler.UpdatePrompt)

	// Test data
	promptID := uuid.New()
	templateID := uuid.New()
	testPrompt := &models.Prompt{
		ID:              promptID,
		Name:            "Updated Prompt",
		TemplateID:      templateID,
		TemplateVersion: 2,
		Values:          map[string]string{"key": "updated_value"},
		Content:         "Updated content",
		ProfileID:       "test-profile",
	}

	// Mock expectations
	mockStore.On("Update", mock.AnythingOfType("*models.Prompt")).Return(testPrompt, nil)

	// Request body
	body := map[string]interface{}{
		"name":             "Updated Prompt",
		"template_id":      templateID.String(),
		"template_version": 2,
		"variable_values":  map[string]string{"key": "updated_value"},
		"content":          "Updated content",
		"profile_id":       "test-profile",
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/prompts/%s", promptID.String()), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var responsePrompt models.Prompt
	json.Unmarshal(w.Body.Bytes(), &responsePrompt)
	assert.Equal(t, promptID, responsePrompt.ID)
	assert.Equal(t, "Updated Prompt", responsePrompt.Name)
	mockStore.AssertExpectations(t)
}

func TestDeletePrompt(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.DELETE("/prompts/:id", handler.DeletePrompt)

	// Test data
	promptID := uuid.New()

	// Mock expectations
	mockStore.On("Delete", promptID).Return(nil)

	// Make request
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/prompts/%s", promptID.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Prompt deleted successfully", response["message"])
	mockStore.AssertExpectations(t)
}

func TestDeletePrompt_NotFound(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.DELETE("/prompts/:id", handler.DeletePrompt)

	// Test data
	promptID := uuid.New()

	// Mock expectations - return error for not found
	mockStore.On("Delete", promptID).Return(fmt.Errorf("prompt not found"))

	// Make request
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/prompts/%s", promptID.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockStore.AssertExpectations(t)
}

func TestEvaluatePrompt(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.POST("/evaluate", handler.EvaluatePrompt)

	// Mock expectations - updated for new signature
	mockEval.On("Evaluate", mock.Anything, "test prompt").
		Return(7.5, "This is an improved prompt!", 9.0, nil)

	// Request body
	body := map[string]string{
		"prompt": "test prompt",
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req, _ := http.NewRequest("POST", "/evaluate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	assert.Equal(t, 7.5, responseBody["inputGrade"].(float64))
	assert.Equal(t, "This is an improved prompt!", responseBody["suggestedPrompt"].(string))
	assert.Equal(t, 9.0, responseBody["outputGrade"].(float64))
	mockEval.AssertExpectations(t)
}

func TestEvaluatePrompt_InvalidPayload(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.POST("/evaluate", handler.EvaluatePrompt)

	// Invalid request body (missing prompt field)
	body := map[string]string{
		"invalid": "field",
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req, _ := http.NewRequest("POST", "/evaluate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "Invalid request payload", responseBody["error"])
}

func TestCreatePrompt_InvalidPayload(t *testing.T) {
	// Setup
	mockStore := &mockStorage{}
	mockEval := &mockEvaluator{}
	handler := createTestHandler(mockStore, mockEval)

	r := gin.New()
	r.Use(addMockStorageMiddleware(mockStore))
	r.POST("/prompts", handler.CreatePrompt)

	// Invalid request body (invalid JSON)
	invalidJSON := []byte(`{"name": "test", "invalid": }`)

	// Make request
	req, _ := http.NewRequest("POST", "/prompts", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Contains(t, responseBody["error"], "invalid character")
}