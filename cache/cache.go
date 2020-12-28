package cache

type contextKeyType int

const (
	ExecutorContextKey contextKeyType = iota
)

type Executor interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
}

type InMemory struct {
	data map[string][]byte
}

func (i *InMemory) Get(key string) ([]byte, error) {
	return i.data[key], nil
}

func (i *InMemory) Set(key string, data []byte) error {
	i.data[key] = data
	return nil
}

func NewInMemoryCache( /*config*/ ) *InMemory {
	return &InMemory{
		data: make(map[string][]byte, 100),
	}
}
