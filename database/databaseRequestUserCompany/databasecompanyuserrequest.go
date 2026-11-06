package databaseRequestUserCompany

import "time"

type CompanyUserRequest struct {
	RequestID         string     `json:"requestId"`
	CompanyID         string     `json:"companyId"`
	InvitedUserID     string     `json:"invitedUserId"`
	RequestedByUserID string     `json:"requestedByUserId"`
	RoleRequested     string     `json:"roleRequested"`
	Status            string     `json:"status"`
	Message           string     `json:"message"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	ResolvedAt        *time.Time `json:"resolvedAt,omitempty"`
}

type CompanyUserRequestSearchResult struct {
	RequestID         string     `json:"requestId"`
	CompanyID         string     `json:"companyId"`
	InvitedUserID     string     `json:"invitedUserId"`
	RequestedByUserID string     `json:"requestedByUserId"`
	RoleRequested     string     `json:"roleRequested"`
	Status            string     `json:"status"`
	Message           string     `json:"message"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	ResolvedAt        *time.Time `json:"resolvedAt,omitempty"`

	InvitedUserName  string `json:"invitedUserName"`
	InvitedUserEmail string `json:"invitedUserEmail"`
}

type CompanyUserRequestWithCompany struct {
	RequestID         string     `json:"requestId"`
	CompanyID         string     `json:"companyId"`
	CompanyName       string     `json:"companyName"`
	InvitedUserID     string     `json:"invitedUserId"`
	RequestedByUserID string     `json:"requestedByUserId"`
	RoleRequested     string     `json:"roleRequested"`
	Status            string     `json:"status"`
	Message           string     `json:"message"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	ResolvedAt        *time.Time `json:"resolvedAt,omitempty"`
}
