package ctypes

const (
	NodeTypeGraphRef = iota
	NodeTypeAction
)

type CompiledBot struct {
	// The graph source map keeps track of each graph's original id
	// Modules get compiled into graphs for execution
	GraphSourceMap map[int]string `json:"module_source_map"`

	// Event nodes stores the graph/node combo for each event that can be handled by this bot
	EventNodes map[string][]int `json:"event_nodes"`

	// Graphs stores all available graphs
	Graphs []CompiledGraph `json:"graphs"`
}

type CompiledGraph struct {
	// The node source map keeps track of each node's original id
	NodeSourceMap map[int]string `json:"node_source_map"`

	// The link source map keeps track of each link's original id
	LinkSourceMap map[int]string `json:"link_source_map"`

	// Nodes is an array of nodes, where the id corresponds to the index of the node
	Nodes []Node `json:"nodes"`

	// Nodes is an array of links, where the id corresponds to the index of the link
	Links []Link `json:"links"`
}

type Node struct {
	Type int `json:"type"`

	// If this node is a graph ref node, this id should be populated
	Graph  *int        `json:"graph,omitempty"`
	Config interface{} `json:"config"`
}

type Link struct {
	Priority int         `json:"priority"`
	Config   interface{} `json:"config"`
}
