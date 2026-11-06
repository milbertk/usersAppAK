package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// ----------------------------------------------------------------------
// Anthropic (Claude) Messages API client — plain net/http, no SDK.
//
// Environment variables (.env):
//   ANTHROPIC_API_KEY     required
//   ANTHROPIC_MODEL       optional (default: defaultAnthropicModel)
//   ANTHROPIC_MAX_TOKENS  optional (default: 1024)
// ----------------------------------------------------------------------

const (
	anthropicURL            = "https://api.anthropic.com/v1/messages"
	anthropicVersion        = "2023-06-01"
	defaultAnthropicModel   = "claude-haiku-5-5"
	defaultAnthropicMaxTok  = 1024
	anthropicRequestTimeout = 60 * time.Second
)

var loadEnvOnce sync.Once

func loadEnv() {
	loadEnvOnce.Do(func() {
		// Ignore the error: variables may already come from the real environment.
		_ = godotenv.Load()
	})
}

// AnthropicMessage is one turn of the conversation sent to Claude.
type AnthropicMessage struct {
	Role    string `json:"role"` // user | assistant
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []AnthropicMessage `json:"messages"`
}

type anthropicResponse struct {
	ID         string `json:"id"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Content    []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// AnthropicResult is what the business layer receives.
type AnthropicResult struct {
	Text         string
	Model        string
	StopReason   string
	InputTokens  int
	OutputTokens int
	LatencyMs    int
}

// AnthropicModel returns the configured model.
func AnthropicModel() string {
	loadEnv()
	if m := strings.TrimSpace(os.Getenv("ANTHROPIC_MODEL")); m != "" {
		return m
	}
	return defaultAnthropicModel
}

func anthropicMaxTokens() int {
	loadEnv()
	if v := strings.TrimSpace(os.Getenv("ANTHROPIC_MAX_TOKENS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultAnthropicMaxTok
}

// CallAnthropic sends the conversation to Claude and returns the text answer.
func CallAnthropic(ctx context.Context, system string, messages []AnthropicMessage) (*AnthropicResult, error) {
	loadEnv()

	apiKey := strings.Trim(strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY")), `"`)
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not configured")
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("messages are required")
	}

	body, err := json.Marshal(anthropicRequest{
		Model:     AnthropicModel(),
		MaxTokens: anthropicMaxTokens(),
		System:    system,
		Messages:  messages,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build anthropic request: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, anthropicRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create anthropic request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	start := time.Now()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic request failed: %v", err)
	}
	defer resp.Body.Close()

	latency := int(time.Since(start).Milliseconds())

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read anthropic response: %v", err)
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("invalid anthropic response (status %d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return nil, fmt.Errorf("anthropic error (%d %s): %s", resp.StatusCode, parsed.Error.Type, parsed.Error.Message)
		}
		return nil, fmt.Errorf("anthropic error (status %d)", resp.StatusCode)
	}

	var sb strings.Builder
	for _, block := range parsed.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}

	return &AnthropicResult{
		Text:         strings.TrimSpace(sb.String()),
		Model:        parsed.Model,
		StopReason:   parsed.StopReason,
		InputTokens:  parsed.Usage.InputTokens,
		OutputTokens: parsed.Usage.OutputTokens,
		LatencyMs:    latency,
	}, nil
}
