package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

const (
	PMErrLink     = "link"
	PMErrNode     = "node"
	PMErrDispatch = "dispatch"
	PMErrManifest = "manifest"
)

type DBPMError struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	PackageID uuid.UUID    `db:"package_id" json:"package_id"`
	Error     PackageError `db:"error" json:"error"`
	CreatedAt *CustomTime  `db:"created_at,omitempty" json:"created_at,omitempty"`
}

type PackageError struct {
	Type       string `json:"type"`
	StatusCode int    `json:"code"`
	Body       string `json:"body"`
}

func (g PackageError) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *PackageError) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &PackageError{}
	_ sql.Scanner   = &PackageError{}
)
