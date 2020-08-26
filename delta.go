package ctypes

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

// TODO make db types
type DBDelta struct {
	ID          uuid.UUID `db:"id" json:"id"`
	AccountID   uuid.UUID
	Operations  []DeltaOperation
	BlueprintID uuid.UUID
	CreatedAt   *time.Time
}

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
		link := module.GetLink(operation.MoveLink.ID)

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
		// TODO implement method
	case DOUpdateLink:
		// TODO implement method

	case DOUpdateNodePackageConfig:
		if node, ok := module.Nodes[operation.UpdateNodePackageConfig.ID]; ok {
			node.ConfigJSON = &operation.UpdateNodePackageConfig.Config
			module.Nodes[operation.UpdateNodePackageConfig.ID] = node
		} else {
			return errors.New("could not update node package config because node did not exist")
		}

	case DOUpdateLinkPackageConfig:
		link := module.GetLink(operation.UpdateLinkPackageConfig.ID)

		if link != nil {
			link.ConfigJSON = operation.UpdateLinkPackageConfig.Config
		} else {
			return errors.New("could not update link package config because link did not exist")
		}
	}

	return nil
}
