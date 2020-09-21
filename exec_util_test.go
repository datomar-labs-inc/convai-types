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

func TestExecutionResult_GetMemoryUpdates(t *testing.T) {
	envID := Newuuid.UUID()

	type fields struct {
		EnvironmentID  uuid.UUID
		InitialContext *Context
		Steps          []Step
	}
	tests := []struct {
		name   string
		fields fields
		want   []MemoryUpdate
	}{
		{
			name: "Test",
			fields: fields{
				EnvironmentID:  envID,
				InitialContext: &ContextTestTree,
				Steps: []Step{
					{
						Node: &NodeExecutionResult{
							Transformations: []Transformation{
								{
									Path:      "environment.data.sub.subfld",
									Value:     "woopah2",
									Operation: OpSet,
								},
								{
									Path:      "environment.data.sub.newfld",
									Value:     "woopah3",
									Operation: OpSet,
								},
								{
									Path:      "user.data.teeeest",
									Value:     48.69,
									Operation: OpSet,
								},
							},
						},
					},
				},
			},
			want: []MemoryUpdate{
				{
					ContextID:     cttEnvID,
					EnvironmentID: envID,
					ContainerType: MCTypeSession,
					ContainerName: "data",
					Transformations: []Transformation{
						{
							Path:      "environment.data.sub.subfld",
							Value:     "woopah2",
							Operation: OpSet,
						},
						{
							Path:      "environment.data.sub.newfld",
							Value:     "woopah3",
							Operation: OpSet,
						},
					},
				},
				{
					ContextID:     cttUserID,
					EnvironmentID: envID,
					ContainerType: MCTypeSession,
					ContainerName: "data",
					Transformations: []Transformation{
						{
							Path:      "user.data.teeeest",
							Value:     48.69,
							Operation: OpSet,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ExecutionResult{
				EnvironmentID:  tt.fields.EnvironmentID,
				InitialContext: tt.fields.InitialContext,
				Steps:          tt.fields.Steps,
			}
			if got := s.GetMemoryUpdates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMemoryUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}
