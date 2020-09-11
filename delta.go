package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"time"

	"upper.io/db.v3/postgresql"

	"github.com/google/uuid"
)

const (
	UpdateTypeStandard = iota
	UpdateTypeUndo
	UpdateTypeRollback
	UpdateTypeRedo
)

const (
	DOMoveNode = iota
	DOMoveLink
	DOCreateNode
	DOCreateLink
	DODeleteNode
	DODeleteLink
	DOUpdateLink
	DOUpdateNode
	DOUpdateLinkPackageConfig
	DOUpdateNodePackageConfig
	DOCreateModule
	DODeleteModule
	DOUpdateModule
	DOCreateEnvironment
	DOUpdateEnvironment
	DODeleteEnvironment
	DOUpdateEnvironmentPackageConfig
	DOUpdateBot
)

type DBDelta struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	AccountID   uuid.UUID       `db:"account_id" json:"account_id"`
	UpdateType  int             `db:"update_type" json:"update_type"`
	Operations  DeltaOperations `db:"delta" json:"delta"`
	BlueprintID uuid.UUID       `db:"blueprint_id" json:"blueprint_id"`
	CreatedAt   *time.Time      `db:"created_at,omitempty" json:"created_at,omitempty"`
}

type DeltaOperations []DeltaOperation

func (g DeltaOperations) Value() (driver.Value, error) {
	return postgresql.EncodeJSONB(g)
}

func (g *DeltaOperations) Scan(src interface{}) error {
	return postgresql.DecodeJSONB(g, src)
}

var (
	_ driver.Valuer = &DeltaOperations{}
	_ sql.Scanner   = &DeltaOperations{}
)

type DeltaOperation struct {
	ModuleID uuid.UUID `json:"module_id"`
	Type     int       `json:"type"`

	MoveNode                *DeltaMoveNode                `json:"move_node,omitempty"`
	MoveLink                *DeltaMoveLink                `json:"move_link,omitempty"`
	CreateNode              *DeltaCreateNode              `json:"create_node,omitempty"`
	CreateLink              *DeltaCreateLink              `json:"create_link,omitempty"`
	DeleteNode              *DeltaDeleteNode              `json:"delete_node,omitempty"`
	DeleteLink              *DeltaDeleteLink              `json:"delete_link,omitempty"`
	UpdateNode              *DeltaUpdateNode              `json:"update_node,omitempty"`
	UpdateLink              *DeltaUpdateLink              `json:"update_link,omitempty"`
	UpdateNodePackageConfig *DeltaUpdateNodePackageConfig `json:"update_node_package_config,omitempty"`
	UpdateLinkPackageConfig *DeltaUpdateLinkPackageConfig `json:"update_link_package_config,omitempty"`

	CreateModule *DeltaCreateModule `json:"create_module,omitempty"`
	UpdateModule *DeltaUpdateModule `json:"update_module,omitempty"`
}

type DeltaMoveNode struct {
	ID  uuid.UUID `json:"node_id"`
	Pos Point     `json:"pos"`
}

type DeltaMoveLink struct {
	ID uuid.UUID `json:"link_id"`
	A  Point     `json:"a"`
	B  Point     `json:"b"`
}

type DeltaCreateNode GraphNode

type DeltaCreateLink GraphLink

type DeltaDeleteNode struct {
	ID uuid.UUID `json:"id"`
}

type DeltaDeleteLink struct {
	ID uuid.UUID `json:"id"`
}

type DeltaUpdateNode struct {
	ID    uuid.UUID `json:"id"`
	Label *string   `json:"label,omitempty"`

	// Properties that are common to all types of nodes
	PackageID *uuid.UUID `json:"package_id,omitempty" msgpack:"p,omitempty"`

	// Properties of a node that references functionality in a package
	TypeID  *string `json:"type_id,omitempty" msgpack:"n,omitempty"`
	Version *string `json:"version,omitempty" msgpack:"v,omitempty"`

	// Properties of a node that acts as a reference to a module
	ModuleID      *uuid.UUID `json:"module_id,omitempty" msgpack:"m,omitempty"`
	ModuleVersion *string    `json:"module_version,omitempty" msgpack:"mv,omitempty"`

	// Properties of a node that acts as an event entrypoint
	EventTypeID *string `json:"event_type_id,omitempty" msgapck:"e,omitempty"`
}

// TODO figure out what types this should have
type DeltaUpdateLink struct {
	ID        uuid.UUID  `json:"id" msgpack:"i"`
	Label     *string    `json:"label" msgpack:"b"`
	PackageID *uuid.UUID `json:"package_id,omitempty" msgpack:"p"`
	TypeID    *string    `json:"type_id,omitempty" msgpack:"l"`
	Version   *string    `json:"version,omitempty" msgpack:"v"`
	Priority  *int       `json:"priority,omitempty"`
	A         *LinkPoint `json:"a,omitempty"`
	B         *LinkPoint `json:"b,omitempty"`
}

type DeltaUpdateNodePackageConfig struct {
	ID     uuid.UUID `json:"id"`
	Config string    `json:"config"`
}

type DeltaUpdateLinkPackageConfig struct {
	ID     uuid.UUID `json:"id"`
	Config string    `json:"config"`
}

type DeltaCreateModule struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type DeltaUpdateModule struct {
	Name string `json:"name"`
}

func ApplyDeltaToModule(module *GraphModule, delta *DBDelta) error {
	for _, op := range delta.Operations {
		err := ApplyOperationToModule(module, &op)
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyOperationToModule(module *GraphModule, operation *DeltaOperation) error {
	if module.Nodes == nil {
		module.Nodes = map[uuid.UUID]GraphNode{}
	}

	switch operation.Type {
	case DOMoveNode:
		if node, ok := module.Nodes[operation.MoveNode.ID]; ok {
			node.Layout = operation.MoveNode.Pos
			module.Nodes[operation.MoveNode.ID] = node
		} else {
			return errors.New("could not update node position because node did not exist")
		}

	case DOMoveLink:
		link, _ := module.GetLink(operation.MoveLink.ID)

		if link != nil {
			link.A.Position = operation.MoveLink.A
			link.B.Position = operation.MoveLink.B
		} else {
			return errors.New("could not update link position because link did not exist")
		}

	case DOCreateNode:
		module.Nodes[operation.CreateNode.ID] = GraphNode(*operation.CreateNode)

	case DOCreateLink:
		module.Links = append(module.Links, GraphLink(*operation.CreateLink))

	case DODeleteNode:
		delete(module.Nodes, operation.DeleteNode.ID)

	case DODeleteLink:
		module.DeleteLink(operation.DeleteLink.ID)

	case DOUpdateNode:
		if node, ok := module.Nodes[operation.UpdateNode.ID]; ok {
			// Use reflection to copy over non-nil fields from the update to the link where the names match
			update := reflect.ValueOf(operation.UpdateNode).Elem()
			nde := reflect.ValueOf(&node).Elem()

			// Loop over each field in the update and copy properties over
			for i := 1; i < update.NumField(); i++ {
				updateField := update.Field(i)

				if !updateField.IsNil() {
					nde.FieldByName(update.Type().Field(i).Name).Set(update.Field(i).Elem())
				}
			}

			module.Nodes[operation.UpdateNode.ID] = node
		} else {
			return errors.New("could not update node because node did not exist")
		}

	case DOUpdateLink:
		link, i := module.GetLink(operation.UpdateLink.ID)

		if link != nil {

			// Use reflection to copy over non-nil fields from the update to the link where the names match
			update := reflect.ValueOf(operation.UpdateLink).Elem()
			lnk := reflect.ValueOf(link).Elem()

			// Loop over each field in the update and copy properties over
			for i := 1; i < update.NumField(); i++ {
				updateField := update.Field(i)

				if !updateField.IsNil() {
					lnk.FieldByName(update.Type().Field(i).Name).Set(update.Field(i).Elem())
				}
			}

			module.Links[i] = *link
		} else {
			return errors.New("could not update link because link did not exist")
		}

	case DOUpdateNodePackageConfig:
		if node, ok := module.Nodes[operation.UpdateNodePackageConfig.ID]; ok {
			node.ConfigJSON = &operation.UpdateNodePackageConfig.Config
			module.Nodes[operation.UpdateNodePackageConfig.ID] = node
		} else {
			return errors.New("could not update node package config because node did not exist")
		}

	case DOUpdateLinkPackageConfig:
		link, _ := module.GetLink(operation.UpdateLinkPackageConfig.ID)

		if link != nil {
			link.ConfigJSON = operation.UpdateLinkPackageConfig.Config
		} else {
			return errors.New("could not update link package config because link did not exist")
		}
	}

	return nil
}
