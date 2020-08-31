package ctypes

import (
	"testing"
)

func TestGetDataPathKey(t *testing.T) {
	key1 := GetDataPathKey("test.path.with.keys.aplenty")

	if key1 != "with.keys.aplenty" {
		t.Error("expected key to equal with.keys.aplenty, but got", key1)
	}
}
