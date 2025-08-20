package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/models"
	"github.com/rahulguha/promptly/internal/storage"
)

// Handler contains the dependencies for HTTP handlers
type Handler struct {
	Store storage.Storage
}

// GetPrompts handles GET /prompts
func (h *Handler) GetPrompts(c *gin.Context) {
	prompts, err := h.Store.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prompts)
}

// GetPrompt handles GET /prompts/:id
func (h *Handler) GetPrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID format"})
		return
	}

	prompt, err := h.Store.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prompt not found"})
		return
	}

	c.JSON(http.StatusOK, prompt)
}

// CreatePrompt handles POST /prompts
func (h *Handler) CreatePrompt(c *gin.Context) {
	var prompt models.Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPrompt, err := h.Store.Create(&prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPrompt)
}

// UpdatePrompt handles PUT /prompts/:id
func (h *Handler) UpdatePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID format"})
		return
	}

	var prompt models.Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt.ID = id
	updatedPrompt, err := h.Store.Update(&prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPrompt)
}

// DeletePrompt handles DELETE /prompts/:id
func (h *Handler) DeletePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID format"})
		return
	}

	err = h.Store.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt deleted successfully"})
}

// Template handlers

// GetTemplates handles GET /templates
func (h *Handler) GetTemplates(c *gin.Context) {
	templates, err := h.Store.GetAllTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetTemplate handles GET /templates/:id
func (h *Handler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	template, err := h.Store.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateTemplate handles POST /templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get persona to populate display roles
	persona, err := h.Store.GetPersonaByID(template.PersonaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona_id: persona not found"})
		return
	}

	// Prepend persona context with actual values
	personaContext := fmt.Sprintf("User is a %s and wants LLM to play the role of %s. ", 
		persona.UserRoleDisplay, persona.LLMRoleDisplay)
	template.Template = personaContext + template.Template

	createdTemplate, err := h.Store.CreateTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTemplate)
}

// UpdateTemplate handles PUT /templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	updatedTemplate, err := h.Store.UpdateTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTemplate)
}

// CreateTemplateVersion handles POST /templates/:id/version
func (h *Handler) CreateTemplateVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	newVersion, err := h.Store.CreateTemplateVersion(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newVersion)
}

// DeleteTemplate handles DELETE /templates/:id?version=N
func (h *Handler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	versionStr := c.Query("version")
	if versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version parameter is required"})
		return
	}

	version := 0
	if _, err := fmt.Sscanf(versionStr, "%d", &version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version format"})
		return
	}

	err = h.Store.DeleteTemplate(id, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GeneratePromptRequest represents the request for generating a prompt
type GeneratePromptRequest struct {
	TemplateID uuid.UUID         `json:"template_id" binding:"required"`
	Values     map[string]string `json:"variable_values" binding:"required"`
}

// GeneratePrompt handles POST /generate-prompt
func (h *Handler) GeneratePrompt(c *gin.Context) {
	var req GeneratePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the template
	template, err := h.Store.GetTemplateByID(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Create filtered values map - only include variables that exist in template
	filteredValues := make(map[string]string)
	for _, variable := range template.Variables {
		if value, exists := req.Values[variable]; exists {
			filteredValues[variable] = value
		}
	}

	// Generate the final prompt content by substituting variables
	content := template.Template
	for variable, value := range filteredValues {
		placeholder := "{{" + variable + "}}"
		content = strings.ReplaceAll(content, placeholder, value)
	}

	// Create new prompt with rendered content
	prompt := &models.Prompt{
		TemplateID:      req.TemplateID,
		TemplateVersion: template.Version,
		Values:          filteredValues,
		Content:         content,
	}

	createdPrompt, err := h.Store.Create(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPrompt)
}

// Persona handlers

// GetPersonas handles GET /personas
func (h *Handler) GetPersonas(c *gin.Context) {
	personas, err := h.Store.GetAllPersonas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, personas)
}

// GetPersona handles GET /personas/:id
func (h *Handler) GetPersona(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona ID format"})
		return
	}

	persona, err := h.Store.GetPersonaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Persona not found"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// CreatePersona handles POST /personas
func (h *Handler) CreatePersona(c *gin.Context) {
	var persona models.Persona
	if err := c.ShouldBindJSON(&persona); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPersona, err := h.Store.CreatePersona(&persona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPersona)
}

// UpdatePersona handles PUT /personas/:id
func (h *Handler) UpdatePersona(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona ID format"})
		return
	}

	var persona models.Persona
	if err := c.ShouldBindJSON(&persona); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persona.ID = id
	updatedPersona, err := h.Store.UpdatePersona(&persona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPersona)
}

// DeletePersona handles DELETE /personas/:id
func (h *Handler) DeletePersona(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona ID format"})
		return
	}

	err = h.Store.DeletePersona(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Persona deleted successfully"})
}