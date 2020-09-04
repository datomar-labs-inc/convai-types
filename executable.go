package ctypes

import (
	"github.com/google/uuid"
)

// Executable is the compiled binary required to run a bot
// It contains a manifest with all packages used, as well as all nodes, links, and modules
type Executable struct {
	Bot         CompiledBot `json:"bot"`
	ContextTree Context     `json:"context_tree"`
	Packages    []Package   `json:"packages"`
}

// ExecutionRequest is what the api will be called with when an execution should be performed
type ExecutionRequest struct {
	ID               uuid.UUID        `json:"id"`
	Event            string           `json:"event"`
	Mock             bool             `json:"mock"`
	Text             string           `json:"text"`
	DefaultDispatch  *string          `json:"default_dispatch,omitempty"`
	Transformations  []Transformation `json:"transformations"`
	ContextTreeSlice ContextTreeSlice `json:"context_tree_slice"`
}
