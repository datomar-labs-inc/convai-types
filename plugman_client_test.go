package ctypes

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// TODO right now these tests require manually running an instance of datomar-labs-inc/Convai/packman on port 8080

func newPackmanClient() *PackmanClient {
	return NewPackmanClient("http://localhost:8080")
}

// TODO: automate test
func TestPackmanClient_ListManifest(t *testing.T) {
	pc := newPackmanClient()

	res, err := pc.ListManifest([]uuid.UUID{})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}
