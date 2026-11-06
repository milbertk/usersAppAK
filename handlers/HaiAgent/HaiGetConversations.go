package haiagent

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/milbertk/security"

	baiagent "github.com/milbertk/usersAppAK/brules/bAIAgent"
)

// HaiGetConversations returns the conversations of the logged user.
// Optional query: ?limit=50
func HaiGetConversations(w http.ResponseWriter, r *http.Request) {
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

	limit, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("limit")))

	conversations, err := baiagent.BGetAIConversations(loggedUID, limit)
	if err != nil {
		insertAIAgentTracking(r, loggedUID, claims.Email, "AIGetConversations", map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		sendAIAgentResponse(w, statusFromAIError(err.Error()), err.Error(), nil)
		return
	}

	insertAIAgentTracking(r, loggedUID, claims.Email, "AIGetConversations", map[string]interface{}{
		"success": true,
		"count":   len(conversations),
	})

	sendAIAgentResponse(w, http.StatusOK, "Conversations retrieved successfully.", conversations)
}
