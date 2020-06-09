package ctypes

import (
	"github.com/google/uuid"
)

type DBEvent struct {
	ID        string    `db:"id" json:"id"`
	PackageID uuid.UUID `db:"package_id" json:"package_id"`
	Name      string    `db:"name" json:"name"`
	Docs      string    `db:"docs" json:"docs"`
	Style     NodeStyle `db:"style" json:"style"`
}
