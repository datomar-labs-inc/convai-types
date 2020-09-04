package ctypes

import (
	"time"

	"github.com/google/uuid"
)

type ExecutionOptions struct {
	MaxStackSize                  int // The maximum allowed stack size
	MaximumNodeCount              int // The maximum number of nodes that can be executed
	MaximumNodeExecutionDuration  time.Duration
	MaximumLinkEvaluationDuration time.Duration
	Timeout                       time.Duration

	Event string
	Data  interface{}
}

type Frame struct {
	ModuleID uuid.UUID
	NodeID   uuid.UUID
	Data     interface{}
}

type Step struct {
	Node     *NodeExecutionResult
	Links    map[uuid.UUID]LinkEvaluationResult
	Duration time.Duration
}

func GetAllTransformations(steps []Step) (transformations []Transformation) {
	for _, step := range steps {
		if step.Node != nil && step.Node.Transformations != nil {
			transformations = append(transformations, step.Node.Transformations...)
		}
	}

	return
}

type NodeExecutionResult struct {
	Transformations []Transformation `json:"transformations"`
	Logs            []LogEntry       `json:"logs"`
	Errors          []Error          `json:"errors"`
	HaltExecution   bool             `json:"halt_execution"`
	GoTo            *GoTo            `json:"go_to"`
}

type GoTo struct {
	ModuleID uuid.UUID `json:"module_id"`
	NodeID   uuid.UUID `json:"node_id"`
	Stack    []Frame   `json:"stack"`
}

type LinkEvaluationResult struct {
	Logs     []LogEntry        `json:"logs"`
	Errors   []Error           `json:"errors"`
	Passable bool              `json:"passable"` // Can the execution of this module proceed down this link
	Link     CompiledGraphLink `json:"link"`
}
