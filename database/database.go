package database

import (
	"context"
	"errors"
)

type contextKeyType int

const (
	ExecutorContextKey contextKeyType = iota
)

type Executor interface {
	LookupByID(ctx context.Context, tableName string, id uint64) (interface{}, error)
	LookupAll(ctx context.Context, tableName string) ([]interface{}, error)
}

// Sample executor implementation

var (
	ErrNotFound = errors.New("Could not find record")
)

type InMemory struct {
	data map[string]map[uint64]interface{}
}

func (i *InMemory) LookupByID(ctx context.Context, tableName string, id uint64) (interface{}, error) {
	if table, exists := i.data[tableName]; exists {
		if val, exists := table[id]; exists {
			return val, nil
		}
	}
	return nil, ErrNotFound
}

func (i *InMemory) LookupAll(ctx context.Context, tableName string) ([]interface{}, error) {
	if table, exists := i.data[tableName]; exists {
		values := make([]interface{}, 0, len(table))
		for _, value := range table {
			values = append(values, value)
		}
		return values, nil
	}
	return nil, ErrNotFound
}

func NewInMemoryDB( /*config*/ ) *InMemory {
	return &InMemory{
		data: make(map[string]map[uint64]interface{}, 5),
	}
}
