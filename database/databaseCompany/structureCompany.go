package databasecompany

import "time"

type Company struct {
	CompanyID string `db:"company_id"`

	OwnerFirebaseUID string `db:"owner_firebase_uid"`

	CompanyName          string `db:"company_name"`
	LegalName            string `db:"legal_name"`
	IdentificationType   string `db:"identification_type"`
	IdentificationNumber string `db:"identification_number"`

	Email   string `db:"email"`
	Phone   string `db:"phone"`
	Website string `db:"website"`

	CountryCode string `db:"country_code"`
	CountryName string `db:"country_name"`

	StateCode string `db:"state_code"`
	StateName string `db:"state_name"`

	City       string `db:"city"`
	District   string `db:"district"`
	PostalCode string `db:"postal_code"`

	Address1   string `db:"address1"`
	Address2   string `db:"address2"`
	Directions string `db:"directions"`

	Category         string `db:"category"`
	InternalCategory string `db:"internal_category"`

	Employees int `db:"employees"`

	IsVirtual bool `db:"is_virtual"`
	Active    bool `db:"active"`

	Plan   string `db:"plan"`
	Status string `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
