package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type DBNode struct {
	ID        string    `db:"id" json:"id"`
	PackageID uuid.UUID `db:"package_id" json:"package_id"`
	Version   string    `db:"verison" json:"version"`
	Name      string    `db:"name" json:"name"`
	Docs      string    `db:"docs" json:"docs"`
	Style     NodeStyle `db:"style" json:"style"`
}

type NodeStyle struct {
	Color string   `json:"color"` // Valid hex code color
	Icons []string `json:"icons"` // File name (files will be served in a special format by the plugin)
}

// NodeCall is Convai requesting that a package perform a node execution and return the result
type NodeCall struct {
	RequestID       uuid.UUID         `json:"request_id"` // The ID of the current request
	ID              string            `json:"id"`         // The ID of the node type, used by the plugin to determine which node
	Version         string            `json:"version"`    // Which version of this node was this config created on
	Config          MemoryContainer   `json:"config"`     // How this specific node was configured by the bot builder
	PackageSettings MemoryContainer   `json:"package_settings"`
	Memory          []MemoryContainer `json:"memory"`   // Any other memory containers that this package is allowed to see
	Sequence        int               `json:"sequence"` // The number of nodes that have been executed during this execution
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
	ID            string    `json:"id"`
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
