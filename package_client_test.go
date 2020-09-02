package ctypes

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// TODO right now these tests require manually running an instance of datomar-labs-inc/convai-package-template on port 5555

func newPackageClient() *PackageClient {
	return NewPackageClient(&DBPackage{
		BaseURL:    "http://localhost:5555",
		SigningKey: "bubbles",
	})
}

// TODO: automate test
func TestPackageClient_Dispatch(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.Dispatch(&DispatchRequest{
		Dispatches: []DispatchCall{
			{
				RequestID:       uuid.Must(uuid.NewRandom()),
				ID:              "example_dispatch",
				ContextTree:     ContextTreeSlice{},
				MessageBody:     XMLResponse{},
				PackageSettings: MemoryContainer{},
				Sequence:        0,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_DispatchMock(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.DispatchMock(&DispatchRequest{
		Dispatches: []DispatchCall{
			{
				RequestID:       uuid.Must(uuid.NewRandom()),
				ID:              "example_dispatch",
				ContextTree:     ContextTreeSlice{},
				MessageBody:     XMLResponse{},
				PackageSettings: MemoryContainer{},
				Sequence:        0,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_ExecuteLink(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.ExecuteLink(&LinkExecutionRequest{
		Calls: []LinkCall{
			{
				TypeID:          "example_link",
				Version:         "0.0.1",
				Config:          MemoryContainer{},
				PackageSettings: MemoryContainer{},
				Memory:          nil,
				Sequence:        0,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_ExecuteLinkMock(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.ExecuteLinkMock(&LinkExecutionRequest{
		Calls: []LinkCall{
			{
				TypeID:          "example_link",
				Version:         "0.0.1",
				Config:          MemoryContainer{},
				PackageSettings: MemoryContainer{},
				Memory:          nil,
				Sequence:        0,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_ExecuteNode(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.ExecuteNode(&NodeCall{
		TypeID:  "example_node",
		Version: "0.0.1",
		Config: MemoryContainer{
			Data: Mem{
				"config": `{"field_1":"value"}`,
			},
		},
		PackageSettings: MemoryContainer{},
		Memory:          nil,
		Sequence:        0,
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_ExecuteNodeMock(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.ExecuteNodeMock(&NodeCall{
		TypeID:  "example_node",
		Version: "0.0.1",
		Config: MemoryContainer{
			Data: Mem{
				"config": `{"field_1":"value"}`,
			},
		},
		PackageSettings: MemoryContainer{},
		Memory:          nil,
		Sequence:        0,
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_FetchManifest(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.FetchManifest()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_GetAsset(t *testing.T) {
	// TODO
}

func TestPackageClient_GetAssetBytes(t *testing.T) {
	// TODO
}
