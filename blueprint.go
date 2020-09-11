package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBBlueprint struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	Modules   DBModuleList `db:"modules" json:"modules"`
	Version   string       `db:"version" json:"version"`
	BotID     uuid.UUID    `db:"bot_id" json:"bot_id"`
	CreatedAt *time.Time   `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time   `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type DBModuleList []DBModuleListItem

type DBModuleListItem struct {
	ModuleID uuid.UUID `db:"id" json:"id"`
	Version  string    `db:"version" json:"version"`
}
