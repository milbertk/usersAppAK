package databaseAIConversations

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectConversationByIDAndUser returns the conversation only if it belongs to the user.
func SelectConversationByIDAndUser(conversationID, userUID string) (*AIConversation, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	conversationID = strings.TrimSpace(conversationID)
	userUID = strings.TrimSpace(userUID)

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
		WHERE conversation_id::text = $1
		  AND user_uid = $2
	`

	var c AIConversation

	err = db.QueryRow(query, conversationID, userUID).Scan(
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

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("conversation does not exist")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get ai conversation: %v", err)
	}

	return &c, nil
}
