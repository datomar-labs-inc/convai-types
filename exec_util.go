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
	Key       string `json:"key"`
	Value     string `json:"value"`
	Operation int    `json:"operation"`
}
