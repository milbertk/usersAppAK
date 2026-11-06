package databaseAIConversations

import "time"

// AIConversation maps public.ai_conversations
type AIConversation struct {
	ConversationID    string     `json:"conversationId"`
	UserUID           string     `json:"userUid"`
	UserEmail         string     `json:"userEmail"`
	CompanyID         string     `json:"companyId"`
	Platform          string     `json:"platform"`
	Title             string     `json:"title"`
	Model             string     `json:"model"`
	Status            string     `json:"status"`
	TotalInputTokens  int        `json:"totalInputTokens"`
	TotalOutputTokens int        `json:"totalOutputTokens"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

// AIConversationMessage maps public.ai_conversation_messages
type AIConversationMessage struct {
	MessageID      string    `json:"messageId"`
	ConversationID string    `json:"conversationId"`
	UserUID        string    `json:"userUid"`
	Role           string    `json:"role"` // user | assistant
	Content        string    `json:"content"`
	Model          string    `json:"model,omitempty"`
	InputTokens    int       `json:"inputTokens"`
	OutputTokens   int       `json:"outputTokens"`
	StopReason     string    `json:"stopReason,omitempty"`
	LatencyMs      int       `json:"latencyMs"`
	ErrorMessage   string    `json:"errorMessage,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}
