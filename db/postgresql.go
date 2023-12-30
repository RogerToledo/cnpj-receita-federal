package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/me/rfb/config"
)

// Aqui ele retorna um db, não seria melhor o nome da função ser NewDB?
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

// Se você usar uma variável de ambiente com a connection string, não precisa dessa função
// Ex: os.Getenv("DATABASE_URL")
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

// Sugestão:

/*
func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
*/

// Outra sugestão, baseado nesse código do Ronaldo:
// https://github.com/olxbr/ronaldo/blob/master/db/postgres/postgres.go
// Funções: Client, Close
