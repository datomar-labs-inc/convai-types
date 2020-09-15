package ctypes

import (
	b64 "encoding/base64"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type DBPackage struct {
	ID             uuid.UUID   `db:"id" json:"id" validate:"required"`
	Name           string      `db:"name" json:"name" validate:"required"`
	Description    string      `db:"description" json:"description"`
	OrganizationID uuid.UUID   `db:"-" json:"organization_id"`
	BaseURL        string      `db:"base_url" json:"base_url" validate:"required,url"`
	SigningKey     string      `db:"signing_key" json:"signing_key,omitempty"`
	CreatedAt      *CustomTime `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      *CustomTime `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Package struct {
	DBPackage

	Nodes      []DBNode     `json:"nodes"`
	Links      []DBLink     `json:"links"`
	Events     []DBEvent    `json:"events"`
	Dispatches []DBDispatch `json:"dispatches"`
}

type PackageDifferences struct {
	UpdatedNodes      []DBNode     `json:"updated_nodes"`
	NewNodes          []DBNode     `json:"new_nodes"`
	UpdatedLinks      []DBLink     `json:"updated_links"`
	NewLinks          []DBLink     `json:"new_links"`
	UpdatedEvents     []DBEvent    `json:"updated_events"`
	NewEvents         []DBEvent    `json:"new_events"`
	UpdatedDispatches []DBDispatch `json:"updated_dispatches"`
	NewDispatches     []DBDispatch `json:"new_dispatches"`
}

func ComputePackageDifferences(old, new *Package) *PackageDifferences {
	diff := PackageDifferences{}

	for _, newNode := range new.Nodes {
		alreadyExisted := false

		newNode.PackageID = new.ID

		for _, oldNode := range old.Nodes {
			if oldNode.TypeID == newNode.TypeID && oldNode.PackageID == newNode.PackageID && oldNode.Version == newNode.Version {
				alreadyExisted = true

				// This node already existed. It only needs to be changed if it was updated
				if oldNode.Name != newNode.Name || oldNode.Documentation != newNode.Documentation || !reflect.DeepEqual(oldNode.Style, newNode.Style) {
					diff.UpdatedNodes = append(diff.UpdatedNodes, newNode)
				}

				break
			}
		}

		if !alreadyExisted {
			diff.NewNodes = append(diff.NewNodes, newNode)
		}
	}

	for _, newLink := range new.Links {
		alreadyExisted := false

		newLink.PackageID = new.ID

		for _, oldLink := range old.Links {
			if oldLink.TypeID == newLink.TypeID && oldLink.PackageID == newLink.PackageID && oldLink.Version == newLink.Version {
				alreadyExisted = true

				// This link already existed. It only needs to be changed if it was updated
				if oldLink.Name != newLink.Name || oldLink.Documentation != newLink.Documentation || !reflect.DeepEqual(oldLink.Style, newLink.Style) {
					diff.UpdatedLinks = append(diff.UpdatedLinks, newLink)
				}

				break
			}
		}

		if !alreadyExisted {
			diff.NewLinks = append(diff.NewLinks, newLink)
		}
	}

	for _, newEvent := range new.Events {
		alreadyExisted := false

		newEvent.PackageID = new.ID

		for _, oldEvent := range old.Events {
			if oldEvent.ID == newEvent.ID && oldEvent.PackageID == newEvent.PackageID {
				alreadyExisted = true

				// This event already existed. It only needs to be changed if it was updated
				if oldEvent.Name != newEvent.Name || oldEvent.Documentation != newEvent.Documentation || !reflect.DeepEqual(oldEvent.Style, newEvent.Style) {
					diff.UpdatedEvents = append(diff.UpdatedEvents, newEvent)
				}

				break
			}
		}

		if !alreadyExisted {
			diff.NewEvents = append(diff.NewEvents, newEvent)
		}
	}

	for _, newDispatch := range new.Dispatches {
		alreadyExisted := false

		newDispatch.PackageID = new.ID

		for _, oldDispatch := range old.Dispatches {
			if oldDispatch.ID == newDispatch.ID && oldDispatch.PackageID == newDispatch.PackageID {
				alreadyExisted = true

				// This dispatch already existed. It only needs to be changed if it was updated
				if oldDispatch.Name != newDispatch.Name || oldDispatch.Documentation != newDispatch.Documentation {
					diff.UpdatedDispatches = append(diff.UpdatedDispatches, newDispatch)
				}

				break
			}
		}

		if !alreadyExisted {
			diff.NewDispatches = append(diff.NewDispatches, newDispatch)
		}
	}

	return &diff
}

func (p *DBPackage) Cursor() string {
	s := fmt.Sprintf("%d,%s", p.CreatedAt.UnixNano(), p.ID)
	return b64.StdEncoding.EncodeToString([]byte(s))
}
