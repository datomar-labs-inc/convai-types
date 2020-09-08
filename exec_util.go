package ctypes

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	OpSet = iota
	OpDelete
)

type ExecutionResult struct {
	ID             uuid.UUID     `json:"id"`
	EnvironmentID  uuid.UUID     `json:"environment_id"`
	BotID          uuid.UUID     `json:"bot_id"`
	StartTime      time.Time     `json:"start_time"`
	FinishTime     time.Time     `json:"finish_time"`
	Duration       time.Duration `json:"duration"`
	InitialContext *Context      `json:"initial_context"`
	Steps          []Step        `json:"steps"`
}

func (s *ExecutionResult) AllTransformations() (transformations []Transformation) {
	return GetAllTransformations(s.Steps)
}

// Mongoify will return a version of an ExecutionResult with an _id field to override Mongo's default _id field
// Returns a new struct with a MongoID field. It looks all weird because the struct is being defined inline
func (s *ExecutionResult) Mongoify() *struct {
	MongoID uuid.UUID `json:"_id"`
	*ExecutionResult
} {
	return &struct {
		MongoID uuid.UUID `json:"_id"`
		*ExecutionResult
	}{
		MongoID:         s.ID,
		ExecutionResult: s,
	}
}

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

	if len(parts) < 3 {
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
	key := ""

	parts := strings.Split(path, ".")[2:]

	for i, str := range parts {
		key += strings.TrimSpace(str)

		if i != len(parts)-1 {
			key += "."
		}
	}

	return key
}

func DataPathHasMultipartKey(path string) bool {
	return len(GetDataPathKeyParts(path)) > 1
}

func GetDataPathKeyParts(path string) []string {
	return strings.Split(path, ".")[2:]
}
