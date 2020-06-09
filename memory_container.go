package ctypes

type Mem map[string]interface{}

const (
	LifetimeModule = iota
	LifetimeExecution
	LifetimeSession
	LifetimeContext
)

type MemoryContainer struct {
	Name     string `json:"name"`
	Lifetime int    `json:"lifetime"`
	Exposed  bool   `json:"exposed"`
	Data     Mem    `json:"data"`
}

func (m *MemoryContainer) Put(key string, value interface{}) *MemoryContainer {
	m.Data[key] = value
	return m
}