package hcompany

import (
	"encoding/json"
	"net/http"
	"strings"

	buserofcompany "github.com/milbertk/usersAppAK/brules/busersOfCompany"
	"github.com/milbertk/usersAppAK/helpers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type UpdateCompanyUserResponse struct {
	Status  string                                  `json:"status"`
	Message string                                  `json:"message"`
	Data    *buserofcompany.UpdateCompanyUserResult `json:"data,omitempty"`
}

func HUpdateCompanyUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		sendUpdateCompanyUserResponseuc(
			w,
			http.StatusMethodNotAllowed,
			"ERROR",
			"method not allowed",
			nil,
		)
		return
	}

	requesterUserID, err := helpers.RequireLocalJWT(r)
	if err != nil {
		sendUpdateCompanyUserResponseuc(
			w,
			http.StatusUnauthorized,
			"ERROR",
			"unauthorized",
			nil,
		)
		return
	}

	var input buserofcompany.UpdateCompanyUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendUpdateCompanyUserResponseuc(
			w,
			http.StatusBadRequest,
			"ERROR",
			"invalid request body",
			nil,
		)
		return
	}

	result, err := buserofcompany.BUpdateCompanyUser(input, requesterUserID)
	if err != nil {
		message := err.Error()
		lowerMessage := strings.ToLower(message)

		switch {
		case strings.Contains(lowerMessage, "unauthorized"):
			sendUpdateCompanyUserResponseuc(
				w,
				http.StatusForbidden,
				"ERROR",
				message,
				nil,
			)

		case strings.Contains(lowerMessage, "required"),
			strings.Contains(lowerMessage, "invalid"):
			sendUpdateCompanyUserResponseuc(
				w,
				http.StatusBadRequest,
				"ERROR",
				message,
				nil,
			)

		default:
			sendUpdateCompanyUserResponseuc(
				w,
				http.StatusInternalServerError,
				"ERROR",
				common.InternalMessage(err, "The company user could not be updated"),
				nil,
			)
		}

		return
	}

	sendUpdateCompanyUserResponseuc(
		w,
		http.StatusOK,
		"OK",
		"company user updated successfully",
		result,
	)
}

func sendUpdateCompanyUserResponseuc(
	w http.ResponseWriter,
	httpStatus int,
	status string,
	message string,
	data *buserofcompany.UpdateCompanyUserResult,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	_ = json.NewEncoder(w).Encode(UpdateCompanyUserResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
