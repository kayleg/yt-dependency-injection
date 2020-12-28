package provider

import (
	"../cache"
	"../database"
)

type contextKeyType string

const (
	// ContextKey is used by a request to get access to a RepositoryProvider from the request's context
	ContextKey contextKeyType = "dep-provider"
)

type RepositoryProvider interface {
	Database() database.Executor
	Cache() cache.Executor
}
