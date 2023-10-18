package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"prototype/query"
	"prototype/subscriber"

	"github.com/ServiceWeaver/weaver"
	"github.com/go-sql-driver/mysql"
)

//go:embed index.html
var indexHTML string

var db *sql.DB

type User struct {
	Name  string
	Email string
}

// net/http: https://pkg.go.dev/net/http#pkg-overview
// serviceweaver: https://serviceweaver.dev/docs.html#components

func main() {
	config := mysql.Config{
		User:   "root",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "users",
	}

	var err error
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	query      weaver.Ref[query.GithubUserQuery]
	subscriber weaver.Ref[subscriber.SubscriberInterface]
	listener   weaver.Listener `weaver:"proto"`
}

func addUser(user User) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}
	return id, nil
}

func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("prototype run", "addr:", a.listener)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if _, err := fmt.Fprint(w, indexHTML); err != nil {
			a.Logger(r.Context()).Error("error writing index.html", "err", err)
		}
	})

	http.HandleFunc("/github_user_query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		user, err := a.query.Get().Query(r.Context(), query)
		if err != nil {
			a.Logger(r.Context()).Error("error getting query results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		current_user := User{user[0], user[2]}
		addUser(current_user)

		bytes, err := json.Marshal(user)
		if err != nil {
			a.Logger(r.Context()).Error("error marshaling search results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			a.Logger(r.Context()).Error("error writing search results", "err", err)
		}
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		user, err := a.query.Get().Query(r.Context(), query)
		if err != nil {
			a.Logger(r.Context()).Error("error getting query results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = a.subscriber.Get().InitSubscriber(r.Context(), user[2])
		if err != nil {
			a.Logger(r.Context()).Error("error subsscribing", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(user)
		if err != nil {
			a.Logger(r.Context()).Error("error marshaling search results", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
			a.Logger(r.Context()).Error("error writing search results", "err", err)
		}
	})

	return http.Serve(a.listener, nil)
}
