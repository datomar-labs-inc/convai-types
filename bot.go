package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBBot struct {
	ID             uuid.UUID   `db:"id" json:"id"`
	Name           string      `db:"name" json:"name"`
	OrganizationID uuid.UUID   `db:"organization_id" json:"organization_id"`
	PackageIDs     []uuid.UUID `db:"package_ids" json:"package_ids"`
	CreatedAt      *time.Time  `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt      *time.Time  `db:"updated_at,omitempty" json:"updated_at"`
}
