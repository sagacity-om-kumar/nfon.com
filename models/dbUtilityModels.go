package models

import "time"

type AuditLogRowDataModel struct {
	ID          int       `db:"id"`
	Module      string    `db:"module"`
	API         string    `db:"api"`
	Page        string    `db:"page"`
	AccessedBy  int       `db:"accessed_by"`
	AccessedDtm time.Time `db:"accessed_dtm"`
}
