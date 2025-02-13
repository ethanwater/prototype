package database

import (
	"context"
	"os"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/go-sql-driver/mysql"
)

type Database interface {
	Init(context.Context) error
	PingDBConnection(context.Context) error
	FetchAccount(context.Context, string) (Account, error)
	AddAccount(context.Context, Account) error
}

type impl struct {
	weaver.Implements[Database]
	weaver.WithConfig[config]

	db *sql.DB
}

type config struct {
	Driver string
	Source string
}

const MaxIdleConns, MaxOpenConns = 10, 20

func (s *impl) Init(ctx context.Context) error {
	database, _ := sql.Open(s.Config().Driver, s.Config().Source)
	s.db = database
	s.db.SetMaxIdleConns(MaxIdleConns)
	s.db.SetMaxOpenConns(MaxOpenConns)

	ping := s.db.Ping()
	if ping != nil {
		fmt.Print("vivian: [failure] mysql -> failed to establish connection. check if the sql server is live\n")
		os.Exit(503)
	}

	s.Logger(ctx).Debug("vivian: [launch] mysql", "connection", s.db.Ping() == nil)
	return s.db.Ping()
}

func (s *impl) PingDBConnection(ctx context.Context) error {
	err := s.db.Ping(); if err != nil {
		s.Logger(ctx).Error("vivian: [failure] mysql", "failed to establish connection", err)
		return err 
	}
	return nil
}

type Account struct {
	weaver.AutoMarshal
	ID       int
	Alias    string
	Name     string
	Email    string
	Password string
	Tier     int
}

func (s *impl) FetchAccount(_ context.Context, email string) (Account, error) {
	var mux sync.Mutex
	mux.Lock()
	defer mux.Unlock()
	var acc Account
	_, err := s.db.Exec("USE vivian_users")
	if err != nil {
		log.Fatal("Error selecting database:", err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return Account{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the email parameter
	err = stmt.QueryRow(email).Scan(&acc.ID, &acc.Alias, &acc.Name, &acc.Email, &acc.Password, &acc.Tier)
	if err != nil {
		// Handle the error appropriately
		if err == sql.ErrNoRows {
			// No rows found, return a specific error or handle it accordingly
			return Account{}, fmt.Errorf("no account found for email: %w", err)
		}
		return Account{}, fmt.Errorf("failed to fetch account: %w", err)
	}
	return acc, nil
}

func (s *impl) AddAccount(_ context.Context, account Account) error {
	var mux sync.Mutex
	mux.Lock()
	defer mux.Unlock()
	_, err := s.db.Exec("USE vivian_users")
	if err != nil {
		log.Fatal("Error selecting database", err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM users")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	acc := account
	err = stmt.QueryRow(acc).Scan(&acc.ID, &acc.Alias, &acc.Name, &acc.Email, &acc.Password, &acc.Tier)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no account found for email: %w", err)
		}
		return fmt.Errorf("failed to fetch account: %w", err)
	}

	return err
}
