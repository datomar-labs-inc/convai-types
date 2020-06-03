package ctypes

type Package struct {
	Nodes      []PackageNode     `json:"nodes"`
	Links      []PackageLink     `json:"links"`
	Events     []RunnableEvent   `json:"events"`
	Dispatches []PackageDispatch `json:"dispatches"`
}
