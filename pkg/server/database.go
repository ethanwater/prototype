package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"
)

func EstablishLinkDatabase(ctx context.Context) *sql.DB {
	toml, err := toml.LoadFile("config.toml")

	if err != nil {
		fmt.Println("Error ", err.Error())
	}

	config := mysql.Config{
		User:   "root",
		Net:    "tcp",
		Addr:   toml.Get("database.address").(string),
		DBName: toml.Get("database.name").(string),
	}

	database, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		os.Exit(1)
	}

	go func() {
		ping := database.Ping()
		if ping != nil {
			os.Exit(2)
		}
	}()

	return database
}
