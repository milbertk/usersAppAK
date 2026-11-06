package baiagent

import (
	"fmt"
	"strings"

	databaseaiconversations "github.com/milbertk/usersAppAK/database/databaseAIConversations"
)

// BGetAIConversations returns the conversations of the logged user.
func BGetAIConversations(userUID string, limit int) ([]databaseaiconversations.AIConversation, error) {
	userUID = strings.TrimSpace(userUID)
	if userUID == "" {
		return nil, fmt.Errorf("logged user is required")
	}

	conversations, err := databaseaiconversations.SelectConversationsByUser(userUID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %v", err)
	}

	return conversations, nil
}
