package ctypes

import (
	"time"

	"github.com/google/uuid"
)

var (
	envID       = uuid.Must(uuid.NewRandom())
	userGroupID = uuid.Must(uuid.NewRandom())
	userID      = uuid.Must(uuid.NewRandom())
)

var ContextTestTree = Context{
	Name: "environment",
	ID:   envID,
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
	Children: []Context{
		{
			Name: "user_group",
			ID:   userGroupID,
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
			Children: []Context{
				{
					Name: "user",
					ID:   userID,
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
		},
	},
}

func StrPtr(str string) *string {
	return &str
}

func TimePtr(time time.Time) *CustomTime {
	return &CustomTime{time}
}
