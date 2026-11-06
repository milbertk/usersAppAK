package databaseaccount

import "time"

type DatabaseAccountS struct {
	ID            string    `db:"id" json:"id"`
	OwnerID       string    `db:"ownerid" json:"ownerid"`
	CreateDate    time.Time `db:"createdate" json:"createdate"`
	StatusAccount string    `db:"statusaccount" json:"statusaccount"`
}
