package ctypes

import (
	"fmt"
	"regexp"
	"strings"
)

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
	Path      string      `json:"path"`
	Value     interface{} `json:"value"`
	Operation int         `json:"operation"`
}

func (t *Transformation) PathValid() bool {
	parts := strings.Split(t.Path, ".")

	if len(parts) != 3 {
		return false
	}

	for _, p := range parts {
		matched, err := regexp.MatchString("^[a-zA-Z_$@\\-]{1}[a-zA-Z_$@\\-0-9]+$", p)
		if err != nil {
			fmt.Println("Failed to regex", err.Error())
			return false
		}

		if !matched {
			return false
		}
	}

	return true
}

func (t *Transformation) GetContextLevelName() string {
	return strings.TrimSpace(strings.Split(t.Path, ".")[0])
}

func (t *Transformation) GetMemoryContainerName() string {
	return strings.TrimSpace(strings.Split(t.Path, ".")[1])
}

func (t *Transformation) GetKey() string {
	return strings.TrimSpace(strings.Split(t.Path, ".")[2])
}