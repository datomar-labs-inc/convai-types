package ctypes

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestContextTreeSlice_GetContextByName(t *testing.T) {
	// User structure
	cts := ContextTreeSlice{
		Name: "Environment",
		Child: &ContextTreeSlice{
			Name: "User Group",
			Child: &ContextTreeSlice{
				Name: "Channel User",
				Child: &ContextTreeSlice{
					Name: "Fake Child",
				},
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

	if res.Name != cts.Child.Name {
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

	if res2.Name != cts.Child.Child.Name {
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

	if res4.Name != cts.Child.Child.Child.Name {
		t.Error("expected result to equal first child of first child of first child of cts")
	}
}

func TestContextTreeSlice_GetContextByRef(t *testing.T) {
	ref1 := "1"
	ref2 := "2"
	ref3 := "3"

	// User structure
	cts := ContextTreeSlice{
		Name: "Environment",
		Ref:  &ref1,
		Child: &ContextTreeSlice{
			Name: "User Group",
			Child: &ContextTreeSlice{
				Name: "Channel User",
				Ref:  &ref2,

				Child: &ContextTreeSlice{
					Name: "Fake Child",
					Ref:  &ref3,
				},
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

func TestContext_WithTransformations(t *testing.T) {
	badTransforms := []Transformation{
		{
			Path:      "user.data0  ff.test",
			Value:     "woah there",
			Operation: OpSet,
		},
		{
			Path:      "environment.data.0superData",
			Value:     "holy moly",
			Operation: OpSet,
		},
	}

	for i, ivp := range badTransforms {
		if ivp.PathValid() {
			t.Error("expected path ", i, "to be invalid")
		}
	}

	_, err := ContextTestTree.WithTransformations(badTransforms)
	if err == nil || !strings.Contains(err.Error(), "invalid transformation path") {
		t.Error("expected invalid transformation path error")
		return
	}

	transforms := []Transformation{
		{
			Path:      "user.data.test",
			Value:     "woah there",
			Operation: OpSet,
		},
		{
			Path:      "environment.data.superData",
			Value:     "holy moly",
			Operation: OpSet,
		},
	}

	for i, ivp := range transforms {
		if !ivp.PathValid() {
			t.Error("expected path", i, "to be valid")
		}
	}

	newTree, err := ContextTestTree.WithTransformations(transforms)
	if err != nil {
		t.Error(err)
		return
	}

	// Test
	userContext, exists := newTree.GetContextByName("user")
	if !exists {
		t.Error("expected user context to exist")
		return
	}

	if data, ok := userContext.Memory[0].Data["test"]; !ok || data != "woah there" {
		t.Error("Data was not set correctly on user")
		return
	}
}

func TestContext_GetData(t *testing.T) {
	// Data should not exist
	d1, ok1 := ContextTestTree.GetData("environment.data.test")
	if ok1 != false || d1 != nil {
		t.Error("expected environment.data.test to not exist, but it did")
	}

	// Data should exist
	d2, ok2 := ContextTestTree.GetData("user.data.str")
	if !ok2 || d2 != "heyo" {
		t.Error("expected user.data.str to equal heyo, but it did not")
	}

	// Data should exist and be a float
	d3, ok3 := ContextTestTree.GetData("user_group.data.fl")
	if !ok3 || d3 != 10.5 {
		t.Error("expected user_group.data.fl to equal 0.557 but it did not")
	}
}

func TestContext_GetDataString(t *testing.T) {
	// Data should not exist
	d1, ok1 := ContextTestTree.GetDataString("environment.data.test")
	if ok1 != false || d1 != "" {
		t.Error("expected environment.data.test to not exist, but it did")
	}

	// Data should exist
	d2, ok2 := ContextTestTree.GetDataString("user.data.str")
	if !ok2 || d2 != "heyo" {
		t.Error("expected user.data.str to equal heyo, but it did not")
	}

	// Data should exist and be a float
	d3, ok3 := ContextTestTree.GetDataString("user_group.data.fl")
	if !ok3 || d3 != "10.5" {
		t.Error("expected user_group.data.fl to equal 0.557 but it did not")
	}

	d4, ok4 := ContextTestTree.GetDataString("environment.data.sub.subfld")
	if !ok4 || d4 != "woopah" {
		t.Error("expected environment.data.sub.subfld to equal woopah, but it did not")
	}
}

func TestContext_ExecuteTemplateString(t *testing.T) {
	tmpl, err := ContextTestTree.ExecuteTemplateString("{{ environment.data.sub.subfld }}")
	if err != nil {
		t.Error(err)
	}

	if tmpl != "woopah" {
		t.Error("invalid evaluation, expected woopah, got", tmpl)
	}

	tmpl, err = ContextTestTree.ExecuteTemplateString("{{ user_group.data.fl | divided_by: 2.0 }}")
	if err != nil {
		t.Error(err)
	}

	if tmpl != "5.25" {
		t.Error("invalid evaluation, expected 5.25, got", tmpl)
	}
}

func TestContext_GetTemplateData(t *testing.T) {
	d := ContextTestTree.GetTemplateData()

	expected := Mem{
		"environment": Mem{
			"data": Mem{
				"sub": Mem{
					"subfld": "woopah",
				},
			},
		},
		"user": Mem{
			"data": Mem{
				"str":      "heyo",
				"numstr":   "98",
				"num":      5,
				"fl":       0.557,
				"flstring": "5.89",
			},
		},
		"user_group": Mem{
			"data": Mem{
				"str":      "heyo",
				"numstr":   "98",
				"num":      5,
				"fl":       10.5,
				"flstring": "5.89",
			},
		},
	}

	expJsb, _ := json.Marshal(expected)
	gotJsb, _ := json.Marshal(d)

	if string(expJsb) != string(gotJsb) {
		t.Error("did not get expected values")
	}
}

func TestContext_IDPath(t *testing.T) {
	idPath := ContextTestTree.IDPath()

	if idPath != fmt.Sprintf("%s.%s.%s", cttEnvID.String(), cttUserGroupID.String(), cttUserID.String()) {
		t.Error("id path invalid")
	}
}
