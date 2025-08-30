
package models

import "github.com/google/uuid"

type Persona struct {
    ID              uuid.UUID `json:"persona_id"`
    UserRoleDisplay string    `json:"user_role_display"`
    LLMRoleDisplay  string    `json:"llm_role_display"`
}

type PromptTemplate struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	PersonaID       uuid.UUID `json:"persona_id"`
	Version         int       `json:"version"`
	MetaRole        string    `json:"meta_role"`
	Task            string    `json:"task"`
	AnswerGuideline string    `json:"answer_guideline"`
	Template        string    `json:"template"`
	Variables       []string  `json:"variables"`
}

type Prompt struct {
    ID              uuid.UUID         `json:"id"`
    Name            string            `json:"name"`
    TemplateID      uuid.UUID         `json:"template_id"`
    TemplateVersion int               `json:"template_version"`
    Values          map[string]string `json:"variable_values"`
    Content         string            `json:"content"`
}

type RoleResponse struct {
    UserRoles []string          `json:"user_roles,omitempty"`
    LLMRoles  []string          `json:"llm_roles,omitempty"`
}
