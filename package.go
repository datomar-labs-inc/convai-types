package ctypes

import (
	"github.com/google/uuid"
)

type DBPackage struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	BaseURL        string    `db:"base_url" json:"base_url"`
	SigningKey     string    `db:"signing_key" json:"signing_key"`
}

type Package struct {
	Nodes      []PackageNode     `json:"nodes"`
	Links      []PackageLink     `json:"links"`
	Events     []RunnableEvent   `json:"events"`
	Dispatches []PackageDispatch `json:"dispatches"`
}
