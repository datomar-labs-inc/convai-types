package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type PublishBlueprintRequest struct {
	BlueprintID uuid.UUID `json:"blueprint_id"`
	Version     *Semver   `json:"version,omitempty"` // If Version is nil, the latest version will be published

	// If NewVersion is not nil, the old version will be copied into a new unlocked version
	// If publishing the latest version, a new version must be provided
	// If not publishing the latest version, the new version must be below the latest version
	NewVersion *Semver `json:"new_version,omitempty"`
}

type DBBlueprint struct {
	ID                   uuid.UUID    `db:"id" json:"id"`
	Modules              DBModuleList `db:"modules" json:"modules"`
	Version              Semver       `db:"version" json:"version"`
	BotID                uuid.UUID    `db:"bot_id" json:"bot_id"`
	ContinuedFromVersion *Semver      `db:"continued_from_version,omitempty" json:"continued_from_version,omitempty"`
	Locked               bool         `db:"locked" json:"locked"`
	CreatedAt            *CustomTime  `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt            *CustomTime  `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type DBBlueprintGraphItem struct {
	Version   Semver      `db:"version" json:"version"`
	Parent    *Semver     `db:"parent,omitempty" json:"parent,omitempty"`
	Locked    bool        `db:"locked" json:"locked"`
	CreatedAt *CustomTime `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *CustomTime `db:"updated_at,omitempty" json:"updated_at,omitempty"`

	Children []DBBlueprintGraphItem `db:"-" json:"children"`
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

type DBModuleList map[uuid.UUID]DBModuleListItem

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

// DBModuleListItem is a single module that is locally scoped to a blueprint
type DBModuleListItem struct {
	ModuleID uuid.UUID   `json:"id"`
	Name     string      `json:"name"`
	Graph    GraphModule `json:"graph"`
}
