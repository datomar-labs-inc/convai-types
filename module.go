package ctypes

import (
	"github.com/google/uuid"
)

type GraphModule struct {
	ID    uuid.UUID               `json:"id" msgpack:"i"`
	Nodes map[uuid.UUID]GraphNode `json:"nodes" msgpack:"n"`
	Links map[uuid.UUID]GraphLink `json:"links" msgpack:"l"`
}
