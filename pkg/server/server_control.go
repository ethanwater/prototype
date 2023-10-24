package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

var databaseAdd = metrics.NewCounter("USER added", "values addded to database")

type handleInterface interface {
	InnerEcho(context.Context, string) (string, error)
}
type serverControls struct {
	weaver.Implements[handleInterface]
}

func (serverControls) kill(ctx context.Context, app *Server) http.Handler {
	app.mu.Lock()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger(ctx).Warn("vivian: kill server...", "status code", 0, http.ErrServerClosed)
		defer app.mu.Unlock()
		os.Exit(0)
	})
}

func (s *serverControls) InnerEcho(ctx context.Context, query string) (string, error) {
	return query, nil
}

func (serverControls) echo(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		echo, err := app.serverControls.Get().InnerEcho(r.Context(), query)
		if err != nil {
			app.Logger(r.Context()).Error("error getting query results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(echo)
		if err != nil {
			app.Logger(r.Context()).Error("error marshaling search results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("error writing search results", "err", err)
		}
		app.Logger(ctx).Debug("echo", "echo", echo)
	})

}

func (serverControls) add(ctx context.Context, app *Server) http.Handler {
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
