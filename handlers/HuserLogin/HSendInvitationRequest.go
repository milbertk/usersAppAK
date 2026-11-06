package huserlogin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/milbertk/security"

	bcreateusers "github.com/milbertk/usersAppAK/brules/bcreateUsers"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type SendInvitationRequest struct {
	TargetEmail string `json:"targetEmail"`
	Message     string `json:"message"`
}

type SendInvitationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HSendInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeInvitationJSON(w, http.StatusMethodNotAllowed, SendInvitationResponse{
			Status:  "MethodNotAllowed",
			Message: "Use POST",
		})
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		writeInvitationJSON(w, http.StatusUnauthorized, SendInvitationResponse{
			Status:  "Unauthorized",
			Message: "Missing user claims",
		})
		return
	}

	currentUser, err := databaseusers.SelectByUID(claims.UID)
	if err != nil {
		fmt.Println("SelectByUID error:", err)

		writeInvitationJSON(w, http.StatusUnauthorized, SendInvitationResponse{
			Status:  "Unauthorized",
			Message: "Current user not found",
		})
		return
	}

	var req SendInvitationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeInvitationJSON(w, http.StatusBadRequest, SendInvitationResponse{
			Status:  "BadRequest",
			Message: "Invalid JSON body",
		})
		return
	}

	req.TargetEmail = strings.TrimSpace(req.TargetEmail)
	req.Message = strings.TrimSpace(req.Message)

	if req.TargetEmail == "" {
		writeInvitationJSON(w, http.StatusBadRequest, SendInvitationResponse{
			Status:  "BadRequest",
			Message: "Target email is required",
		})
		return
	}

	if req.Message == "" {
		req.Message = "You have been invited to join the platform."
	}

	err = bcreateusers.BSendInvitationEmail(
		*currentUser,
		req.TargetEmail,
		req.Message,
	)

	if err != nil {
		fmt.Println("BSendInvitationEmail error:", err)

		writeInvitationJSON(w, http.StatusInternalServerError, SendInvitationResponse{
			Status:  "Error",
			Message: common.ClientMessage(err, "The invitation could not be sent"),
		})
		return
	}

	writeInvitationJSON(w, http.StatusOK, SendInvitationResponse{
		Status:  "OK",
		Message: "Invitation sent successfully",
	})
}

func writeInvitationJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
