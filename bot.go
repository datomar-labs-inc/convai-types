package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBBot struct {
	ID                uuid.UUID         `db:"id" json:"id"`
	Name              string            `db:"name" json:"name"`
	OrganizationID    uuid.UUID         `db:"organization_id" json:"organization_id"`
	InstalledPackages InstalledPackages `db:"installed_packages" json:"installed_packages"`
	CreatedAt         *time.Time        `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt         *time.Time        `db:"updated_at,omitempty" json:"updated_at"`
}

type CreateBotRequest struct {
	Name string `json:"name" validate:"required,max=35,min=2"`
}

type InstalledPackages struct {
	Packages []InstalledPackage `json:"package_ids"`
}

type InstalledPackage struct {
	ID uuid.UUID `json:"id"`
}

func (g InstalledPackages) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *InstalledPackages) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &InstalledPackages{}
	_ sql.Scanner   = &InstalledPackages{}
)
