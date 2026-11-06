package databaseuserofcompany

import (
	"time"
)

type CompanyUser struct {
	CompanyUserID string     `json:"companyUserId"`
	CompanyID     string     `json:"companyId"`
	UserID        string     `json:"userId"`
	Role          string     `json:"role"`
	Active        bool       `json:"active"`
	JoinedAt      time.Time  `json:"joinedAt"`
	RemovedAt     *time.Time `json:"removedAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

type CompanyUserSearchResult struct {
	CompanyUserID string  `json:"companyUserId"`
	CompanyID     string  `json:"companyId"`
	UserID        string  `json:"userId"`
	UserName      string  `json:"userName"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Active        bool    `json:"active"`
	JoinedAt      *string `json:"joinedAt,omitempty"`
	RemovedAt     *string `json:"removedAt,omitempty"`
	CreatedAt     *string `json:"createdAt,omitempty"`
	UpdatedAt     *string `json:"updatedAt,omitempty"`
}
