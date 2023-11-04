package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"vivianlab/pkg/models"

	"github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"
)

func EstablishLinkDatabase(ctx context.Context) *sql.DB {
	database := FetchDatabase(ctx)

	go func() {
		ping := database.Ping()
		if ping != nil {
			os.Exit(2)
		}
	}()

	return database
}

func FetchDatabase(ctx context.Context) *sql.DB {
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

	return database
}

func FetchDatabaseData(ctx context.Context) ([]models.Account, error) {
	database := FetchDatabase(ctx)

	rows, err := database.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var accounts []models.Account

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.Alias, &acc.Name, &acc.Email, &acc.Password, &acc.Tier); err != nil {
			return accounts, err
		}
		accounts = append(accounts, acc)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}
	return accounts, nil
}

func FetchAccount(ctx context.Context, email string) (models.Account, error) {
	database := FetchDatabase(ctx)

	rows, err := database.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		//TODO:
		fmt.Println("nope")
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var acc models.Account
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		if err := rows.Scan(&acc.ID, &acc.Alias, &acc.Name, &acc.Email, &acc.Password, &acc.Tier); err != nil {
			return acc, err
		}
	}
	if err = rows.Err(); err != nil {
		return acc, err
	}
	return acc, nil
}
