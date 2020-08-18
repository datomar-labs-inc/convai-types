package ctypes

import (
	"github.com/google/uuid"
)

// Location reference constants are used in the LocationReference structs to note the type of location being referenced
const (
	LRTypeConfig    = "config"
	LRTypePosition  = "position"
	LRTypeModule    = "module"
	LRTypeDuplicate = "duplicate"
)

type CompiledBot struct {
	// Event nodes stores the graph/node combo for each event that can be handled by this bot
	EventNodes map[string][]uuid.UUID `json:"event_nodes"`

	// PackageIDs stores a list of all packages referenced by this compiled bot
	PackageIDs []uuid.UUID            `json:"package_ids"`

	// Modules stores all available modules
	Modules map[uuid.UUID]CompiledGraphModule `json:"modules"`
}

type CompilerResult struct {
	Info     []CompilationNote `json:"info"`
	Warnings []CompilationNote `json:"warnings"`
	Errors   []CompilationNote `json:"errors"`
}

type CompilationNote struct {
	Message string                     `json:"message"`
	Code    int                        `json:"code"`
	PLR     []PackageLocationReference `json:"plr,omitempty"`
	GLR     []GraphLocationReference   `json:"glr,omitempty"`
}

type PackageLocationReference struct {
	PackageID uuid.UUID `json:"package_id"`
}

type GraphLocationReference struct {
	ModuleID uuid.UUID  `json:"module_id"`
	NodeID   *uuid.UUID `json:"node_id,omitempty"`
	LinkID   *uuid.UUID `json:"link_id,omitempty"`
	Type     string     `json:"type"`
}
