package baiagent

import (
	"fmt"
	"strings"

	databaseaiconversations "github.com/milbertk/usersAppAK/database/databaseAIConversations"
)

// BGetAIConversationMessages returns all messages of a conversation,
// only if the conversation belongs to the logged user.
func BGetAIConversationMessages(conversationID, userUID string) ([]databaseaiconversations.AIConversationMessage, error) {
	conversationID = strings.TrimSpace(conversationID)
	userUID = strings.TrimSpace(userUID)

	if userUID == "" {
		return nil, fmt.Errorf("logged user is required")
	}

	if conversationID == "" {
		return nil, fmt.Errorf("conversationId is required")
	}

	// Ownership validation
	if _, err := databaseaiconversations.SelectConversationByIDAndUser(conversationID, userUID); err != nil {
		return nil, err
	}

	messages, err := databaseaiconversations.SelectMessagesByConversation(conversationID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	return messages, nil
}
