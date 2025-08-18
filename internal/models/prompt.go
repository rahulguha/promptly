
package models

import "github.com/google/uuid"

type Persona struct {
    ID              uuid.UUID `json:"persona_id"`
    UserRoleDisplay string    `json:"user_role_display"`
    LLMRoleDisplay  string    `json:"llm_role_display"`
}

type PromptTemplate struct {
    ID        uuid.UUID `json:"id"`
    PersonaID uuid.UUID `json:"persona_id"`
    Template  string    `json:"template"`
    Variables []string  `json:"variables"`
}

type Prompt struct {
    ID         uuid.UUID         `json:"id"`
    TemplateID uuid.UUID         `json:"template_id"`
    Values     map[string]string `json:"values"`
    Content    string            `json:"content"`
}

type RoleResponse struct {
    UserRoles []string          `json:"user_roles,omitempty"`
    LLMRoles  []string          `json:"llm_roles,omitempty"`
}
