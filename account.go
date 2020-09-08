package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type DBAccount struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Email       string     `db:"email" json:"email"`
	AccountType string     `db:"account_type" json:"account_type"`
	AccountKey  *string    `db:"account_key,omitempty" json:"account_key,omitempty"`
	Admin       bool       `db:"admin" json:"admin"`
	CreatedAt   *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}
