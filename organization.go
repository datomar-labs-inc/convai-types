package ctypes

import (
	"github.com/google/uuid"
)

type DBOrganization struct {
	ID        uuid.UUID   `db:"id" json:"id"`
	Name      string      `db:"name" json:"name"`
	CreatedAt *CustomTime `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *CustomTime `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}
