package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBEnvironment struct {
	ID               uuid.UUID       `db:"id" json:"id"`
	Name             string          `db:"name" json:"name"`
	BotID            uuid.UUID       `db:"bot_id" json:"bot_id"`
	Data             EnvironmentData `db:"data" json:"data"`
	BlueprintID      *uuid.UUID      `db:"blueprint_id,omitempty" json:"blueprint_id,omitempty"`
	BlueprintVersion *Semver         `db:"blueprint_version,omitempty" json:"blueprint_version,omitempty"`
	PromotedAt       *CustomTime     `db:"promoted_at,omitempty" json:"promoted_at,omitempty"`
	IsDev            bool            `db:"is_dev" json:"is_dev"`
}

// APIEnvironment includes the environment details, as well as the blueprint (if any)
type APIEnvironment struct {
	*DBEnvironment
	Blueprint *DBBlueprint `db:"blueprint" json:"blueprint"`
}

type DBEnvironments []DBEnvironment

func (g DBEnvironments) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBEnvironments) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DBEnvironments{}
	_ sql.Scanner   = &DBEnvironments{}
)

type EnvironmentData map[uuid.UUID]interface{}

func (g EnvironmentData) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *EnvironmentData) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &EnvironmentData{}
	_ sql.Scanner   = &EnvironmentData{}
)

type UpdateEnvironmentPackageConfigRequest struct {
	PackageID uuid.UUID   `json:"package_id"`
	Data      interface{} `json:"data"`
}
