package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Transaction() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil

}

func (r *repository) Upsert(rfb []string) error {
	query := Insert
	cnpj := fmt.Sprintf("%s%s%s%s", rfb[0], rfb[1], rfb[2], rfb[4])
	cnpjBasico := rfb[0]
	cnpjOrdem := rfb[1]
	cnpjDV := rfb[2]
	identificador := rfb[3]
	nomeFantasia := rfb[4]
	situacaoCadastral := rfb[5]
	dataSituacaoCadastral := formatDate(rfb[6])
	motivoSituacaoCadastral := rfb[7]
	nomeCidadeExterior := rfb[8]
	pais := rfb[9]
	dataInicio := formatDate(rfb[10])
	cnaePrincipal := rfb[11]
	cnaeSecundario := rfb[12]
	tipoLogradouro := rfb[13]
	logradouro := rfb[14]
	numero := rfb[15]
	complemento := rfb[16]
	bairro := rfb[17]
	cep := rfb[18]
	uf := rfb[19]
	municipio := rfb[20]
	ddd1 := rfb[21]
	telefone1 := rfb[22]
	ddd2 := rfb[23]
	telefone2 := rfb[24]
	dddFax := rfb[25]
	fax := rfb[26]
	email := rfb[27]
	situacaoEspecial := rfb[28]
	dataSituacaoEspecial := formatDate(rfb[29])
	date := time.Now()
	dateNow := date.Format("2006-01-02 03:04:05")

	_, err := r.db.Exec(
		query,
		cnpj,
		cnpjBasico,
		cnpjOrdem,
		cnpjDV,
		identificador,
		nomeFantasia,
		situacaoCadastral,
		dataSituacaoCadastral,
		motivoSituacaoCadastral,
		nomeCidadeExterior,
		pais,
		dataInicio,
		cnaePrincipal,
		cnaeSecundario,
		tipoLogradouro,
		logradouro,
		numero,
		complemento,
		bairro,
		cep,
		uf,
		municipio,
		ddd1,
		telefone1,
		ddd2,
		telefone2,
		dddFax,
		fax,
		email,
		situacaoEspecial,
		dataSituacaoEspecial,
		dateNow,
	)
	
	return err
}

func formatDate(s string) string {
	if s != "" {
		date, _ := time.Parse("20060102", s)
		df := date.Format("2006-01-02")
		return df
	}

	return time.Time{}.Format("2006-01-02")
}
