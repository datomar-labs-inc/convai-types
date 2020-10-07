package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blang/semver"
	"github.com/google/uuid"
	"upper.io/db.v3/postgresql"
)

var INITIAL_VERSION = Semver{semver.MustParse("0.0.1")}

type DBModule struct {
	ID             uuid.UUID   `db:"id" json:"id"`
	Version        string      `db:"version" json:"version"`
	Changelog      string      `db:"changelog" json:"changelog"`
	Name           string      `db:"name" json:"name"`
	Graph          GraphModule `db:"graph" json:"graph"`
	OrganizationID uuid.UUID   `db:"organization_id" json:"organization_id"`
	CreatedAt      *CustomTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt      *CustomTime `db:"updated_at,omitempty" json:"updated_at"`
}

type GraphNode struct {
	ID uuid.UUID `json:"id" msgpack:"i"`

	// Properties that are common to all types of nodes
	PackageID uuid.UUID `json:"package_id,omitempty" msgpack:"p,omitempty"`
	Label     string    `json:"label" msgpack:"la"`

	// Properties of a node that references functionality in a package
	TypeID     *string `json:"type_id,omitempty" msgpack:"n,omitempty"`
	Version    *string `json:"version,omitempty" msgpack:"v,omitempty"`
	ConfigJSON *string `json:"config_json,omitempty" msgpack:"c,omitempty"`

	// Properties of a node that acts as a reference to a module
	ModuleID      *uuid.UUID `json:"module_id,omitempty" msgpack:"m,omitempty"`
	ModuleVersion *string    `json:"module_version,omitempty" msgpack:"mv,omitempty"`

	// Properties of a node that acts as an event entrypoint
	EventTypeID *string `json:"event_type_id,omitempty" msgapck:"e,omitempty"`

	Layout Point `json:"layout"`
}

// CompiledGraphLink is the variant of link that lives in a compiled executable
type GraphLink struct {
	ID         uuid.UUID `json:"id" msgpack:"i"`
	Label      string    `json:"label" msgpack:"la"`
	PackageID  uuid.UUID `json:"package_id" msgpack:"p"`
	TypeID     string    `json:"type_id" msgpack:"l"`
	Version    string    `json:"version" msgpack:"v"`
	Priority   int       `json:"priority"`
	ConfigJSON string    `json:"config_json" msgpack:"c"`
	A          LinkPoint `json:"a"`
	B          LinkPoint `json:"b"`
}

type LinkPoint struct {
	NodeID   *uuid.UUID `json:"node_id"`
	IsOutput bool       `json:"is_input"`
	Position Point      `json:"pos"`
}

func (l *GraphLink) Validate() error {
	if l.ID == uuid.Nil {
		return errors.New("link cannot have nil id")
	}

	if l.TypeID == "" {
		return errors.New("link type id missing")
	}

	if l.A.NodeID != nil && l.B.NodeID != nil && *l.A.NodeID == *l.B.NodeID {
		return errors.New("link source and destination cannot be identical")
	}

	if l.A.Position.X == l.B.Position.X && l.A.Position.Y == l.B.Position.Y {
		return errors.New("both link points cannot occupy the same position")
	}

	if !json.Valid([]byte(l.ConfigJSON)) {
		return errors.New("config json is invalid")
	}

	if _, err := semver.Parse(l.Version); err != nil {
		return fmt.Errorf("invalid version: %v", err)
	}

	return nil
}

// Single point used to derive a layout/draw an object
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type CompiledGraphModule struct {
	Nodes map[uuid.UUID]CompiledGraphNode `json:"nodes" msgpack:"n"`
	Links []CompiledGraphLink             `json:"links" msgpack:"l"`
}

type GraphModule struct {
	ID    uuid.UUID               `json:"id" msgpack:"i"`
	Label string                  `json:"label" msgpack:"la"`
	Nodes map[uuid.UUID]GraphNode `json:"nodes" msgpack:"n"`
	Links []GraphLink             `json:"links" msgpack:"l"`
}

func (m *GraphModule) Validate() error {
	if m.ID == uuid.Nil {
		return errors.New("module cannot have nil id")
	}

	if m.Nodes == nil {
		return errors.New("module cannot have nil nodes")
	}

	if m.Links == nil {
		return errors.New("module cannot have nil links")
	}

	// TODO validate links do not have double inputs

	// Check for duplicate links
	dups := map[string]bool{}

	for _, l := range m.Links {
		if l.B.NodeID != nil && l.A.NodeID != nil {
			key := fmt.Sprintf("%s-%s", l.A.NodeID.String(), l.B.NodeID.String())

			if dups[key] {
				return fmt.Errorf("duplicate links %s", key)
			} else {
				dups[key] = true
			}
		}
	}

	return nil
}

func (m *GraphModule) GetLink(id uuid.UUID) (*GraphLink, int) {
	for i, l := range m.Links {
		if l.ID == id {
			return &l, i
		}
	}

	return nil, 0
}

func (m *GraphModule) DeleteLink(id uuid.UUID) {
	idx := -1

	for i, l := range m.Links {
		if l.ID == id {
			idx = i
			break
		}
	}

	if idx != -1 {
		m.Links = append(m.Links[:idx], m.Links[idx+1:]...)
	}
}

func (g GraphModule) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *GraphModule) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &GraphModule{}
	_ sql.Scanner   = &GraphModule{}
)
