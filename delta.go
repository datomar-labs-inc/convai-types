package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

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

type ApplyDeltaRequest struct {
	Operations  DeltaOperations `json:"operations"`
	BlueprintID *uuid.UUID      `json:"blueprint_id,omitempty"`
}

type DBDelta struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	AccountID   uuid.UUID       `db:"account_id" json:"account_id"`
	UpdateType  int             `db:"update_type" json:"update_type"`
	Operations  DeltaOperations `db:"delta" json:"delta"`
	BlueprintID *uuid.UUID      `db:"blueprint_id,omitempty" json:"blueprint_id,omitempty"`
	CreatedAt   *CustomTime     `db:"created_at,omitempty" json:"created_at,omitempty"`
}

func (d *DBDelta) GetGroovePath() string {
	if d.BlueprintID != nil {
		return fmt.Sprintf("%s.%s", d.BlueprintID.String(), d.ID.String())
	} else {
		idParts := ""

		for _, op := range d.Operations {
			if op.Type == DOUpdateEnvironmentPackageConfig {
				idParts += op.UpdateEnvironmentPackageConfig.EnvironmentID.String()
			}
		}

		if idParts == "" {
			idParts = "unknown"
		}

		return fmt.Sprintf("%s.%s", idParts, d.ID.String())
	}
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
	ModuleID *uuid.UUID `json:"module_id,omitempty"`
	Type     uint       `json:"type"`

	MoveNode                       *DeltaMoveNode                       `json:"move_node,omitempty"`
	MoveLink                       *DeltaMoveLink                       `json:"move_link,omitempty"`
	CreateNode                     *DeltaCreateNode                     `json:"create_node,omitempty"`
	CreateLink                     *DeltaCreateLink                     `json:"create_link,omitempty"`
	DeleteNode                     *DeltaDeleteNode                     `json:"delete_node,omitempty"`
	DeleteLink                     *DeltaDeleteLink                     `json:"delete_link,omitempty"`
	UpdateNode                     *DeltaUpdateNode                     `json:"update_node,omitempty"`
	UpdateLink                     *DeltaUpdateLink                     `json:"update_link,omitempty"`
	UpdateNodePackageConfig        *DeltaUpdateNodePackageConfig        `json:"update_node_package_config,omitempty"`
	UpdateLinkPackageConfig        *DeltaUpdateLinkPackageConfig        `json:"update_link_package_config,omitempty"`
	UpdateEnvironmentPackageConfig *DeltaUpdateEnvironmentPackageConfig `json:"update_environment_package_config,omitempty"`

	CreateModule *DeltaCreateModule `json:"create_module,omitempty"`
	UpdateModule *DeltaUpdateModule `json:"update_module,omitempty"`
}

// Makes a best effort attempt to validate the delta with all immediately available information
// returns the first encountered error and stops searching
func (d DeltaOperations) Validate() error {
	for _, op := range d {
		err := op.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DeltaOperation) Validate() error {
	// First validate module_id present-ness
	// There are duplicated switch statements in this function because it actually ends up saving space 
	// vs duplicating this check for all switch cases individually
	switch d.Type {
	case DOMoveLink, DOCreateNode, DOCreateLink,
		DODeleteNode, DODeleteLink, DOUpdateLink,
		DOUpdateNode, DOUpdateLinkPackageConfig,
		DOUpdateNodePackageConfig, DOMoveNode, DOUpdateModule:
		if d.ModuleID == nil || *d.ModuleID == uuid.Nil {
			return fmt.Errorf("invalid module_id %v", d.ModuleID)
		}
	}

	switch d.Type {
	case DOMoveNode:
		if d.MoveNode == nil {
			return errors.New("move node cannot be null")
		}

		return d.MoveNode.Validate()
	case DOMoveLink:
		if d.MoveLink == nil {
			return errors.New("move link cannot be null")
		}

		return d.MoveLink.Validate()
	case DOCreateNode:
		if d.CreateNode == nil {
			return errors.New("create node cannot be null")
		}

		return d.CreateNode.Validate()
	case DOCreateLink:
		if d.CreateLink == nil {
			return errors.New("create link cannot be null")
		}

		return d.CreateLink.Validate()
	case DODeleteNode:
		if d.DeleteNode == nil {
			return errors.New("DeleteNode cannot be null")
		}

		return d.DeleteNode.Validate()
	case DODeleteLink:
		if d.DeleteLink == nil {
			return errors.New("DeleteLink cannot be null")
		}

		return d.DeleteLink.Validate()
	case DOUpdateLink:
		if d.UpdateLink == nil {
			return errors.New("UpdateLink cannot be null")
		}

		return d.UpdateLink.Validate()
	case DOUpdateNode:
		if d.UpdateNode == nil {
			return errors.New("UpdateNode cannot be null")
		}

		return d.UpdateNode.Validate()
	case DOUpdateLinkPackageConfig:
		if d.UpdateLinkPackageConfig == nil {
			return errors.New("UpdateLinkPackageConfig cannot be null")
		}

		return d.UpdateLinkPackageConfig.Validate()
	case DOUpdateNodePackageConfig:
		if d.UpdateNodePackageConfig == nil {
			return errors.New("UpdateNodePackageConfig cannot be null")
		}

		return d.UpdateNodePackageConfig.Validate()
	case DOCreateModule:
		if d.CreateModule == nil {
			return errors.New("CreateModule cannot be null")
		}

		return d.CreateModule.Validate()
	case DODeleteModule:
		// if d.DeleteModule == nil {
		// 	return errors.New("DeleteModule cannot be null")
		// }
		//
		// return d.DeleteModule.Validate()
	case DOUpdateModule:
		if d.UpdateModule == nil {
			return errors.New("UpdateModule cannot be null")
		}

		return d.UpdateModule.Validate()
	case DOCreateEnvironment:
		// if d.CreateEnvironment == nil {
		// 	return errors.New("CreateEnvironment cannot be null")
		// }
		//
		// return d.CreateEnvironment.Validate()
	case DOUpdateEnvironment:
		// if d.UpdateEnvironment == nil {
		// 	return errors.New("UpdateEnvironment cannot be null")
		// }
		//
		// return d.UpdateEnvironment.Validate()
	case DODeleteEnvironment:
		// if d.DeleteEnvironment == nil {
		// 	return errors.New("DeleteEnvironment cannot be null")
		// }
		//
		// return d.DeleteEnvironment.Validate()
	case DOUpdateEnvironmentPackageConfig:
		if d.UpdateEnvironmentPackageConfig == nil {
			return errors.New("UpdateEnvironmentPackageConfig cannot be null")
		}

		return d.UpdateEnvironmentPackageConfig.Validate()
	case DOUpdateBot:
		// if d.UpdateBot == nil {
		// 	return errors.New("UpdateBot cannot be null")
		// }
		//
		// return d.UpdateBot.Validate()
	default:
		return fmt.Errorf("unknown operation %d", d.Type)
	}

	return nil
}

type DeltaMoveNode struct {
	ID  uuid.UUID `json:"node_id"`
	Pos Point     `json:"pos"`
}

func (d *DeltaMoveNode) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate node point

	return nil
}

type DeltaMoveLink struct {
	ID uuid.UUID `json:"link_id"`
	A  Point     `json:"a"`
	B  Point     `json:"b"`
}

func (d *DeltaMoveLink) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate link points

	return nil
}

type DeltaCreateNode GraphNode

func (d *DeltaCreateNode) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate create node

	return nil
}

type DeltaCreateLink GraphLink

func (d *DeltaCreateLink) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate link creation

	return nil
}

type DeltaDeleteNode struct {
	ID uuid.UUID `json:"id"`
}

func (d *DeltaDeleteNode) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate node deletion

	return nil
}

type DeltaDeleteLink struct {
	ID uuid.UUID `json:"id"`
}

func (d *DeltaDeleteLink) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate link deletion

	return nil
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

func (d *DeltaUpdateNode) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate node update

	return nil
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

func (d *DeltaUpdateLink) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate link update

	return nil
}

type DeltaUpdateNodePackageConfig struct {
	ID     uuid.UUID `json:"id"`
	Config string    `json:"config"`
}

func (d *DeltaUpdateNodePackageConfig) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate

	return nil
}

type DeltaUpdateLinkPackageConfig struct {
	ID     uuid.UUID `json:"id"`
	Config string    `json:"config"`
}

func (d *DeltaUpdateLinkPackageConfig) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate

	return nil
}

type DeltaCreateModule struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (d *DeltaCreateModule) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("invalid id")
	}

	// TODO validate

	return nil
}

type DeltaUpdateModule struct {
	Name string `json:"name"`
}

func (d *DeltaUpdateModule) Validate() error {
	// TODO validate

	return nil
}

type DeltaUpdateEnvironmentPackageConfig struct {
	EnvironmentID uuid.UUID   `json:"environment_id"`
	PackageID     uuid.UUID   `json:"package_id"`
	Data          interface{} `json:"data"`
}

func (d *DeltaUpdateEnvironmentPackageConfig) Validate() error {
	// TODO validate

	return nil
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
