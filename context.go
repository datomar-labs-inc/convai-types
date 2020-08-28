package ctypes

import (
	"errors"

	"github.com/google/uuid"

	"github.com/datomar-labs-inc/convai-types/deepcopy"
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

// Retrieve context from a tree slice by name
func (c *Context) GetContextByName(name string) (*Context, bool) {
	if c.Name == name {
		return c, true
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
func (c *Context) GetContextByRef(ref string) (*Context, bool) {
	if c.Ref != nil {
		for _, r := range c.Ref {
			if r == ref {
				return c, true
			}
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

// GetLastTreeItem returns the deepest child context in the tree.
// Note, this only works when each context has 0 or 1 children
func (c *Context) GetLastTreeItem() *Context {
	if len(c.Children) > 0 {
		return c.GetLastTreeItem()
	}

	return c
}

// AddChildContext adds a child to the current context, and returns the current context
func (c *Context) AddChildContext(context *Context) *Context {
	context.ParentID = c.ID
	context.Parent = c
	c.Children = append(c.Children, *context)
	return c
}

// WithTransformations returns a new context tree with transformations applied
func (c *Context) WithTransformations(transformations []Transformation) (*Context, error) {
	newMemory := []MemoryContainer{}

	for _, mc := range c.Memory {
		memCopy, err := deepcopy.DeepCopy(mc.Data)
		if err != nil {
			panic(err)
		}

		newMem := MemoryContainer{
			Name:    mc.Name,
			Type:    mc.Type,
			Exposed: mc.Exposed,
			Data:    memCopy,
		}

		for _, transformation := range transformations {
			if transformation.PathValid() && transformation.GetContextLevelName() == c.Name && transformation.GetMemoryContainerName() == mc.Name {
				newMem.Transform(transformation)
			} else if !transformation.PathValid() {
				return nil, errors.New("invalid transformation path: " + transformation.Path)
			}
		}

		newMemory = append(newMemory, newMem)
	}

	newContext := Context{
		Name:     c.Name,
		ID:       c.ID,
		ParentID: c.ParentID,
		Ref:      c.Ref,
		Memory:   newMemory,
	}

	if len(c.Children) > 0 {
		var children []Context

		for _, child := range c.Children {
			transformed, err := child.WithTransformations(transformations)
			if err != nil {
				return nil, err
			}

			transformed.ParentID = newContext.ID
			transformed.Parent = &newContext

			children = append(children, *transformed)
		}

		newContext.Children = children
	}

	return &newContext, nil
}