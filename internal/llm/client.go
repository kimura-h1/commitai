// Package llm wraps the OpenAI-compatible API to generate commit messages.
// It supports both OpenAI and local Ollama (via OpenAI-compatible endpoint).
package llm

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	// ollamaPrefix は --model フラグでOllamaを指定するプレフィックス。
	// 例: --model ollama/llama3
	ollamaPrefix = "ollama/"

	// ollamaBaseURL はOllamaのOpenAI互換エンドポイント。
	ollamaBaseURL = "http://localhost:11434/v1"

	// ollamaDummyKey はOllamaに必要なダミーキー（実際には検証されない）。
	ollamaDummyKey = "ollama"
)

// Client holds an OpenAI-compatible client and the resolved model name.
type Client struct {
	inner    *openai.Client
	model    string
	provider string // "openai" | "ollama"
}

// NewClient creates a Client for the given API key and model.
//
// Ollama モード: model が "ollama/" で始まる場合、ローカルの Ollama に接続する。
//
//	例: model = "ollama/llama3" → Ollama で llama3 を使用
//
// OpenAI モード: それ以外は通常の OpenAI API を使用する。
func NewClient(apiKey, model string) *Client {
	if strings.HasPrefix(model, ollamaPrefix) {
		actualModel := strings.TrimPrefix(model, ollamaPrefix)

		cfg := openai.DefaultConfig(ollamaDummyKey)
		cfg.BaseURL = ollamaBaseURL

		return &Client{
			inner:    openai.NewClientWithConfig(cfg),
			model:    actualModel,
			provider: "ollama",
		}
	}

	return &Client{
		inner:    openai.NewClient(apiKey),
		model:    model,
		provider: "openai",
	}
}

// IsOllama は Ollama モードで動作しているかを返す。
func (c *Client) IsOllama() bool {
	return c.provider == "ollama"
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
		Temperature: 0.2, // 低めに設定して出力を安定させる
	})
	if err != nil {
		// Ollama が起動していない場合、わかりやすいメッセージを返す
		if c.provider == "ollama" && strings.Contains(err.Error(), "connection refused") {
			return "", fmt.Errorf(
				"Ollamaが起動していません。\n" +
					"  以下のコマンドで起動してから再実行してください:\n" +
					"  ollama serve",
			)
		}
		return "", fmt.Errorf("%s API エラー: %w", c.provider, err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%s からの応答が空でした", c.provider)
	}

	return resp.Choices[0].Message.Content, nil
}
