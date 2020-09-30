package ctypes

import (
	"time"

	"github.com/google/uuid"
)

var (
	CTTEnvID       = uuid.Must(uuid.NewRandom())
	CTTUserGroupID = uuid.Must(uuid.NewRandom())
	CTTUserID      = uuid.Must(uuid.NewRandom())
)

var ContextTestTree = Context{
	Name: "environment",
	ID:   CTTEnvID,
	Memory: []MemoryContainer{
		{
			Name:    "data",
			Type:    MCTypeSession,
			Exposed: false,
			Data: Mem{
				"sub": Mem{
					"subfld": "woopah",
				},
			},
		},
	},
	Child: &Context{
		Name: "user_group",
		ID:   CTTUserGroupID,
		Memory: []MemoryContainer{
			{
				Name:    "data",
				Type:    MCTypeSession,
				Exposed: false,
				Data: Mem{
					"str":      "heyo",
					"numstr":   "98",
					"num":      5,
					"fl":       10.5,
					"flstring": "5.89",
				},
			},
		},
		Child: &Context{
			Name: "user",
			ID:   CTTUserID,
			Memory: []MemoryContainer{
				{
					Name:    "data",
					Type:    MCTypeSession,
					Exposed: false,
					Data: Mem{
						"str":      "heyo",
						"numstr":   "98",
						"num":      5,
						"fl":       0.557,
						"flstring": "5.89",
					},
				},
			},
		},
	},
}

func StrPtr(str string) *string {
	return &str
}

func TimePtr(time time.Time) *CustomTime {
	return &CustomTime{time}
}

func StringSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}