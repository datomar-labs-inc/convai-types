package ctypes

import (
	"github.com/google/uuid"
)

type PackageLink struct {
	Name          string    `json:"name"`
	ID            string    `json:"id"`
	Version       string    `json:"version"` // Valid semantic version
	Style         LinkStyle `json:"style"`
	Documentation string    `json:"documentation"` // Markdown format
}

// LinkCall is Convai requesting that a package perform a link execution and return the result
type LinkCall struct {
	RequestID       uuid.UUID         `json:"request_id"` // The id of the current request
	ID              string            `json:"id"`         // The ID of the link type, used by the plugin to determine which link should be executed
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