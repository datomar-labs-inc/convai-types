package ctypes

const (
	OpSet = iota
	OpDelete
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LogEntry struct {
	Message    string `json:"message"`
	Level      int    `json:"level"`
	ExecOffset int    `json:"exec_offset"` // The number of milliseconds since the start of node execution
}

type Transformation struct {
	MemoryContainerName string      `json:"memory_container_name"`
	Key                 string      `json:"key"`
	Value               interface{} `json:"value"`
	Operation           int         `json:"operation"`
}
