package databaseAIConversations

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert stores one message (user or assistant) in ai_conversation_messages.
func (m *AIConversationMessage) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		INSERT INTO public.ai_conversation_messages (
			message_id,
			conversation_id,
			user_uid,
			role,
			content,
			model,
			input_tokens,
			output_tokens,
			stop_reason,
			latency_ms,
			error_message,
			created_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4,
			NULLIF($5, ''),
			$6,
			$7,
			NULLIF($8, ''),
			$9,
			NULLIF($10, ''),
			NOW()
		)
		RETURNING message_id, created_at
	`

	err = db.QueryRow(
		query,
		m.ConversationID,
		m.UserUID,
		m.Role,
		m.Content,
		m.Model,
		m.InputTokens,
		m.OutputTokens,
		m.StopReason,
		m.LatencyMs,
		m.ErrorMessage,
	).Scan(&m.MessageID, &m.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert ai conversation message: %v", err)
	}

	return nil
}
