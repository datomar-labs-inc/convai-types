package ctypes

type Mem map[string]string

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

func (m *MemoryContainer) Put(key, value string) *MemoryContainer {
	m.Data[key] = value
	return m
}

func (m *MemoryContainer) PutInt(key, value string) *MemoryContainer {
	m.Data[key] = value
	return m
}
