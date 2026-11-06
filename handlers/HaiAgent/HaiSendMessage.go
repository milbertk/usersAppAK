package haiagent

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/milbertk/security"

	baiagent "github.com/milbertk/usersAppAK/brules/bAIAgent"
)

type AISendMessageRequest struct {
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
	Platform       string `json:"platform"`
	CompanyID      string `json:"companyId"`
	CompanyName    string `json:"companyName"`
	CurrentPath    string `json:"currentPath"`
	Language       string `json:"language"`
}

// HaiSendMessage receives a message from the frontend, sends it to the
// Claude agent and returns the answer. Only for logged users.
func HaiSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendAIAgentResponse(w, http.StatusMethodNotAllowed, "Use POST", nil)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendAIAgentResponse(w, http.StatusUnauthorized, "Missing user claims", nil)
		return
	}

	loggedUID := strings.TrimSpace(claims.UID)

	var req AISendMessageRequest

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		insertAIAgentTracking(r, loggedUID, claims.Email, "AIChatSendMessage", map[string]interface{}{
			"success": false,
			"message": "Invalid JSON body: " + err.Error(),
		})
		sendAIAgentResponse(w, http.StatusBadRequest, "Invalid JSON body: "+err.Error(), nil)
		return
	}

	result, err := baiagent.BSendAIMessage(r.Context(), baiagent.SendAIMessageRequest{
		ConversationID: req.ConversationID,
		Message:        req.Message,
		Platform:       req.Platform,
		CompanyID:      req.CompanyID,
		CompanyName:    req.CompanyName,
		CurrentPath:    req.CurrentPath,
		Language:       req.Language,
		LoggedUID:      loggedUID,
		LoggedEmail:    claims.Email,
	})

	trackPayload := map[string]interface{}{
		"conversationId": strings.TrimSpace(req.ConversationID),
		"platform":       strings.TrimSpace(req.Platform),
		"companyId":      strings.TrimSpace(req.CompanyID),
		"currentPath":    strings.TrimSpace(req.CurrentPath),
		"language":       strings.TrimSpace(req.Language),
		"messageLength":  utf8.RuneCountInString(strings.TrimSpace(req.Message)),
	}

	if err != nil {
		trackPayload["success"] = false
		trackPayload["message"] = err.Error()
		insertAIAgentTracking(r, loggedUID, claims.Email, "AIChatSendMessage", trackPayload)

		sendAIAgentResponse(w, statusFromAIError(err.Error()), err.Error(), nil)
		return
	}

	trackPayload["success"] = true
	trackPayload["conversationId"] = result.ConversationID
	trackPayload["newConversation"] = result.NewConversation
	trackPayload["model"] = result.Model
	trackPayload["inputTokens"] = result.InputTokens
	trackPayload["outputTokens"] = result.OutputTokens
	trackPayload["latencyMs"] = result.LatencyMs
	trackPayload["stopReason"] = result.StopReason
	insertAIAgentTracking(r, loggedUID, claims.Email, "AIChatSendMessage", trackPayload)

	sendAIAgentResponse(w, http.StatusOK, "AI answer generated successfully.", result)
}
