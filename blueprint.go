package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBBlueprint struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	Modules   DBModuleList `db:"modules" json:"modules"`
	Version   string       `db:"version" json:"version"`
	BotID     uuid.UUID    `db:"bot_id" json:"bot_id"`
	CreatedAt *CustomTime  `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *CustomTime  `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (g DBBlueprint) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBBlueprint) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DBBlueprint{}
	_ sql.Scanner   = &DBBlueprint{}
)

type DBModuleList []DBModuleListItem

func (g DBModuleList) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBModuleList) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DBModuleList{}
	_ sql.Scanner   = &DBModuleList{}
)

type DBModuleListItem struct {
	ModuleID uuid.UUID `db:"id" json:"id"`
	Version  string    `db:"version" json:"version"`
}
