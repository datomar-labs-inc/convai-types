package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

type MemoryUpdate struct {
	ContextID       uuid.UUID        `json:"i"`
	EnvironmentID   uuid.UUID        `json:"e"`
	ContainerType   int              `json:"ct"`
	ContainerName   string           `json:"c"`
	Transformations []Transformation `json:"t"`
}

type TransformMemoryInput struct {
	Transformations []Transformation `json:"transformations"`
}

type Mem map[string]interface{}

func (m Mem) ToTransformationsPrefixed(pathPrefix string) (trs []Transformation) {
	for k, v := range m {
		trs = append(trs, Transformation{
			Path:      fmt.Sprintf("%s.%s", pathPrefix, k),
			Value:     v,
			Operation: OpSet,
		})
	}

	return
}

func (m Mem) ToTransformations() (trs []Transformation) {
	return m.ToTransformationsPrefixed(".container")
}

func (m Mem) Transform(transformation ...Transformation) Mem {
	for _, t := range transformation {
		switch t.Operation {
		case OpSet:
			m[t.GetKey()] = t.Value
		case OpDelete:
			delete(m, t.GetKey())
		}
	}

	return m
}

const (
	MCTypeExecution = iota // Execution type memory can be modified, and is only stored in the execution log
	MCTypeSession          // Session type memory can be modified, and is stored in redis
	MCTypeContext          // Context type memory can be modified, and is stored in mongo
	MCTypeReadOnly         // ReadOnly type memory cannot be modified, and is only stored in the execution log
	MCTypeSecure           // Secure type memory can be modified, and is stored in the vault
)

type MemoryContainer struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Exposed bool   `json:"exposed"`
	Data    Mem    `json:"data"`
}

type DBMemoryContainer struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Exposed bool   `json:"exposed"`
}

type DBMemoryContainers []DBMemoryContainer

func (m *MemoryContainer) Put(key string, value interface{}) *MemoryContainer {
	m.Data[key] = value
	return m
}

func (m *MemoryContainer) Transform(transformation ...Transformation) *MemoryContainer {
	// If this memory container is not allowed to be modified, do not modify it
	if m.Type == MCTypeReadOnly {
		return m
	}

	if m.Data == nil {
		m.Data = map[string]interface{}{}
	}

	for _, t := range transformation {
		m.Data = m.Data.Transform(t)
	}

	return m
}

func (g DBMemoryContainers) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DBMemoryContainers) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DBMemoryContainers{}
	_ sql.Scanner   = &DBMemoryContainers{}
)
