package ctypes

import (
	"testing"
)

func TestContextTreeSlice_GetContextByName(t *testing.T) {
	// User structure
	cts := ContextTreeSlice{
		Name: "Environment",
		Children: []ContextTreeSlice{
			{
				Name: "User Group",
				Children: []ContextTreeSlice{
					{
						Name: "Channel User",
						Children: []ContextTreeSlice{
							{
								Name: "Fake Child",
							},
						},
					},
				},
			},
			{
				Name:     "Bob",
				Children: []ContextTreeSlice{},
			},
		},
	}

	/*
	* Check first level of recursion
	 */
	res, exists := cts.GetContextByName("User Group")
	if !exists {
		t.Error("expected user group context to exist")
		return
	}

	if res.Name != cts.Children[0].Name {
		t.Error("expected result to equal first child of cts")
	}

	/*
	* Check second first level of recursion
	 */
	res3, exists3 := cts.GetContextByName("Bob")
	if !exists3 {
		t.Error("expected user group context to exist")
		return
	}

	if res3.Name != cts.Children[1].Name {
		t.Error("expected result to equal first child of cts")
	}

	/*
	* Check second level of recursion
	 */
	res2, exists2 := cts.GetContextByName("Channel User")
	if !exists2 {
		t.Error("expected channel user context to exist")
		return
	}

	if res2.Name != cts.Children[0].Children[0].Name {
		t.Error("expected result to equal first child of first child of cts")
	}

	/*
	* Check third level of recursion
	 */
	res4, exists4 := cts.GetContextByName("Fake Child")
	if !exists4 {
		t.Error("expected channel user context to exist")
		return
	}

	if res4.Name != cts.Children[0].Children[0].Children[0].Name {
		t.Error("expected result to equal first child of first child of first child of cts")
	}
}

func TestContextTreeSlice_GetContextByRef(t *testing.T) {
	ref1 := "1"
	ref2 := "2"
	ref3 := "3"
	ref4 := "4"
	// User structure
	cts := ContextTreeSlice{
		Name: "Environment",
		Ref:  &ref1,
		Children: []ContextTreeSlice{
			{
				Name: "User Group",
				Children: []ContextTreeSlice{
					{
						Name: "Channel User",
						Ref:  &ref2,
						Children: []ContextTreeSlice{
							{
								Name: "Fake Child",
								Ref:  &ref3,
							},
						},
					},
				},
			},
			{
				Name:     "Bob",
				Ref:      &ref4,
				Children: []ContextTreeSlice{},
			},
		},
	}

	/*
	* Check first level hierarchy
	*/
	res, exists := cts.GetContextByRef(ref1)
	if !exists {
		t.Error("expected ref1 context to exist")
		return
	}

	if res.Ref != cts.Ref {
		t.Error("expected result to equal first child of cts")
	}

	/*
	* Check second first level of recursion
	 */
	res3, exists3 := cts.GetContextByRef(ref4)
	if !exists3 {
		t.Error("expected user group context to exist")
		return
	}

	if res3.Ref != cts.Children[1].Ref {
		t.Error("expected result to equal first child of cts")
	}

	/*
	* Check second level of recursion
	 */
	res2, exists2 := cts.GetContextByRef(ref2)
	if !exists2 {
		t.Error("expected channel user context to exist")
		return
	}

	if res2.Name != "Channel User" {
		t.Error("expected result to equal first child of first child of cts")
	}

	/*
	* Check third level of recursion
	 */
	res4, exists4 := cts.GetContextByRef(ref3)
	if !exists4 {
		t.Error("expected channel user context to exist")
		return
	}

	if res4.Name != "Fake Child" {
		t.Error("expected result to equal first child of first child of first child of cts")
	}
}
