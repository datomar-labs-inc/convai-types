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
	PackageID   *uuid.UUID  `db:"package_id,omitempty" json:"package_id,omitempty"` // If this account belongs to a package
	Admin       bool        `db:"admin" json:"admin"`
	CreatedAt   *CustomTime `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   *CustomTime `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type AccountAndPermissionInfo struct {
	DBAccount
	Organizations DBOrganizationAccounts `db:"organizations" json:"organizations"`
	Bots          DBBotAccounts          `db:"bots" json:"bots"`
}

func (a *AccountAndPermissionInfo) UserHasPermissionLevelForOrganization(permLevel PermissionLevel, orgID uuid.UUID) bool {
	return a.GetPermissionLevelForOrganization(orgID) >= permLevel
}

// Contains tells whether a contains x.
func (a *AccountAndPermissionInfo) GetPermissionLevelForOrganization(orgID uuid.UUID) PermissionLevel {
	for _, op := range a.Organizations {
		if op.OrganizationID == orgID {
			return op.Permissions
		}
	}
	return PermUnauthorized
}

func (a *AccountAndPermissionInfo) UserHasPermissionLevelForBot(permLevel PermissionLevel, botID uuid.UUID) bool {
	return a.GetPermissionLevelForBot(botID) >= permLevel
}

func (a *AccountAndPermissionInfo) GetPermissionLevelForBot(botID uuid.UUID) PermissionLevel {
	for _, bp := range a.Bots {
		if bp.BotID == botID {
			return bp.Permissions
		}
	}
	return PermUnauthorized
}

type PermissionLevel int

const (
	PermUnauthorized    PermissionLevel = -1
	PermReadOnly                        = 0
	PermReadWrite                       = 1
	PermReadWriteDelete                 = 2
)

type DBOrganizationAccount struct {
	AccountID      uuid.UUID       `db:"account_id" json:"account_id"`
	OrganizationID uuid.UUID       `db:"organization_id" json:"organization_id"`
	Permissions    PermissionLevel `db:"permissions" json:"permissions"`
}

type DBOrganizationAccounts []DBOrganizationAccount

func (a DBOrganizationAccounts) OrgIDList() (ids []uuid.UUID) {
	for _, oa := range a {
		ids = append(ids, oa.OrganizationID)
	}

	return
}

func (a DBOrganizationAccounts) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(a)
}

func (a *DBOrganizationAccounts) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(a, src)
}

type DBBotAccount struct {
	AccountID   uuid.UUID       `db:"account_id" json:"account_id"`
	BotID       uuid.UUID       `db:"bot_id" json:"bot_id"`
	Permissions PermissionLevel `db:"permissions" json:"permissions"`
}

type DBBotAccounts []DBBotAccount

func (a DBBotAccounts) BotIDList() (ids []uuid.UUID) {
	for _, oa := range a {
		ids = append(ids, oa.BotID)
	}

	return
}

func (a DBBotAccounts) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(a)
}

func (a *DBBotAccounts) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(a, src)
}

var (
	_ driver.Valuer = &DBBotAccounts{}
	_ sql.Scanner   = &DBBotAccounts{}
)

var (
	_ driver.Valuer = &DBOrganizationAccounts{}
	_ sql.Scanner   = &DBOrganizationAccounts{}
)
