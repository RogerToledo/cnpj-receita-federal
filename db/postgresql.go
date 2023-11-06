package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/me/rfb/config"
)

func NewConnect() (*sql.DB, error) {
	connString := connString()
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(10)
	conn.SetConnMaxIdleTime(20)
	conn.SetConnMaxLifetime(time.Hour)

	return conn, nil
}

func connString() string {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	user := config.GetString("postgresql_user")
	password := config.GetString("postgresql_pwd")
	host := config.GetString("postgresql_host")
	database := config.GetString("postgresql_database")

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
}
