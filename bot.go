package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBBot struct {
	ID                uuid.UUID         `db:"id" json:"id"`
	Name              string            `db:"name" json:"name"`
	OrganizationID    uuid.UUID         `db:"organization_id" json:"organization_id"`
	InstalledPackages InstalledPackages `db:"installed_packages" json:"installed_packages"`
	CreatedAt         *CustomTime       `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt         *CustomTime       `db:"updated_at,omitempty" json:"updated_at"`
}


type BotsByOrganization struct {
	Organization DBOrganization `json:"organization"`
	Bots []DBBot `json:"bots"`
}

type DBBots []DBBot
func (b DBBots) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(b)
}
func (b *DBBots) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(b, src)
}
var (
	_ driver.Valuer = &DBBots{}
	_ sql.Scanner   = &DBBots{}
)

// APIBot is what is returned when fetching a single bot
type APIBot struct {
	*DBBot
	Environments DBEnvironments `db:"environments" json:"environments"`
	Packages     []Package      `db:"packages" json:"packages"`
}

func (b *APIBot) GetEnvironment(id uuid.UUID) *DBEnvironment {
	for _, e := range b.Environments {
		if e.ID == id {
			return &e
		}
	}

	return nil
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
