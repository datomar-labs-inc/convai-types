package ctypes

import (
	"database/sql"
	"database/sql/driver"

	"upper.io/db.v3/postgresql"
)

type Mem map[string]interface{}

const (
	MCTypeExecution = iota
	MCTypeSession
	MCTypeContext
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

func (m *MemoryContainer) Transform(transformation Transformation) *MemoryContainer {
	switch transformation.Operation {
	case OpSet:
		m.Data[transformation.GetKey()] = transformation.Value
	case OpDelete:
		delete(m.Data, transformation.GetKey())
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

