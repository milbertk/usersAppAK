package databaseusers

type FirebaseUserSearchResult struct {
	UID         string `json:"uid"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}
