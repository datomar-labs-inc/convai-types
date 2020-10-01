package ctypes

import (
	"testing"

	"github.com/google/uuid"
)

func TestExpandUUID(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())

	stripped := StripUUID(id)
	expanded := ExpandUUID(stripped)

	if expanded.String() != id.String() {
		t.Error("invalid uuid")
	}
}
