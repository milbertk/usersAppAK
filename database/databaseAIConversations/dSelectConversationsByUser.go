package databaseAIConversations

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectConversationsByUser returns the user's conversations, newest first.
func SelectConversationsByUser(userUID string, limit int) ([]AIConversation, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	userUID = strings.TrimSpace(userUID)
	if userUID == "" {
		return nil, fmt.Errorf("userUID is required")
	}

	if limit <= 0 || limit > 200 {
		limit = 50
	}

	query := `
		SELECT
			conversation_id,
			user_uid,
			COALESCE(user_email, '') AS user_email,
			COALESCE(company_id, '') AS company_id,
			platform,
			COALESCE(title, '') AS title,
			model,
			status,
			total_input_tokens,
			total_output_tokens,
			created_at,
			updated_at
		FROM public.ai_conversations
		WHERE user_uid = $1
		ORDER BY COALESCE(updated_at, created_at) DESC
		LIMIT $2
	`

	rows, err := db.Query(query, userUID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get ai conversations: %v", err)
	}
	defer rows.Close()

	results := []AIConversation{}

	for rows.Next() {
		var c AIConversation

		err := rows.Scan(
			&c.ConversationID,
			&c.UserUID,
			&c.UserEmail,
			&c.CompanyID,
			&c.Platform,
			&c.Title,
			&c.Model,
			&c.Status,
			&c.TotalInputTokens,
			&c.TotalOutputTokens,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ai conversation row: %v", err)
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
