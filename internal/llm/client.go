// Package llm wraps the OpenAI API to generate commit messages.
package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Client holds the OpenAI client and the target model name.
type Client struct {
	inner *openai.Client
	model string
}

// NewClient creates a Client with the given API key and model.
func NewClient(apiKey, model string) *Client {
	return &Client{
		inner: openai.NewClient(apiKey),
		model: model,
	}
}

// Generate sends the system and user prompts to the API and returns the
// generated commit message (trimmed of surrounding whitespace).
func (c *Client) Generate(ctx context.Context, system, user string) (string, error) {
	resp, err := c.inner.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: system,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: user,
			},
		},
		MaxTokens:   256,
		Temperature: 0.2, // Low temperature → more deterministic output
	})
	if err != nil {
		return "", fmt.Errorf("openai API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}
