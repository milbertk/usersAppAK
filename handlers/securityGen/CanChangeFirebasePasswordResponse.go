package securitygen

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/milbertk/security"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type CanChangeFirebasePasswordResponse struct {
	Status       string   `json:"status"`
	Message      string   `json:"message"`
	CanChange    bool     `json:"canChange"`
	Providers    []string `json:"providers"`
	ProviderType string   `json:"providerType"`
}

func HCanChangeFirebasePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCanChangeFirebasePasswordResponse(
			w,
			http.StatusMethodNotAllowed,
			"Use GET",
			false,
			[]string{},
			"unknown",
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendCanChangeFirebasePasswordResponse(
			w,
			http.StatusUnauthorized,
			"Missing user claims",
			false,
			[]string{},
			"unknown",
		)
		return
	}

	userID := strings.TrimSpace(claims.UID)

	userDb, err := databaseusers.SelectByUID(userID)
	if err != nil {
		sendCanChangeFirebasePasswordResponse(
			w,
			http.StatusBadRequest,
			common.InternalMessage(err, "Error loading user"),
			false,
			[]string{},
			"unknown",
		)
		return
	}

	providerID := normalizeProviderID(userDb.ProviderID)

	providers := getCanChangeFirebaseProviderIDs(providerID)
	canChange := hasPasswordProvider(providerID)

	if canChange {
		sendCanChangeFirebasePasswordResponse(
			w,
			http.StatusOK,
			"Password change is available for this account.",
			true,
			providers,
			"password",
		)
		return
	}

	sendCanChangeFirebasePasswordResponse(
		w,
		http.StatusOK,
		"Password change is not available because this account uses an external provider.",
		false,
		providers,
		"external",
	)
}

func normalizeProviderID(providerID string) string {
	providerID = strings.ToLower(strings.TrimSpace(providerID))

	if providerID == "" {
		return "unknown"
	}

	return providerID
}

func hasPasswordProvider(providerID string) bool {
	providerID = normalizeProviderID(providerID)

	return providerID == "password"
}

func getCanChangeFirebaseProviderIDs(providerID string) []string {
	providerID = normalizeProviderID(providerID)

	return []string{providerID}
}

func sendCanChangeFirebasePasswordResponse(
	w http.ResponseWriter,
	status int,
	message string,
	canChange bool,
	providers []string,
	providerType string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if providers == nil {
		providers = []string{}
	}

	_ = json.NewEncoder(w).Encode(CanChangeFirebasePasswordResponse{
		Status:       http.StatusText(status),
		Message:      message,
		CanChange:    canChange,
		Providers:    providers,
		ProviderType: providerType,
	})
}
