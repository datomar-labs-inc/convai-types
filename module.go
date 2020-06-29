package ctypes

import (
	"github.com/google/uuid"
)

// Data fields for an asset originating from a package
type PackageAsset struct {
	PackageAssetID      uuid.UUID // ID of the asset in the package
	PackageAssetVersion string    // Version of the asset
	PackageID           string    // ID of the package this asset is from
}

type Module struct {
	ID      uuid.UUID
	Name    string
	Graph   Graph
	PackageAsset
}

type Graph struct {
	Objects      []GraphObject
	LastObjectID int // Couldnt find a better way to figure this out
}

// How do we handle styling at this point? Part of config or its own data field?
// Object on a graph (link, node)
type GraphObject struct {
	ID          int         // ID of the object on the graph/in the module
	Name        string      // Name of the object (having this outside the config allows nil config objects)
	Type        int         // Node, Link?
	ConfigID    *uuid.UUID  // DB ID of the config
	Layout      []Point     // Where the object is located on the graph
	Connections []uuid.UUID // A list of objects this is connected to on the graph
	PackageAsset
}

// Single point used to derive a layout/draw an object
type Point struct {
	X int
	Y int
}

type GraphModule struct {
	ID    uuid.UUID               `json:"id" msgpack:"i"`
	Nodes map[uuid.UUID]GraphNode `json:"nodes" msgpack:"n"`
	Links map[uuid.UUID]GraphLink `json:"links" msgpack:"l"`
}
