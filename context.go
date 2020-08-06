package ctypes

import (
	"github.com/google/uuid"
)

// Context is the actual data format for a context, NOT DATABASE FRIENDLY
type Context struct {
	Name     string            `json:"name"`
	ID       uuid.UUID         `json:"id"`
	ParentID uuid.UUID         `json:"parent_id"`
	Ref      []string          `json:"ref"` // Maps to a ref table in the database, will be unpacked/queried into a slice
	Memory   []MemoryContainer `json:"memory"`

	// References
	Parent   *Context  `json:"-"`
	Children []Context `json:"-"`
}

type DBContext struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	ParentID         *uuid.UUID         `db:"parent_id,omitempty" json:"parent_id,omitempty"`
	Name             string             `db:"name" json:"name"`
	MemoryContainers DBMemoryContainers `db:"memory_containers" json:"memory_containers"`

	// The ref property is not present on the contexts table, but is frequently queried with it
	Ref *string `db:"ref,omitempty" json:"ref,omitempty"`

	Children []DBContext `db:"-" json:"children,omitempty"`
}

type DBContextRef struct {
	ContextID uuid.UUID `db:"context_id" json:"context_id"`
	Ref       string    `db:"ref" json:"ref"`
}

// ContextTreeSlice is used to rectify a context tree. It is one piece of a context tree
type ContextTreeSlice struct {
	Name     string             `json:"name"`
	Ref      *string            `json:"ref"` // Ref in this case is only one value because it's used to search
	Children []ContextTreeSlice `json:"children"`
}

// Retrieve context from a tree slice by name
func (c ContextTreeSlice) GetContextByName(name string) (*ContextTreeSlice, bool) {
	if c.Name == name {
		return &c, true
	}

	if len(c.Children) == 0 {
		return nil, false
	}

	// Recursive style checking all the children
	for _, cc := range c.Children {
		nc, exists := cc.GetContextByName(name)
		if exists {
			return nc, exists
		}
	}

	return nil, false
}

// Retrieve a context from a tree slice by Ref
func (c ContextTreeSlice) GetContextByRef(ref string) (*ContextTreeSlice, bool) {
	if c.Ref != nil {
		if *c.Ref == ref {
			return &c, true
		}
	}

	if len(c.Children) == 0 {
		return nil, false
	}

	// Recursive style checking all the children
	for _, cc := range c.Children {
		nc, exists := cc.GetContextByRef(ref)
		if exists {
			return nc, exists
		}
	}

	return nil, false
}
