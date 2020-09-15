package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBAccount struct {
	ID          uuid.UUID   `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Email       string      `db:"email" json:"email"`
	AccountType string      `db:"account_type" json:"account_type"`
	AccountKey  *string     `db:"account_key,omitempty" json:"account_key,omitempty"`
	PackageID   *uuid.UUID  `db:"package_id,omitempty" json:"package_id,omitempty"`
	Admin       bool        `db:"admin" json:"admin"`
	CreatedAt   *CustomTime `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   *CustomTime `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type DBOrganizationAccount struct {
	AccountID      uuid.UUID `db:"account_id" json:"account_id"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	Permissions    int       `db:"permissions" json:"permissions"`
}

type DBOrganizationAccounts []DBOrganizationAccount

func (g DBOrganizationAccounts) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBOrganizationAccounts) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

type DBBotAccount struct {
	AccountID   uuid.UUID `db:"account_id" json:"account_id"`
	BotID       uuid.UUID `db:"bot_id" json:"bot_id"`
	Permissions int       `db:"permissions" json:"permissions"`
}

type DBBotAccounts []DBBotAccount

func (g DBBotAccounts) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBBotAccounts) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DBBotAccounts{}
	_ sql.Scanner   = &DBBotAccounts{}
)

var (
	_ driver.Valuer = &DBOrganizationAccounts{}
	_ sql.Scanner   = &DBOrganizationAccounts{}
)
