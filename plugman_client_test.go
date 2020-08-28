package ctypes

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// TODO right now these tests require manually running an instance of datomar-labs-inc/Convai/plugman on port 8080

func newPlugmanClient() *PlugmanClient {
	return NewPlugmanClient("http://localhost:8080")
}

// TODO: automate test
func TestPlugmanClient_ListManifest(t *testing.T) {
	pc := newPlugmanClient()

	res, err := pc.ListManifest([]uuid.UUID{})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}
