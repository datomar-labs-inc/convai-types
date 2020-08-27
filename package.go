package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBPackage struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	BaseURL        string    `db:"base_url" json:"base_url"`
	SigningKey     string    `db:"signing_key" json:"signing_key,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type Package struct {
	DBPackage

	Nodes      []DBNode     `json:"nodes"`
	Links      []DBLink     `json:"links"`
	Events     []DBEvent    `json:"events"`
	Dispatches []DBDispatch `json:"dispatches"`
}
