package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBNode struct {
	TypeID        string    `db:"id" json:"type_id" validate:"required"`
	PackageID     uuid.UUID `db:"package_id" json:"package_id"`
	Version       string    `db:"version" json:"version"`
	Name          string    `db:"name" json:"name"`
	Documentation string    `db:"docs" json:"docs"`
	Style         NodeStyle `db:"style" json:"style"`
}

type CompiledGraphNode struct {
	ID uuid.UUID `json:"id" msgpack:"i"`

	// Properties that are common to all types of nodes
	PackageID uuid.UUID `json:"package_id,omitempty" msgpack:"p,omitempty"`

	// Properties of a node that references functionality in a package
	TypeID     *string `json:"type_id,omitempty" msgpack:"n,omitempty"`
	Version    *string `json:"version,omitempty" msgpack:"v,omitempty"`
	ConfigJSON *string `json:"config_json,omitempty" msgpack:"c,omitempty"`

	// Properties of a node that acts as a reference to a module
	ModuleID      *uuid.UUID `json:"module_id,omitempty" msgpack:"m,omitempty"`
	ModuleVersion *string    `json:"module_version,omitempty" msgpack:"mv,omitempty"`

	// Properties of a node that acts as an event entrypoint
	EventTypeID *string `json:"event_type_id,omitempty" msgapck:"e,omitempty"`
}

type NodeStyle struct {
	Color string   `json:"color" validate:"hexcolor"` // Valid hex code color
	Icons []string `json:"icons"`                     // File name (files will be served in a special format by the plugin)
}

// NodeCall is Convai requesting that a package perform a node execution and return the result
type NodeCall struct {
	RequestID       uuid.UUID `json:"request_id"`       // The TypeID of the current request
	TypeID          string    `json:"type_id"`          // The TypeID of the node type, used by the plugin to determine which node
	Version         string    `json:"version"`          // Which version of this node was this config created on
	Config          string    `json:"config"`           // How this specific node was configured by the bot builder (JSON format)
	PackageSettings string    `json:"package_settings"` // Settings for this package (JSON format)
	Tree            *Context   `json:"tree"`            // A context tree containing all data visible by the package
	Sequence        int       `json:"sequence"`         // The number of nodes that have been executed during this execution
}

// NodeCallResult is what a package returns after executing a node
type NodeCallResult struct {
	RequestID       uuid.UUID        `json:"request_id"` // The package is required to return the request id for security
	Transformations []Transformation `json:"transformations"`
	Logs            []LogEntry       `json:"logs"`
	Errors          []Error          `json:"errors"`
}

type NodeExecutionRequest struct {
	Calls []NodeCall `json:"calls"` // All nodes in a pack must complete execution before returning, so don't pack nodes
}

type NodeExecutionResponse struct {
	Results []NodeCallResult `json:"results"` // Results should be returned in the same order that the calls were provided
}

type PackageNode struct {
	Name          string    `json:"name"`
	TypeID        string    `json:"type_id"`
	Version       string    `json:"version"` // Valid semantic version
	Style         NodeStyle `json:"style"`
	Documentation string    `json:"documentation"` // Markdown format
}

func (g NodeStyle) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *NodeStyle) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &NodeStyle{}
	_ sql.Scanner   = &NodeStyle{}
)
