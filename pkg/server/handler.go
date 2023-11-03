package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var echoFailErr = "vivian: ERROR!"
var addFailErr = "vivian: ERROR!"

func Split(r rune) bool {
	return r == '&' || r == ','
}

func EchoResponse(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		err := app.echo.Get().EchoResponse(ctx, query)
		if err != nil {
			app.Logger(r.Context()).Error(echoFailErr, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(query)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling resullts", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}
	})
}

//func AddAccount(ctx context.Context, app *App) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		query := strings.TrimSpace(r.URL.Query().Get("q"))
//		result, err := app.add.Get().DatabaseAddAccount(ctx, query)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! cannot add user", "err", err)
//		}
//
//		bytes, err := json.Marshal(result)
//		if err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
//			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
//		}
//	})
//}

func FetchUsers(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := FetchDatabaseData(ctx)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! cannot add user", "err", err)
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}

	})
}

func AccountLogin(ctx context.Context, app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result bool
		query := strings.TrimSpace(r.URL.RawQuery)
		creds := strings.FieldsFunc(query, Split)

		if len(creds) < 2 {
			app.Logger(ctx).Error("vivian: ERROR! missing field, killing handle..")
			return
		}

		email := strings.TrimSpace(creds[0])
		password := strings.TrimSpace(creds[1])

		result, err := app.login.Get().Login(ctx, email, password)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure logging", "err", err)
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure marshalling results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			app.Logger(r.Context()).Error("vivian: ERROR! failure writing search results", "err", err)
		}

	})
}
