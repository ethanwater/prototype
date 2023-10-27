package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

const echoFailErr string = "vivian: FAIL! echo"
const addFailErr string = "vivian: FAIL! database.add_user: "

var databaseAdd = metrics.NewCounter("USER added", "values addded to database")

type handleInterface interface {
	InnerEcho(context.Context, string) (string, error)
	Add(context.Context, string) (string, error)
	//Ping(context.Context, string) (string, error)
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

func (s *serverControls) Add(ctx context.Context, query string) (string, error) {
	return query, nil
}

//	func (s *serverControls) Ping(ctx context.Context, database string) (string, error) {
//		return database, nil
//	}
//
//	func (serverControls) ping(ctx context.Context, app *Server) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ping, err := app.serverControls.Get().Ping(r.Context(), app.db_name)
//			if err != nil {
//				app.Logger(r.Context()).Error("cant retrieve database")
//				http.Error(w, err.Error(), http.StatusInternalServerError)
//			}
//
//			bytes, err := json.Marshal(ping)
//			if err != nil {
//				app.Logger(r.Context()).Error("error marshaling results", "err", err)
//				http.Error(w, err.Error(), http.StatusInternalServerError)
//				return
//			}
//			if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
//				app.Logger(r.Context()).Error("error writing search results", "err", err)
//			}
//
//		})
//	}
func (serverControls) echo(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		echo, err := app.serverControls.Get().InnerEcho(r.Context(), query)
		if err != nil {
			app.Logger(r.Context()).Error(echoFailErr, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(echo)
		if err != nil {
			app.Logger(r.Context()).Error(echoFailErr, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error(echoFailErr, "err", err)
		}
		app.Logger(ctx).Debug("echo", "echo", echo)
	})

}

func (serverControls) add(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var status string
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		result, err := app.database.Exec("INSERT INTO users (name) VALUES (?)", query)
		if err != nil {
			status = addFailErr + query
			app.Logger(ctx).Debug(status, "err", err)
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				status = addFailErr + query
				app.Logger(ctx).Debug(status, "err", err)
			} else {
				databaseAdd.Inc()
				status = "vivian: SUCCESS! database.add_user: " + query
				app.Logger(ctx).Debug(status, "id", id, "user", query)
			}
		}

		//adduser, err := app.serverControls.Get().Add(r.Context(), query)
		//if err != nil {
		//	app.Logger(r.Context()).Error("error getting query results", "err", err)
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		bytes, err := json.Marshal(status)
		if err != nil {
			app.Logger(r.Context()).Error("error marshaling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("error writing search results", "err", err)
		}
	})
}
