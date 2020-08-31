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
	return ValidateDataPath(t.Path)
}

func (t *Transformation) GetContextLevelName() string {
	return GetDataPathContextLevelName(t.Path)
}

func (t *Transformation) GetMemoryContainerName() string {
	return GetDataPathMemoryContainerName(t.Path)
}

func (t *Transformation) GetKey() string {
	return GetDataPathKey(t.Path)
}

func ValidateDataPath(path string) bool {
	parts := strings.Split(path, ".")

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

func GetDataPathContextLevelName(path string) string {
	return strings.TrimSpace(strings.Split(path, ".")[0])
}

func GetDataPathMemoryContainerName(path string) string {
	return strings.TrimSpace(strings.Split(path, ".")[1])
}

func GetDataPathKey(path string) string {
	return strings.TrimSpace(strings.Split(path, ".")[2])
}

