package databaseAIConversations

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectMessagesByConversation returns the latest `limit` messages of a
// conversation in chronological order (oldest first).
// limit <= 0 returns all messages.
func SelectMessagesByConversation(conversationID string, limit int) ([]AIConversationMessage, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return nil, fmt.Errorf("conversationID is required")
	}

	if limit <= 0 {
		limit = 100000
	}

	query := `
		SELECT * FROM (
			SELECT
				message_id,
				conversation_id,
				user_uid,
				role,
				content,
				COALESCE(model, '') AS model,
				COALESCE(input_tokens, 0) AS input_tokens,
				COALESCE(output_tokens, 0) AS output_tokens,
				COALESCE(stop_reason, '') AS stop_reason,
				COALESCE(latency_ms, 0) AS latency_ms,
				COALESCE(error_message, '') AS error_message,
				created_at
			FROM public.ai_conversation_messages
			WHERE conversation_id::text = $1
			ORDER BY created_at DESC
			LIMIT $2
		) last_messages
		ORDER BY created_at ASC
	`

	rows, err := db.Query(query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get ai messages: %v", err)
	}
	defer rows.Close()

	results := []AIConversationMessage{}

	for rows.Next() {
		var m AIConversationMessage

		err := rows.Scan(
			&m.MessageID,
			&m.ConversationID,
			&m.UserUID,
			&m.Role,
			&m.Content,
			&m.Model,
			&m.InputTokens,
			&m.OutputTokens,
			&m.StopReason,
			&m.LatencyMs,
			&m.ErrorMessage,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ai message row: %v", err)
		}

		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
