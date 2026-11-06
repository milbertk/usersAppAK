package databaseuserofcompany

type CompanyByUserResult struct {
	CompanyID   string `json:"companyId"`
	CompanyName string `json:"companyName"`
	UserID      string `json:"userId"`
	Kind        string `json:"kind"`
}
