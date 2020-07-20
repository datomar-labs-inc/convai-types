package ctypes

import (
	"github.com/google/uuid"
)

type CompiledBot struct {
	// Event nodes stores the graph/node combo for each event that can be handled by this bot
	EventNodes map[string][]uuid.UUID `json:"event_nodes"`

	// Modules stores all available modules
	Modules map[uuid.UUID]GraphModule `json:"modules"`
}
