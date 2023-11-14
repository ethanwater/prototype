package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"vivianlab/models"

	"github.com/go-sql-driver/mysql"
)

const name = "vivian_users"
const address = "127.0.0.1:3306"

func EstablishLinkDatabase(ctx context.Context) (*sql.DB, error) {
	database := FetchDatabase(ctx)

	ping := database.Ping()
	if ping != nil {
		return database, ping
	}

	return database, nil
}

func FetchDatabase(ctx context.Context) *sql.DB {
	config := mysql.Config{
		User:   "root",
		Net:    "tcp",
		Addr:   address,
		DBName: name,
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

	// Use a prepared statement
	stmt, err := database.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var acc models.Account
	// Execute the prepared statement with the email parameter
	err = stmt.QueryRow(email).Scan(&acc.ID, &acc.Alias, &acc.Name, &acc.Email, &acc.Password, &acc.Tier)
	if err != nil {
		// Handle the error appropriately
		if err == sql.ErrNoRows {
			// No rows found, return a specific error or handle it accordingly
			return models.Account{}, fmt.Errorf("no account found for email: %w", err)
		}
		return models.Account{}, fmt.Errorf("failed to fetch account: %w", err)
	}
	return acc, nil
}
