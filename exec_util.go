package ctypes

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/opentracing-contrib/go-zap/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

const (
	OpSet = iota
	OpDelete
)

const (
	LogLevelTrace   = 10
	LogLevelDebug   = 20
	LogLevelInfo    = 30
	LogLevelWarning = 40
	LogLevelError   = 50
)

type ExecutionResult struct {
	ID             uuid.UUID     `json:"id"`
	EnvironmentID  uuid.UUID     `json:"environment_id"`
	BotID          uuid.UUID     `json:"bot_id"`
	StartTime      CustomTime    `json:"start_time"`
	FinishTime     CustomTime    `json:"finish_time"`
	Duration       time.Duration `json:"duration"`
	InitialContext *Context      `json:"initial_context"`
	Steps          []Step        `json:"steps"`
}

func (s *ExecutionResult) GetMemoryUpdates() []MemoryUpdate {
	var updates = make(map[string]MemoryUpdate)

	for _, step := range s.AllTransformations() {
		ctx, ok := s.InitialContext.GetContextByName(step.GetContextLevelName())
		if ok && ctx.ID != uuid.Nil {
			mc := ctx.GetMemoryContainerByName(step.GetMemoryContainerName())

			if mc != nil {
				key := fmt.Sprintf("%s.%s", step.GetContextLevelName(), step.GetMemoryContainerName())

				if ud, ok := updates[key]; !ok {
					updates[key] = MemoryUpdate{
						ContextID:       ctx.ID,
						EnvironmentID:   s.EnvironmentID,
						ContainerType:   mc.Type,
						ContainerName:   mc.Name,
						Transformations: []Transformation{step},
					}
				} else {
					ud.Transformations = append(updates[key].Transformations, step)
					updates[key] = ud
				}

			} else {
				// There was a memory modification to a non existent memory container :/
			}
		} else {
			// There was a memory modification to a non existent context :/
		}
	}

	var udList []MemoryUpdate

	for _, ud := range updates {
		udList = append(udList, ud)
	}

	return udList
}

func (s *ExecutionResult) AllTransformations() (transformations []Transformation) {
	return GetAllTransformations(s.Steps)
}

// Mongoify will return a version of an ExecutionResult with an _id field to override Mongo's default _id field
// Returns a new struct with a MongoID field. It looks all weird because the struct is being defined inline
func (s *ExecutionResult) Mongoify() bson.M {
	return mustMappify(&struct {
		MongoID uuid.UUID `json:"_id"`
		*ExecutionResult
	}{
		MongoID:         s.ID,
		ExecutionResult: s,
	})
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LogEntry struct {
	Message    string `json:"message"`
	Level      int    `json:"level"`
	ExecOffset int64  `json:"exec_offset"` // The number of milliseconds since the start of node execution
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
			log.Error("failed to match string regex while validating data path", zap.Error(err))
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
