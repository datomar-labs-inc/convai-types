package ctypes

import (
	"time"
)

func StrPtr(str string) *string {
	return &str
}

func TimePtr(time time.Time) *CustomTime {
	return &CustomTime{time}
}
