package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulguha/promptly/internal/config"
	"github.com/rahulguha/promptly/internal/models"
	"github.com/rahulguha/promptly/internal/storage"
)

func extractVariables(text string) []string {
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

	matches := re.FindAllStringSubmatch(text, -1)

	vars := make(map[string]bool)
	for _, match := range matches {
		vars[match[1]] = true
	}


uniqueVars := make([]string, 0, len(vars))
	for v := range vars {
	
uniqueVars = append(uniqueVars, v)
	}
	return uniqueVars
}

// Handler contains the dependencies for HTTP handlers
type Handler struct {
	DBManager *storage.DBManager
	Cfg       *config.Config
}

// GetPrompts handles GET /prompts
func (h *Handler) GetPrompts(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	prompts, err := store.(storage.Storage).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prompts)
}

// GetPrompt handles GET /prompts/:id
func (h *Handler) GetPrompt(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID format"})
		return
	}

	prompt, err := store.(storage.Storage).GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prompt not found"})
		return
	}

	c.JSON(http.StatusOK, prompt)
}

// CreatePrompt handles POST /prompts
func (h *Handler) CreatePrompt(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	var prompt models.Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if prompt.ProfileID == "" {
		prompt.ProfileID = DefaultProfileID
	}

	createdPrompt, err := store.(storage.Storage).Create(&prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPrompt)
}

// UpdatePrompt handles PUT /prompts/:id
func (h *Handler) UpdatePrompt(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
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
	updatedPrompt, err := store.(storage.Storage).Update(&prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPrompt)
}

// DeletePrompt handles DELETE /prompts/:id
func (h *Handler) DeletePrompt(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID format"})
		return
	}

	err = store.(storage.Storage).Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt deleted successfully"})
}

// Template handlers

// GetTemplates handles GET /templates
func (h *Handler) GetTemplates(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}

	// Log complete request details
	fmt.Println("--- New GetTemplates Request ---")
	fmt.Printf("Request URL: %s %s\n", c.Request.Method, c.Request.URL.String())
	fmt.Println("Request Headers:")
	for key, values := range c.Request.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	fmt.Println("-----------------------------")

	profileID := c.Query("profile_id")
	templates, err := store.(storage.Storage).GetAllTemplates(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log response payload
	fmt.Println("--- GetTemplates Response ---")
	// Marshal the payload to JSON for proper logging
	responsePayload, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling response payload: %v\n", err)
	} else {
		fmt.Printf("Payload:\n%s\n", string(responsePayload))
	}
	fmt.Println("--------------------------")

	c.JSON(http.StatusOK, templates)
}



// GetTemplate handles GET /templates/:id
func (h *Handler) GetTemplate(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	template, err := store.(storage.Storage).GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}
// BuildMetaPrompt constructs a meta prompt given user and LLM roles
func BuildMetaPrompt(userRole, llmRole string) string {
    return fmt.Sprintf(`
I am a %s.
You are a %s. 
Please respond clearly, in a way that fits my background as a %s, 
while staying in your role as a %s.
`, userRole, llmRole, userRole, llmRole)
}

func buildTemplate(metaRole, task, answerGuideline string) string {
	var sb strings.Builder

	if metaRole != "" {
		sb.WriteString("[Meta Role]\n")
		sb.WriteString(metaRole)
		sb.WriteString("\n\n")
	}

	if task != "" {
		sb.WriteString("[Task]\n")
		sb.WriteString(task)
		sb.WriteString("\n\n")
	}

	if answerGuideline != "" {
		sb.WriteString("[Answer Guideline]\n")
		sb.WriteString(answerGuideline)
	}

	return sb.String()
}

// CreateTemplate handles POST /templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if template.ProfileID == "" {
		template.ProfileID = DefaultProfileID
	}

	// Get persona to populate display roles
	persona, err := store.(storage.Storage).GetPersonaByID(template.PersonaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona_id: persona not found"})
		return
	}

	// Prepend persona context with actual values
	metaRole := BuildMetaPrompt(persona.UserRoleDisplay, persona.LLMRoleDisplay)
	template.MetaRole = metaRole

	template.Template = buildTemplate(template.MetaRole, template.Task, template.AnswerGuideline)

	// Extract variables from task and answer guideline
	taskVars := extractVariables(template.Task)
	guidelineVars := extractVariables(template.AnswerGuideline)

	// Combine and get unique variables
	allVars := make(map[string]bool)
	for _, v := range taskVars {
		allVars[v] = true
	}
	for _, v := range guidelineVars {
		allVars[v] = true
	}

	uniqueVars := make([]string, 0, len(allVars))
	for v := range allVars {
	
uniqueVars = append(uniqueVars, v)
	}
	template.Variables = uniqueVars

	createdTemplate, err := store.(storage.Storage).CreateTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTemplate)
}


// UpdateTemplate handles PUT /templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
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

	// Get persona to populate display roles
	persona, err := store.(storage.Storage).GetPersonaByID(template.PersonaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona_id: persona not found"})
		return
	}

	// Prepend persona context with actual values
	metaRole := BuildMetaPrompt(persona.UserRoleDisplay, persona.LLMRoleDisplay)
	template.MetaRole = metaRole

	template.Template = buildTemplate(template.MetaRole, template.Task, template.AnswerGuideline)

	// Extract variables from task and answer guideline
	taskVars := extractVariables(template.Task)
	guidelineVars := extractVariables(template.AnswerGuideline)

	// Combine and get unique variables
	allVars := make(map[string]bool)
	for _, v := range taskVars {
		allVars[v] = true
	}
	for _, v := range guidelineVars {
		allVars[v] = true
	}

	uniqueVars := make([]string, 0, len(allVars))
	for v := range allVars {
	
uniqueVars = append(uniqueVars, v)
	}
	template.Variables = uniqueVars

	updatedTemplate, err := store.(storage.Storage).UpdateTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTemplate)
}

// CreateTemplateVersion handles POST /templates/:id/version
func (h *Handler) CreateTemplateVersion(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
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

	// Get persona to populate display roles
	persona, err := store.(storage.Storage).GetPersonaByID(template.PersonaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona_id: persona not found"})
		return
	}

	// Prepend persona context with actual values
	metaRole := BuildMetaPrompt(persona.UserRoleDisplay, persona.LLMRoleDisplay)
	template.MetaRole = metaRole

	template.Template = buildTemplate(template.MetaRole, template.Task, template.AnswerGuideline)

	// Extract variables from task and answer guideline
	taskVars := extractVariables(template.Task)
	guidelineVars := extractVariables(template.AnswerGuideline)

	// Combine and get unique variables
	allVars := make(map[string]bool)
	for _, v := range taskVars {
		allVars[v] = true
	}
	for _, v := range guidelineVars {
		allVars[v] = true
	}

	uniqueVars := make([]string, 0, len(allVars))
	for v := range allVars {
	
uniqueVars = append(uniqueVars, v)
	}
	template.Variables = uniqueVars

	newVersion, err := store.(storage.Storage).CreateTemplateVersion(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newVersion)
}

// DeleteTemplate handles DELETE /templates/:id?version=N
func (h *Handler) DeleteTemplate(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
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

	err = store.(storage.Storage).DeleteTemplate(id, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GeneratePromptRequest represents the request for generating a prompt
type GeneratePromptRequest struct {
	TemplateID      uuid.UUID         `json:"template_id" binding:"required"`
	TemplateVersion int               `json:"template_version"`
	Name            string            `json:"name"`
	Values          map[string]string `json:"variable_values" binding:"required"`
}

// GeneratePrompt handles POST /generate-prompt
func (h *Handler) GeneratePrompt(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	var req GeneratePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the incoming request
	fmt.Println("--- New GeneratePrompt Request ---")
	reqBytes, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
	} else {
		fmt.Printf("Request Body:\n%s\n", string(reqBytes))
	}

	// Get the template
	template, err := store.(storage.Storage).GetTemplateByID(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Extract variables directly from the template content to ensure they are always up-to-date
	templateVars := extractVariables(template.Template)

	// Create filtered values map - only include variables that exist in the template
	filteredValues := make(map[string]string)
	for _, variable := range templateVars {
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
		Name:            req.Name,
		TemplateID:      req.TemplateID,
		TemplateVersion: template.Version,
		Values:          filteredValues,
		Content:         content,
	}

	createdPrompt, err := store.(storage.Storage).Create(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the response
	fmt.Println("--- GeneratePrompt Response ---")
	resBytes, err := json.MarshalIndent(createdPrompt, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling response: %v\n", err)
	} else {
		fmt.Printf("Response Body:\n%s\n", string(resBytes))
	}
	fmt.Println("-----------------------------")

	c.JSON(http.StatusCreated, createdPrompt)
}



// Persona handlers

// GetPersonas handles GET /personas
func (h *Handler) GetPersonas(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}

	// Log complete request details
	// fmt.Println("--- New GetPersonas Request ---")
	// fmt.Printf("Request URL: %s %s\n", c.Request.Method, c.Request.URL.String())
	// fmt.Println("Request Headers:")
	// for key, values := range c.Request.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("  %s: %s\n", key, value)
	// 	}
	// }
	// fmt.Println("-----------------------------")

	profileID := c.Query("profile_id")
	personas, err := store.(storage.Storage).GetAllPersonas(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log response payload
	fmt.Println("--- GetPersonas Response ---")
	// Marshal the payload to JSON for proper logging
	responsePayload, err := json.MarshalIndent(personas, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling response payload: %v\n", err)
	} else {
		fmt.Printf("Payload:\n%s\n", string(responsePayload))
	}
	fmt.Println("--------------------------")

	c.JSON(http.StatusOK, personas)
}



// GetPersona handles GET /personas/:id
func (h *Handler) GetPersona(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona ID format"})
		return
	}

	persona, err := store.(storage.Storage).GetPersonaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Persona not found"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

const DefaultProfileID = "00000000-0000-0000-0000-000000000000"

// CreatePersona handles POST /personas
func (h *Handler) CreatePersona(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	var persona models.Persona
	if err := c.ShouldBindJSON(&persona); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if persona.ProfileID == "" {
		persona.ProfileID = DefaultProfileID
	}

	createdPersona, err := store.(storage.Storage).CreatePersona(&persona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPersona)
}

// UpdatePersona handles PUT /personas/:id
func (h *Handler) UpdatePersona(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
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
	updatedPersona, err := store.(storage.Storage).UpdatePersona(&persona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPersona)
}

// DeletePersona handles DELETE /personas/:id
func (h *Handler) DeletePersona(c *gin.Context) {
	store, exists := c.Get("store")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage not initialized"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid persona ID format"})
		return
	}

	err = store.(storage.Storage).DeletePersona(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Persona deleted successfully"})
}