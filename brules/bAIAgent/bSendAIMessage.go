package baiagent

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	databaseaiconversations "github.com/milbertk/usersAppAK/database/databaseAIConversations"
	"github.com/milbertk/usersAppAK/helpers"
)

const (
	maxMessageChars   = 4000 // max characters per user message
	maxHistoryMessage = 20   // last N messages sent to Claude as context
	defaultPlatform   = "UserManagement"
)

type SendAIMessageRequest struct {
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
	Platform       string `json:"platform"`
	CompanyID      string `json:"companyId"`
	CompanyName    string `json:"companyName"`
	CurrentPath    string `json:"currentPath"`
	Language       string `json:"language"`
	LoggedUID      string `json:"-"`
	LoggedEmail    string `json:"-"`
}

type SendAIMessageResponse struct {
	ConversationID  string `json:"conversationId"`
	NewConversation bool   `json:"newConversation"`
	Reply           string `json:"reply"`
	Model           string `json:"model"`
	StopReason      string `json:"stopReason"`
	InputTokens     int    `json:"inputTokens"`
	OutputTokens    int    `json:"outputTokens"`
	LatencyMs       int    `json:"latencyMs"`
}

// BSendAIMessage saves the user message, calls Claude with the history,
// saves the answer and returns it.
func BSendAIMessage(ctx context.Context, req SendAIMessageRequest) (*SendAIMessageResponse, error) {
	req.ConversationID = strings.TrimSpace(req.ConversationID)
	req.Message = strings.TrimSpace(req.Message)
	req.Platform = strings.TrimSpace(req.Platform)
	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.LoggedUID = strings.TrimSpace(req.LoggedUID)
	req.CompanyName = limitText(req.CompanyName, 150)
	req.CurrentPath = limitText(req.CurrentPath, 200)
	req.Language = limitText(req.Language, 10)

	if req.LoggedUID == "" {
		return nil, fmt.Errorf("logged user is required")
	}

	if req.Message == "" {
		return nil, fmt.Errorf("message is required")
	}

	if utf8.RuneCountInString(req.Message) > maxMessageChars {
		return nil, fmt.Errorf("message is too long (max %d characters)", maxMessageChars)
	}

	if req.Platform == "" {
		req.Platform = defaultPlatform
	}

	newConversation := false
	var conversation *databaseaiconversations.AIConversation

	// 1. Load (and validate ownership) or create the conversation
	if req.ConversationID != "" {
		c, err := databaseaiconversations.SelectConversationByIDAndUser(req.ConversationID, req.LoggedUID)
		if err != nil {
			return nil, err
		}
		conversation = c
	} else {
		conversation = &databaseaiconversations.AIConversation{
			UserUID:   req.LoggedUID,
			UserEmail: req.LoggedEmail,
			CompanyID: req.CompanyID,
			Platform:  req.Platform,
			Title:     buildTitle(req.Message),
			Model:     helpers.AnthropicModel(),
			Status:    "ACTIVE",
		}
		if err := conversation.Insert(); err != nil {
			return nil, fmt.Errorf("error creating conversation: %v", err)
		}
		newConversation = true
	}

	// 2. Save the user message
	userMsg := databaseaiconversations.AIConversationMessage{
		ConversationID: conversation.ConversationID,
		UserUID:        req.LoggedUID,
		Role:           "user",
		Content:        req.Message,
	}
	if err := userMsg.Insert(); err != nil {
		return nil, fmt.Errorf("error saving user message: %v", err)
	}

	// 3. Load history (includes the message just saved)
	history, err := databaseaiconversations.SelectMessagesByConversation(conversation.ConversationID, maxHistoryMessage)
	if err != nil {
		return nil, fmt.Errorf("error loading history: %v", err)
	}

	messages := buildAnthropicMessages(history)

	// 4. Call Claude
	result, err := helpers.CallAnthropic(ctx, buildSystemPrompt(AIRequestContext{
		Platform:    conversation.Platform,
		CompanyName: req.CompanyName,
		CurrentPath: req.CurrentPath,
		Language:    req.Language,
	}), messages)
	if err != nil {
		// Keep a record of the failure in the conversation
		failMsg := databaseaiconversations.AIConversationMessage{
			ConversationID: conversation.ConversationID,
			UserUID:        req.LoggedUID,
			Role:           "assistant",
			Content:        "",
			Model:          helpers.AnthropicModel(),
			ErrorMessage:   err.Error(),
		}
		_ = failMsg.Insert()
		return nil, fmt.Errorf("error calling AI agent: %v", err)
	}

	// 5. Save the assistant answer
	assistantMsg := databaseaiconversations.AIConversationMessage{
		ConversationID: conversation.ConversationID,
		UserUID:        req.LoggedUID,
		Role:           "assistant",
		Content:        result.Text,
		Model:          result.Model,
		InputTokens:    result.InputTokens,
		OutputTokens:   result.OutputTokens,
		StopReason:     result.StopReason,
		LatencyMs:      result.LatencyMs,
	}
	if err := assistantMsg.Insert(); err != nil {
		return nil, fmt.Errorf("error saving assistant message: %v", err)
	}

	// 6. Update totals (non-blocking for the user)
	if err := databaseaiconversations.AddConversationTokens(
		conversation.ConversationID,
		result.InputTokens,
		result.OutputTokens,
	); err != nil {
		fmt.Println("Error updating ai conversation totals:", err)
	}

	return &SendAIMessageResponse{
		ConversationID:  conversation.ConversationID,
		NewConversation: newConversation,
		Reply:           result.Text,
		Model:           result.Model,
		StopReason:      result.StopReason,
		InputTokens:     result.InputTokens,
		OutputTokens:    result.OutputTokens,
		LatencyMs:       result.LatencyMs,
	}, nil
}

// buildAnthropicMessages converts DB rows into Claude messages.
// - skips failed/empty assistant rows
// - merges consecutive messages with the same role (Claude requires alternation)
// - makes sure the first message is from the user
func buildAnthropicMessages(history []databaseaiconversations.AIConversationMessage) []helpers.AnthropicMessage {
	out := []helpers.AnthropicMessage{}

	for _, m := range history {
		content := strings.TrimSpace(m.Content)
		if content == "" || (m.Role != "user" && m.Role != "assistant") {
			continue
		}

		if n := len(out); n > 0 && out[n-1].Role == m.Role {
			out[n-1].Content += "\n\n" + content
			continue
		}

		out = append(out, helpers.AnthropicMessage{Role: m.Role, Content: content})
	}

	for len(out) > 0 && out[0].Role != "user" {
		out = out[1:]
	}

	return out
}

func buildTitle(message string) string {
	title := strings.Join(strings.Fields(message), " ")
	if utf8.RuneCountInString(title) > 80 {
		r := []rune(title)
		title = string(r[:80]) + "..."
	}
	return title
}
