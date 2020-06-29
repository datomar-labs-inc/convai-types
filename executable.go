package ctypes

import (
	"github.com/google/uuid"
)

// Executable is the compiled binary required to run a bot
// It contains a manifest with all packages used, as well as all nodes, links, and modules
type Executable struct {
	Manifest ExecutableManifest        `json:"manifest" msgpack:"m"`
	Modules  map[uuid.UUID]GraphModule `json:"modules" msgpack:"d"`
}

// ExecutableManifest contains a list of all plugins referenced by this executable
type ExecutableManifest struct {
	PackageIDs []uuid.UUID `json:"package_ids" msgpack:"i"`
}
