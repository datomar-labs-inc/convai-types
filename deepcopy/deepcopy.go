// Package deepcopy provides a function for deep copying map[string]interface{}
// values. Inspired by the StackOverflow answer at:
// http://stackoverflow.com/a/28579297/1366283
//
// Uses the golang.org/pkg/encoding/gob package to do this and therefore has the
// same caveats.
// See: https://blog.golang.org/gobs-of-data
// See: https://golang.org/pkg/encoding/gob/
package deepcopy

import (
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(uuid.UUID{})
}

// Map performs a deep copy of the given map m.
func DeepCopy(src interface{}) (map[string]interface{}, error) {
	var dst map[string]interface{}

	if src == nil {
		return nil, fmt.Errorf("src cannot be nil")
	}

	jsb, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal src: %s", err)
	}

	err = json.Unmarshal(jsb, &dst)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}

	return dst, nil
}