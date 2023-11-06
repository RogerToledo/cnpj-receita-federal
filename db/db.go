package db

import "database/sql"

type Database interface {
	UpsertTx(transaction *sql.Tx, rfb []string) error
}
