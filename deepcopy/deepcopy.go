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
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(uuid.UUID{})
}

// Map performs a deep copy of the given map m.
func DeepCopy(m map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}