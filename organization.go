package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBOrganization struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}
