package llm

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/pkg/errors"
)


// ClaudeEvaluator is an implementation of the Evaluator interface that uses the Anthropic Claude API.
type ClaudeEvaluator struct {
	client       anthropic.Client
	systemPrompt string
}

// NewClaudeEvaluator creates a new ClaudeEvaluator.
func NewClaudeEvaluator(systemPrompt string) (*ClaudeEvaluator, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY environment variable not set")
	}
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &ClaudeEvaluator{
		client:       client,
		systemPrompt: systemPrompt,
	}, nil
}

// Evaluate evaluates the given prompt using the Anthropic Claude API.
func (e *ClaudeEvaluator) Evaluate(ctx context.Context, prompt string) (inputGrade float64, suggestedPrompt string, outputGrade float64, err error) {
	// Create the request
	resp, err := e.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: e.systemPrompt,
			},
		},
	})

	if err != nil {
		return 0, "", 0, errors.Wrap(err, "failed to create messages")
	}

	// Extract text from response
	var responseText string
	if len(resp.Content) > 0 {
		responseText = resp.Content[0].Text
	} else {
		return 0, "", 0, errors.New("empty response from Claude API")
	}


	// Clean the response to extract JSON
	cleanedResponse := strings.TrimSpace(responseText)

	// Look for JSON block if wrapped in markdown
	if strings.Contains(cleanedResponse, "```json") {
		start := strings.Index(cleanedResponse, "```json") + 7
		end := strings.Index(cleanedResponse[start:], "```")
		if end != -1 {
			cleanedResponse = strings.TrimSpace(cleanedResponse[start : start+end])
		}
	} else if strings.Contains(cleanedResponse, "```") {
		// Handle generic code block
		start := strings.Index(cleanedResponse, "```") + 3
		end := strings.Index(cleanedResponse[start:], "```")
		if end != -1 {
			cleanedResponse = strings.TrimSpace(cleanedResponse[start : start+end])
		}
	}

	// Find JSON object boundaries
	startBrace := strings.Index(cleanedResponse, "{")
	lastBrace := strings.LastIndex(cleanedResponse, "}")

	if startBrace != -1 && lastBrace != -1 && lastBrace > startBrace {
		jsonStr := cleanedResponse[startBrace : lastBrace+1]

		// More robust JSON string cleaning
		// First, properly escape newlines within JSON string values
		jsonStr = fixJSONStringValues(jsonStr)

		cleanedResponse = jsonStr
	}

	type evaluationResponse struct {
		InputPromptGrade float64 `json:"inputPromptGrade"`
		SuggestedPrompt  string  `json:"suggestedPrompt"`
		ImprovedGrade    float64 `json:"ImprovedGrade"`
	}

	var evalResp evaluationResponse
	err = json.Unmarshal([]byte(cleanedResponse), &evalResp)
	if err != nil {
		return 0, "", 0, errors.Wrapf(err, "failed to unmarshal evaluation response. Raw response: %s", responseText)
	}

	return evalResp.InputPromptGrade, evalResp.SuggestedPrompt, evalResp.ImprovedGrade, nil
}

// fixJSONStringValues properly handles newlines and quotes within JSON string values
func fixJSONStringValues(jsonStr string) string {
	var result strings.Builder
	inString := false
	escaped := false

	for _, char := range jsonStr {
		if escaped {
			// Previous character was a backslash, so this character is escaped
			result.WriteRune(char)
			escaped = false
			continue
		}

		if char == '\\' {
			// This is an escape character
			result.WriteRune(char)
			escaped = true
			continue
		}

		if char == '"' && !inString {
			// Starting a string
			inString = true
			result.WriteRune(char)
			continue
		}

		if char == '"' && inString {
			// Ending a string
			inString = false
			result.WriteRune(char)
			continue
		}

		if inString {
			// We're inside a JSON string value
			switch char {
			case '\n':
				result.WriteString("\\n")
			case '\r':
				result.WriteString("\\r")
			case '\t':
				result.WriteString("\\t")
			default:
				result.WriteRune(char)
			}
		} else {
			// We're outside string values, normal JSON structure
			result.WriteRune(char)
		}
	}

	return result.String()
}