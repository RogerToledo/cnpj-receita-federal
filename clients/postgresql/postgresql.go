package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/me/rfb/config"
)

type postgresql struct{}

var (
	instance *sqlx.DB
	pg       = new(postgresql)
)

func New() *postgresql {
	return pg
}

func Setup() error {
	if pg.IsConnected() {
		return nil
	}

	if err := pg.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pg.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database")

	return nil

}

func (p *postgresql) Connect() error {
	conn, err := sqlx.Connect("postgres", connectionString())
	if err != nil {
		return err
	}

	conn.SetMaxOpenConns(10)
	conn.SetConnMaxIdleTime(20)
	conn.SetConnMaxLifetime(time.Hour)

	instance = conn

	return nil
}

func (p *postgresql) Client() interface{} {
	if instance == nil {
		return nil
	}

	return instance
}

func (p *postgresql) Close() error {
	if instance == nil {
		return nil
	}

	return instance.Close()
}

func (p *postgresql) IsConnected() bool {
	return instance != nil && p.Ping() == nil
}

func (p *postgresql) Ping() error {
	if instance == nil {
		return nil
	}

	return instance.Ping()
}

func (p *postgresql) Transaction() (*sqlx.Tx) {
	return instance.MustBeginTx(context.Background(), nil)
}

func connectionString() string {
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
