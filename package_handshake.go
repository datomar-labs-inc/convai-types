package ctypes

type NodeStyle struct {
	Color string   `json:"color"` // Valid hex code color
	Icons []string `json:"icons"` // File name (files will be served in a special format by the plugin)
}

type LinkStyle struct {
	Color string   `json:"color"` // Valid hex code color
	Icons []string `json:"icons"` // File name (files will be served in a special format by the plugin)
}

type RunnableEvent struct {
	Name          string    `json:"name"`
	ID            string    `json:"id"`
	Documentation string    `json:"documentation"` // Markdown format
	Style         NodeStyle `json:"style"`
}

type DispatchStyle struct {
	Color string `json:"color"` // Valid hex code color
	Icon  string `json:"icon"`  // File name (files will be served in a special format by the plugin)
}

type PackageModule struct {
}

type PackageTemplates struct {
}
