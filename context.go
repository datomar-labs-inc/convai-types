package ctypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/osteele/liquid"

	"github.com/datomar-labs-inc/convai-types/deepcopy"
)

var reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
var liquidEngine *liquid.Engine

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
	Name  string            `json:"name"`
	Ref   *string           `json:"ref"` // Ref in this case is only one value because it's used to search
	Child *ContextTreeSlice `json:"child"`
}

// Retrieve context from a tree slice by name
func (c ContextTreeSlice) GetContextByName(name string) (*ContextTreeSlice, bool) {
	if c.Name == name {
		return &c, true
	}

	if c.Child == nil {
		return nil, false
	}

	return c.Child.GetContextByName(name)
}

// Retrieve a context from a tree slice by Ref
func (c ContextTreeSlice) GetContextByRef(ref string) (*ContextTreeSlice, bool) {
	if c.Ref != nil {
		if *c.Ref == ref {
			return &c, true
		}
	}

	if c.Child == nil {
		return nil, false
	}

	return c.Child.GetContextByRef(ref)
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

func (c *Context) GetData(path string) (interface{}, bool) {
	// Validate the path
	if !ValidateDataPath(path) {
		return nil, false
	}

	// Get the context
	ctx, exists := c.GetContextByName(GetDataPathContextLevelName(path))
	if !exists {
		return nil, false
	}

	// Return the data
	for _, mc := range ctx.Memory {
		if mc.Name == GetDataPathMemoryContainerName(path) {
			if DataPathHasMultipartKey(path) {
				var parts []reflect.Value

				for _, kp := range GetDataPathKeyParts(path) {
					parts = append(parts, reflect.ValueOf(kp))
				}

				val, err := index(reflect.ValueOf(mc.Data), parts...)
				if err != nil {
					return nil, false
				}

				return val.Interface(), true
			} else {
				data, ok := mc.Data[GetDataPathKey(path)]
				return data, ok
			}
		}
	}

	return nil, false
}

func (c *Context) GetDataString(path string) (string, bool) {
	d, ok := c.GetData(path)
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%v", d), true
}

func (c *Context) GetDataInt(path string) (int, bool) {
	d, ok := c.GetData(path)
	if !ok {
		return 0, false
	}

	n, err := strconv.Atoi(fmt.Sprintf("%v", d))
	if err != nil {
		return 0, false
	}

	return n, true
}

func (c *Context) GetDataFloat(path string) (float64, bool) {
	d, ok := c.GetData(path)
	if !ok {
		return 0, false
	}

	n, err := strconv.ParseFloat(fmt.Sprintf("%v", d), 64)
	if err != nil {
		return 0, false
	}

	return n, true
}

func (c *Context) GetTemplateData() map[string]interface{} {
	data := map[string]interface{}{}

	currentContext := c

	// Add all tree items to data
	for {
		tlData := map[string]interface{}{}

		for _, mc := range currentContext.Memory {
			tlData[mc.Name] = mustMappify(mc.Data)
		}

		data[currentContext.Name] = tlData

		if len(currentContext.Children) == 0 {
			break
		} else {
			currentContext = &currentContext.Children[0]
		}
	}

	return data
}

func (c *Context) ExecuteTemplate(tmpl string) ([]byte, error) {
	if liquidEngine == nil {
		liquidEngine = liquid.NewEngine()
	}

	return liquidEngine.ParseAndRender([]byte(tmpl), c.GetTemplateData())
}

func (c *Context) ExecuteTemplateString(tmpl string) (string, error) {
	if liquidEngine == nil {
		liquidEngine = liquid.NewEngine()
	}

	return liquidEngine.ParseAndRenderString(tmpl, c.GetTemplateData())
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

func mustMappify(in interface{}) map[string]interface{} {
	jsb, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	var out map[string]interface{}

	err = json.Unmarshal(jsb, &out)
	if err != nil {
		panic(err)
	}

	return out
}

// index returns the result of indexing its first argument by the following
// arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each
// indexed item must be a map, slice, or array.
func index(item reflect.Value, indexes ...reflect.Value) (reflect.Value, error) {
	item = indirectInterface(item)
	if !item.IsValid() {
		return reflect.Value{}, fmt.Errorf("index of untyped nil")
	}
	for _, index := range indexes {
		index = indirectInterface(index)
		var isNil bool
		if item, isNil = indirect(item); isNil {
			return reflect.Value{}, fmt.Errorf("index of nil pointer")
		}
		switch item.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			x, err := indexArg(index, item.Len())
			if err != nil {
				return reflect.Value{}, err
			}
			item = item.Index(x)
		case reflect.Map:
			index, err := prepareArg(index, item.Type().Key())
			if err != nil {
				return reflect.Value{}, err
			}
			if x := item.MapIndex(index); x.IsValid() {
				item = x
			} else {
				item = reflect.Zero(item.Type().Elem())
			}
		case reflect.Invalid:
			// the loop holds invariant: item.IsValid()
			panic("unreachable")
		default:
			return reflect.Value{}, fmt.Errorf("can't index item of type %s", item.Type())
		}
	}
	return item, nil
}

// indirectInterface returns the concrete value in an interface value,
// or else the zero reflect.Value.
// That is, if v represents the interface value x, the result is the same as reflect.ValueOf(x):
// the fact that x was an interface value is forgotten.
func indirectInterface(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Interface {
		return v
	}
	if v.IsNil() {
		return reflect.Value{}
	}
	return v.Elem()
}

// indirect returns the item at the end of indirection, and a bool to indicate
// if it's nil. If the returned bool is true, the returned value's kind will be
// either a pointer or interface.
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}

// indexArg checks if a reflect.Value can be used as an index, and converts it to int if possible.
func indexArg(index reflect.Value, cap int) (int, error) {
	var x int64
	switch index.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = index.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x = int64(index.Uint())
	case reflect.Invalid:
		return 0, fmt.Errorf("cannot index slice/array with nil")
	default:
		return 0, fmt.Errorf("cannot index slice/array with type %s", index.Type())
	}
	if x < 0 || int(x) < 0 || int(x) > cap {
		return 0, fmt.Errorf("index out of range: %d", x)
	}
	return int(x), nil
}

// prepareArg checks if value can be used as an argument of type argType, and
// converts an invalid value to appropriate zero if possible.
func prepareArg(value reflect.Value, argType reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if !canBeNil(argType) {
			return reflect.Value{}, fmt.Errorf("value is nil; should be of type %s", argType)
		}
		value = reflect.Zero(argType)
	}
	if value.Type().AssignableTo(argType) {
		return value, nil
	}
	if intLike(value.Kind()) && intLike(argType.Kind()) && value.Type().ConvertibleTo(argType) {
		value = value.Convert(argType)
		return value, nil
	}
	return reflect.Value{}, fmt.Errorf("value has type %s; should be %s", value.Type(), argType)
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	case reflect.Struct:
		return typ == reflectValueType
	}
	return false
}

func intLike(typ reflect.Kind) bool {
	switch typ {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	}
	return false
}
