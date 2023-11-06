package repository

import (
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Transaction () (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil

}

func (r *repository) UpsertTx(rfb []string) error {
	query := Insert
	cnpj := fmt.Sprintf("%s%s%s%s", rfb[0], rfb[1], rfb[2], rfb[4])

	_, err := r.db.Exec(
		query,
		cnpj,
		rfb[0],
		rfb[1],
		rfb[2],
		rfb[3],
		rfb[4],
		rfb[5],
		rfb[6],
		rfb[7],
		rfb[8],
		rfb[9],
		rfb[10],
		rfb[11],
		rfb[12],
		rfb[13],
		rfb[14],
		rfb[15],
		rfb[16],
		rfb[17],
		rfb[18],
		rfb[18],
		rfb[19],
		rfb[20],
		rfb[21],
		rfb[22],
		rfb[23],
		rfb[24],
		rfb[25],
		rfb[26],
		rfb[27],
		rfb[28],
		rfb[29],
		rfb[30],
	)

	return err
}
