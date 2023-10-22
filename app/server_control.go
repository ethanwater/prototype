package app

import (
	"context"
	"net/http"
	"os"

	"github.com/ServiceWeaver/weaver/metrics"
)

var databaseAdd = metrics.NewCounter("USER added", "values addded to database")

type serverControls struct{}

func (s serverControls) base() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
	})
}
func (serverControls) kill(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger(ctx).Warn("vivian: kill server...", "status code", 0)
		os.Exit(0)
	})
}

func (serverControls) echo(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		app.Logger(ctx).Debug("echo", "echo", query)
	})
}

func (serverControls) add(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		result, err := app.database.Exec("INSERT INTO users (name) VALUES (?)", query)
		if err != nil {
			app.Logger(ctx).Debug("vivian: FAIL! 'database.add_user'", "err", err)
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				app.Logger(ctx).Error("vivian: ERROR! 'cannot establish ID'", err)
			} else {
				databaseAdd.Inc()
				app.Logger(ctx).Debug("vivian: SUCCESS! 'database.add_user'", "id", id, "user", query)
			}
		}
	})
}
