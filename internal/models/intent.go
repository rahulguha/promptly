package models

// Intent represents a single intent from the intent master file.
type Intent struct {
	Intent       string   `json:"intent"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	SystemPrompt string   `json:"system_prompt"`
	Keywords     []string `json:"keywords"`
	Tag          string   `json:"tag"`
}
