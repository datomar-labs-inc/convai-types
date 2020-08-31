package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBEnvironment struct {
	ID          uuid.UUID                 `db:"id" json:"id"`
	Name        string                    `db:"name" json:"name"`
	BotID       uuid.UUID                 `db:"bot_id" json:"bot_id"`
	Data        map[uuid.UUID]interface{} `db:"data" json:"data"`
	BlueprintID *uuid.UUID                `db:"blueprint_id,omitempty" json:"blueprint_id,omitempty"`
	PromotedAt  *time.Time                `db:"promoted_at,omitempty" json:"promoted_at,omitempty"`
	IsDev       bool                      `db:"is_dev" json:"is_dev"`
}
