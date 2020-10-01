package ctypes

import (
	"fmt"
	"strings"
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

func StripUUID(id uuid.UUID) string {
	return strings.ReplaceAll(id.String(), "-", "")
}

func ExpandUUID(strippedUUID string) uuid.UUID {
	return uuid.MustParse(
		fmt.Sprintf("%s-%s-%s-%s-%s",
			substr(strippedUUID, 0, 8),
			substr(strippedUUID, 8, 4),
			substr(strippedUUID, 12, 4),
			substr(strippedUUID, 16, 4),
			substr(strippedUUID, 20, len(strippedUUID)-20),
		),
	)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
