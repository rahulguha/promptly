
package models

import "github.com/google/uuid"

type Prompt struct {
    ID         uuid.UUID         `json:"id"`
    UserRole   string            `json:"user_role"`
    LLMRole    string            `json:"llm_role"`
    Template   string            `json:"template"`
    Variables  []string          `json:"variables"`
    Values     map[string]string `json:"values"`
}

type RoleResponse struct {
    UserRoles []string          `json:"user_roles,omitempty"`
    LLMRoles  []string          `json:"llm_roles,omitempty"`
}
