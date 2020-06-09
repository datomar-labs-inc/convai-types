package ctypes

import (
	"github.com/google/uuid"
)

type DBDispatch struct {
	ID        string    `db:"id" json:"id"`
	PackageID uuid.UUID `db:"package_id" json:"package_id"`
	Name      string    `db:"name" json:"name"`
	Docs      string    `db:"docs" json:"docs"`
}

type DispatchRequest struct {
	Dispatches []DispatchCall `json:"dispatches"`
}

type DispatchResponse struct {
	Results []DispatchCallResult `json:"dispatch_result"`
}

type DispatchCall struct {
	RequestID       uuid.UUID       `json:"request_id"`
	ID              string          `json:"id"`           // The ID of the type of dispatch being called
	MessageBody     string          `json:"message_body"` // XML format message body that the package should parse, post templating
	PackageSettings MemoryContainer `json:"package_settings"`
	Sequence        int             `json:"sequence"` // The order of the message (the order is per request id)
}

type DispatchCallResult struct {
	RequestID  uuid.UUID  `json:"request_id"`
	Successful bool       `json:"successful"` // Did the dispatch operation succeed
	Logs       []LogEntry `json:"logs"`
	Error      *Error     `json:"error"`
}

type PackageDispatch struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	Documentation string `json:"documentation"` // Markdown format
}
