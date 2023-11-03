package server

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

// deprecated
var database_add_user = metrics.NewCounter("USER added", "values addded to database")
var mu sync.Mutex

type Add interface {
	DatabaseAddAccount(context.Context, string) (string, error)
}

type add struct {
	weaver.Implements[Add]
}

func (a *add) DatabaseAddAccount(ctx context.Context, query string) (string, error) {
	logger := a.Logger(ctx)
	mu.Lock()
	defer mu.Unlock()

	dbx, err := FetchDatabaseData(ctx)
	if err != nil {
		logger.Error("err:", err)
	}

	var status string
	for _, user := range dbx {
		if query == user.Alias {
			status = "user: '" + query + "' already exists"
			return status, nil
		}
	}

	db := FetchDatabase(ctx)
	result, err := db.Exec("INSERT INTO users (name) VALUES (?)", query)
	if err != nil {
		status = addFailErr + query
		logger.Debug(status, "err", err)
	} else {
		id, err := result.LastInsertId()
		if err != nil {
			status = addFailErr + query
			logger.Debug(status, "err", err)
		} else {
			database_add_user.Inc()
			status = "vivian: SUCCESS! database.add_user: " + query
			logger.Debug(status, "id", id, "user", query)
		}
	}
	return status, nil
}
