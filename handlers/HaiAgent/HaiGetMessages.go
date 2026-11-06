package haiagent

import (
	"net/http"
	"strings"

	"github.com/milbertk/security"

	baiagent "github.com/milbertk/usersAppAK/brules/bAIAgent"
)

// HaiGetMessages returns the messages of one conversation of the logged user.
// Query: ?conversationId=<uuid>
func HaiGetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendAIAgentResponse(w, http.StatusMethodNotAllowed, "Use GET", nil)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendAIAgentResponse(w, http.StatusUnauthorized, "Missing user claims", nil)
		return
	}

	loggedUID := strings.TrimSpace(claims.UID)
	conversationID := strings.TrimSpace(r.URL.Query().Get("conversationId"))

	messages, err := baiagent.BGetAIConversationMessages(conversationID, loggedUID)
	if err != nil {
		insertAIAgentTracking(r, loggedUID, claims.Email, "AIGetConversationMessages", map[string]interface{}{
			"conversationId": conversationID,
			"success":        false,
			"message":        err.Error(),
		})
		sendAIAgentResponse(w, statusFromAIError(err.Error()), err.Error(), nil)
		return
	}

	insertAIAgentTracking(r, loggedUID, claims.Email, "AIGetConversationMessages", map[string]interface{}{
		"conversationId": conversationID,
		"success":        true,
		"count":          len(messages),
	})

	sendAIAgentResponse(w, http.StatusOK, "Messages retrieved successfully.", messages)
}
