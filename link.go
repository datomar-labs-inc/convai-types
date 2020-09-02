package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBLink struct {
	TypeID        string    `db:"id" json:"type_id"`
	PackageID     uuid.UUID `db:"package_id" json:"package_id"`
	Version       string    `db:"version" json:"version"`
	Name          string    `db:"name" json:"name"`
	Documentation string    `db:"docs" json:"docs"`
	Style         LinkStyle `db:"style" json:"style"`
}

// CompiledGraphLink is the variant of link that lives in a compiled executable
type CompiledGraphLink struct {
	ID          uuid.UUID `json:"id" msgpack:"i"`
	PackageID   uuid.UUID `json:"package_id" msgpack:"p"`
	TypeID      string    `json:"type_id" msgpack:"l"`
	Version     string    `json:"version" msgpack:"v"`
	Priority    int       `json:"priority"`
	Source      uuid.UUID `json:"source"`
	Destination uuid.UUID `json:"destination"`
	ConfigJSON  string    `json:"config_json" msgpack:"c"`
}

type LinkStyle struct {
	Color string   `json:"color"` // Valid hex code color
	Icons []string `json:"icons"` // File name (files will be served in a special format by the plugin)
}

type PackageLink struct {
	Name          string    `json:"name"`
	TypeID        string    `json:"type_id"`
	Version       string    `json:"version"` // Valid semantic version
	Style         LinkStyle `json:"style"`
	Documentation string    `json:"documentation"` // Markdown format
}

// LinkCall is Convai requesting that a package perform a link execution and return the result
type LinkCall struct {
	RequestID       uuid.UUID         `json:"request_id"` // The id of the current request
	TypeID          string            `json:"type_id"`    // The TypeID of the link type, used by the plugin to determine which link should be executed
	Version         string            `json:"version"`    // Which version of this link was this config created on
	Config          MemoryContainer   `json:"config"`     // How this specific link was configured by the bot builder
	PackageSettings MemoryContainer   `json:"package_settings"`
	Memory          []MemoryContainer `json:"memory"`   // Any other memory containers that this package is allowed to see
	Sequence        int               `json:"sequence"` // The number of links that have been executed during this execution
}

// LinkCallResult is what a package returns after executing a link
type LinkCallResult struct {
	RequestID uuid.UUID  `json:"request_id"` // The package is required to return the request id for security
	Logs      []LogEntry `json:"logs"`
	Errors    []Error    `json:"errors"`
	Passable  bool       `json:"passable"` // Can the execution of this module proceed down this link
}

// PackageLink stuff
// POST /links/:lid/execute
// POST /links/:lid/execute-mock
type LinkExecutionRequest struct {
	Calls []LinkCall `json:"calls"` // All links in a call should be processed concurrently
}

type LinkExecutionResponse struct {
	Results []LinkCallResult `json:"results"` // Results should be returned in the same order that the calls were provided
}

func (g LinkStyle) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *LinkStyle) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &LinkStyle{}
	_ sql.Scanner   = &LinkStyle{}
)
