package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBBlueprint struct {
	ID        uuid.UUID    `json:"id"`
	Modules   DBModuleList `json:"modules"`
	Version   string       `json:"version"`
	BotID     uuid.UUID    `json:"bot_id"`
	CreatedAt *time.Time   `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time   `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type DBModuleList []DBModuleListItem

type DBModuleListItem struct {
	ModuleID uuid.UUID `json:"id"`
	Version  string    `json:"version"`
}
