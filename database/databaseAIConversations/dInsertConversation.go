package databaseAIConversations

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert creates a new conversation in ai_conversations.
func (c *AIConversation) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		INSERT INTO public.ai_conversations (
			conversation_id,
			user_uid,
			user_email,
			company_id,
			platform,
			title,
			model,
			status,
			total_input_tokens,
			total_output_tokens,
			created_at,
			updated_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			NULLIF($3, ''),
			$4,
			$5,
			$6,
			$7,
			0,
			0,
			NOW(),
			NULL
		)
		RETURNING conversation_id, created_at
	`

	err = db.QueryRow(
		query,
		c.UserUID,
		c.UserEmail,
		c.CompanyID,
		c.Platform,
		c.Title,
		c.Model,
		c.Status,
	).Scan(&c.ConversationID, &c.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert ai conversation: %v", err)
	}

	return nil
}
