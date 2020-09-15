package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBEnvironment struct {
	ID          uuid.UUID                 `db:"id" json:"id"`
	Name        string                    `db:"name" json:"name"`
	BotID       uuid.UUID                 `db:"bot_id" json:"bot_id"`
	Data        map[uuid.UUID]interface{} `db:"data" json:"data"`
	BlueprintID *uuid.UUID                `db:"blueprint_id,omitempty" json:"blueprint_id,omitempty"`
	PromotedAt  *CustomTime               `db:"promoted_at,omitempty" json:"promoted_at,omitempty"`
	IsDev       bool                      `db:"is_dev" json:"is_dev"`
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

type UpdateEnvironmentPackageConfigRequest struct {
	PackageID uuid.UUID   `json:"package_id"`
	Data      interface{} `json:"data"`
}
