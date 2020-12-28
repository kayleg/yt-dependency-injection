package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"../cache"
	"../database"
	"../provider"
)

type defaultRepositoryProvider struct {
	db           database.Executor
	cacheStorage cache.Executor
}

func (d *defaultRepositoryProvider) Database() database.Executor {
	return d.db
}

func (d *defaultRepositoryProvider) Cache() cache.Executor {
	return d.cacheStorage
}

type apiServer struct {
	provider provider.RepositoryProvider
}

func (a *apiServer) createProviderMiddleware() func(http.Handler) http.Handler {
	providers := map[interface{}]interface{}{
		provider.ContextKey:         a.provider,
		database.ExecutorContextKey: a.provider.Database(),
		cache.ExecutorContextKey:    a.provider.Cache(),
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for key, value := range providers {
				req = req.WithContext(context.WithValue(req.Context(), key, value))
			}
			next.ServeHTTP(w, req)
		})
	}
}

func getAllRecords(rw http.ResponseWriter, req *http.Request) {

	repositories := req.Context().Value(provider.ContextKey).(provider.RepositoryProvider)

	cached, err := repositories.Cache().Get("all-records")
	if err != nil {
		// Do something about the error (log, alert, etc)
		fmt.Println("Failed to get info from cache")
	}

	if cached != nil {
		_, _ = rw.Write(cached)
		return
	}

	records, err := repositories.Database().LookupAll(req.Context(), "records")
	if err != nil {
		if err == database.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer

	if err = json.NewEncoder(&buffer).Encode(records); err != nil {
		fmt.Println("Failed to encode json", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = repositories.Cache().Set("all-records", buffer.Bytes())

	_, _ = rw.Write(buffer.Bytes())

}

func main() {
	repositories := &defaultRepositoryProvider{
		db:           database.NewInMemoryDB(),
		cacheStorage: cache.NewInMemoryCache(),
	}

	server := &apiServer{
		provider: repositories,
	}

	router := http.NewServeMux()
	router.HandleFunc("/", getAllRecords)

	_ = http.ListenAndServe(":8080", server.createProviderMiddleware()(router))

}
