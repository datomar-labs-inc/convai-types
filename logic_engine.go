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
	InitialTransforms             []Transformation // These transformations should be applied as the first step during execution

	Event string
	Data  interface{}
}

type Frame struct {
	ModuleID uuid.UUID
	NodeID   uuid.UUID
	Data     interface{}
}

type Step struct {
	ModuleID *uuid.UUID                         `json:"module_id,omitempty"`
	Node     *NodeExecutionResult               `json:"node,omitempty"`
	Links    map[uuid.UUID]LinkEvaluationResult `json:"links,omitempty"`
	Duration time.Duration                      `json:"duration"`
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
	NodeID          *uuid.UUID       `json:"node_id,omitempty"`
	Transformations []Transformation `json:"transformations,omitempty"`
	Logs            []LogEntry       `json:"logs,omitempty"`
	Errors          []Error          `json:"errors,omitempty"`
	HaltExecution   bool             `json:"halt_execution,omitempty"`
	GoTo            *GoTo            `json:"go_to,omitempty"`
}

type GoTo struct {
	ModuleID uuid.UUID `json:"module_id"`
	NodeID   uuid.UUID `json:"node_id"`
	Stack    []Frame   `json:"stack"`
}

type LinkEvaluationResult struct {
	Logs     []LogEntry        `json:"logs,omitempty"`
	Errors   []Error           `json:"errors,omitempty"`
	Passable bool              `json:"passable"` // Can the execution of this module proceed down this link
	Link     CompiledGraphLink `json:"-"`
}
