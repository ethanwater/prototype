package app

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

//go:embed index.html
var indexHTML string
var databaseAdd = metrics.NewCounter("USER added", "values addded to database")

type handleInterface interface {
	InnerEcho(context.Context, string) (string, error)
}
type serverControls struct {
	weaver.Implements[handleInterface]
}

func (s serverControls) base(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if _, err := fmt.Fprint(w, indexHTML); err != nil {
			app.Logger(ctx).Error("error writing index.html", "err", err)
		}
	})
}
func (serverControls) kill(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger(ctx).Warn("vivian: kill server...", "status code", 0)
		os.Exit(0)
	})
}

func (s *serverControls) InnerEcho(ctx context.Context, query string) (string, error) {
	return query, nil
}

func (serverControls) echo(ctx context.Context, app *App) http.Handler {
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
