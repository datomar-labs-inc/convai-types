package ctypes

import (
	"github.com/google/uuid"
)

// Event is sent from a package to Convai (basically a request)
type Event struct {
	ID              string           `json:"id"`
	ContextTree     ContextTreeSlice `json:"context_tree"`
	Text            string           `json:"text"`
	Transformations []Transformation `json:"transformations"`
}

// Event Request is a request that must be answered with execution result
type EventRequest struct {
	Event Event `json:"event"`
}

// AsyncEventRequest is a request that can be queued and executed without returning the result immediately
type AsyncEventRequest struct {
	Events []Event `json:"events"`
}

// EventResponse is a response to a synchronous event, requiring execution to complete before responding
type EventResponse struct {
}

// Async event response is a response to an AsyncEventRequest
type AsyncEventResponse struct {
	QueueCount int                 `json:"queue_count"`
	Errors     EventResponseErrors `json:"errors"`
}

type EventResponseErrors map[int]Error

type DBEvent struct {
	ID            string    `db:"id" json:"id"`
	PackageID     uuid.UUID `db:"package_id" json:"package_id"`
	Name          string    `db:"name" json:"name"`
	Documentation string    `db:"docs" json:"docs"`
	Style         NodeStyle `db:"style" json:"style"`
}

/*type RunnableEvent struct {
	Name          string           `json:"name"`
	TypeID            string           `json:"id"`
	Documentation string           `json:"documentation"` // Markdown format
	Style         ctypes.NodeStyle `json:"style"`
}
*/
