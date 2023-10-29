package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var echoFailErr = "failed: "
var addFailErr = "failed: "

func Echo(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		err := app.echo.Get().Echo(ctx, query)
		if err != nil {
			app.Logger(r.Context()).Error(echoFailErr, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(query)
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

func Add(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		result, err := app.add.Get().AddUser(ctx, query)
		if err != nil {
			app.Logger(r.Context()).Error("cannot add user", "err", err)
		}

		bytes, err := json.Marshal(result)
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

func FetchUsers(ctx context.Context, app *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := FetchDatabaseData(ctx)
		if err != nil {
			app.Logger(r.Context()).Error("cannot add user", "err", err)
		}
		bytes, err := json.Marshal(result)
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
