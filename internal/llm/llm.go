package llm

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)


// Evaluator is an interface for evaluating prompts.
type Evaluator interface {
	Evaluate(ctx context.Context, prompt string) (inputGrade float64, suggestedPrompt string, outputGrade float64, err error)
}

// OpenAIEvaluator is an implementation of the Evaluator interface that uses the OpenAI API.
type OpenAIEvaluator struct {
	client       *openai.Client
	systemPrompt string
}

// NewOpenAIEvaluator creates a new OpenAIEvaluator.
func NewOpenAIEvaluator(systemPrompt string) (*OpenAIEvaluator, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable not set")
	}
	client := openai.NewClient(apiKey)
	return &OpenAIEvaluator{
		client:       client,
		systemPrompt: systemPrompt,
	}, nil
}

// Evaluate evaluates the given prompt using the OpenAI API.
func (e *OpenAIEvaluator) Evaluate(ctx context.Context, prompt string) (inputGrade float64, suggestedPrompt string, outputGrade float64, err error) {
	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: e.systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return 0, "", 0, errors.Wrap(err, "failed to create chat completion")
	}


	type evaluationResponse struct {
		InputPromptGrade float64 `json:"inputPromptGrade"`
		SuggestedPrompt  string  `json:"suggestedPrompt"`
		ImprovedGrade    float64 `json:"ImprovedGrade"`
	}

	var evalResp evaluationResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &evalResp)
	if err != nil {
		return 0, "", 0, errors.Wrap(err, "failed to unmarshal evaluation response")
	}

	return evalResp.InputPromptGrade, evalResp.SuggestedPrompt, evalResp.ImprovedGrade, nil
}
