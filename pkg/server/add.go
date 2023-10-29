package server

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

var database_add_user = metrics.NewCounter("USER added", "values addded to database")
var mu sync.Mutex

type AddUserInterface interface {
	AddUser(context.Context, string) (string, error)
}

type adduser struct {
	weaver.Implements[AddUserInterface]
}

func (a *adduser) AddUser(ctx context.Context, query string) (string, error) {
	logger := a.Logger(ctx)
	mu.Lock()
	defer mu.Unlock()
	b := EstablishLinkDatabase(ctx)
	result, err := b.Exec("INSERT INTO users (name) VALUES (?)", query)
	var status string
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
