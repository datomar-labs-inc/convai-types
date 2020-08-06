package ctypes

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGraphLink_Validate(t *testing.T) {
	sd := uuid.Must(uuid.NewRandom())

	link := GraphLink{
		A: LinkPoint{
			NodeID: &sd,
		},
		B: LinkPoint{
			NodeID: &sd,
		},
	}

	lv1 := link.Validate()
	if lv1 != nil && strings.Contains(lv1.Error(), "link cannot have nil id") {
		// Test passes
	} else {
		t.Errorf("expected link nil id error, got %v", lv1)
	}

	link.ID = uuid.Must(uuid.NewRandom())

	lv2 := link.Validate()
	if lv2 != nil && strings.Contains(lv2.Error(), "link type id missing") {
		// Test passes
	} else {
		t.Errorf("expected link type id error, got %v", lv2)
	}

	link.LinkTypeID = "basic"

	lv3 := link.Validate()
	if lv3 != nil && strings.Contains(lv3.Error(), "link source and destination cannot be identical") {
		// Test passes
	} else {
		t.Errorf("expected link source error, got %v", lv3)
	}

	aID := uuid.Must(uuid.NewRandom())
	link.A.NodeID = &aID

	lv4 := link.Validate()
	if lv4 != nil && strings.Contains(lv4.Error(), "both link points cannot occupy the same position") {
		// Test passes
	} else {
		t.Errorf("expected link position error, got %v", lv4)
	}

	link.A.Position = Point{
		X: 1,
		Y: 1,
	}

	lv5 := link.Validate()
	if lv5 != nil && strings.Contains(lv5.Error(), "config json is invalid") {
		// Test passes
	} else {
		t.Errorf("expected link config json error, got %v", lv5)
	}

	link.ConfigJSON = "{}"

	lv6 := link.Validate()
	if lv6 != nil && strings.Contains(lv6.Error(), "invalid version") {
		// Test passes
	} else {
		t.Errorf("expected link version error, got %v", lv6)
	}

	link.Version = "1.0.0"

	lv7 := link.Validate()
	if lv7 != nil {
		t.Error(lv7)
	}
}

func TestGraphModule_Validate(t *testing.T) {
	type fields struct {
		ID    uuid.UUID
		Nodes map[uuid.UUID]GraphNode
		Links []GraphLink
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GraphModule{
				ID:    tt.fields.ID,
				Nodes: tt.fields.Nodes,
				Links: tt.fields.Links,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO test node validation
