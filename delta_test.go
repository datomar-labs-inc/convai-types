package ctypes

import (
	"testing"

	"github.com/google/uuid"
)

func TestApplyDeltaToModule(t *testing.T) {

}

func TestApplyOperationToModule(t *testing.T) {
	module := GraphModule{
		ID: uuid.Must(uuid.NewRandom()),
	}

	// ***************************************
	// *		     Create Node             *
	// ***************************************
	createNode := DeltaOperation{
		Type:     DOCreateNode,
		ModuleID: module.ID,

		CreateNode: &DeltaCreateNode{
			ID:         uuid.Must(uuid.NewRandom()),
			PackageID:  uuid.Nil,
			TypeID:     StrPtr("branch"),
			Version:    StrPtr("0.0.1"),
			ConfigJSON: StrPtr("{}"),

			Layout: Point{},
		},
	}

	err := ApplyOperationToModule(&module, &createNode)
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := module.Nodes[createNode.CreateNode.ID]; !ok {
		t.Error("expected node to exist")
		return
	}

	// ***************************************
	// *		      Move Node              *
	// ***************************************
	moveNode := DeltaOperation{
		Type:     DOMoveNode,
		ModuleID: module.ID,

		MoveNode: &DeltaMoveNode{
			ID: createNode.CreateNode.ID,
			Pos: Point{
				X: 69,
				Y: 420,
			},
		},
	}

	err = ApplyOperationToModule(&module, &moveNode)
	if err != nil {
		t.Error(err)
		return
	}

	if node, ok := module.Nodes[moveNode.MoveNode.ID]; !ok {
		t.Error("expected node to exist")
		return
	} else {
		if node.Layout.X != 69 || node.Layout.Y != 420 {
			t.Error("node position was not correctly updated")
			return
		}
	}

	// TODO write remaining tests
}
