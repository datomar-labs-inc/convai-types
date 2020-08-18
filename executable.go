package ctypes

// Executable is the compiled binary required to run a bot
// It contains a manifest with all packages used, as well as all nodes, links, and modules
type Executable struct {
	Bot         CompiledBot `json:"bot"`
	ContextTree Context     `json:"context_tree"`
}
