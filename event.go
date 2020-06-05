package ctypes

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

//
type EventResponse struct {
}

// Async event response is a response to an AsyncEventRequest
type AsyncEventResponse struct {
	QueueCount int                 `json:"queue_count"`
	Errors     EventResponseErrors `json:"errors"`
}

type EventResponseErrors map[int]Error

/*type RunnableEvent struct {
	Name          string           `json:"name"`
	ID            string           `json:"id"`
	Documentation string           `json:"documentation"` // Markdown format
	Style         ctypes.NodeStyle `json:"style"`
}
*/
