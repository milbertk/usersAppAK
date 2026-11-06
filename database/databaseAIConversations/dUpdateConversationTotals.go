package databaseAIConversations

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// AddConversationTokens adds token usage and refreshes updated_at.
func AddConversationTokens(conversationID string, inputTokens, outputTokens int) error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		UPDATE public.ai_conversations
		SET total_input_tokens  = total_input_tokens + $2,
		    total_output_tokens = total_output_tokens + $3,
		    updated_at          = NOW()
		WHERE conversation_id::text = $1
	`

	if _, err := db.Exec(query, conversationID, inputTokens, outputTokens); err != nil {
		return fmt.Errorf("failed to update ai conversation totals: %v", err)
	}

	return nil
}
