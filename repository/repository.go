package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/me/rfb/validation"
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

func (r *repository) Upsert(fields map[string]string, path string) error {
	fields, errs := validation.ValidateData(fields, path)
	if len(errs) != 0 {
		err := insertError(r, errs, path)
		if err != nil {
			return err
		}	

		msg := fmt.Sprintf("Got %d validation errors", len(errs))

		return errors.New(msg)
	}

	query := Insert
	cnpj := fmt.Sprintf("%s%s%s%s",  fields["cnpjBasico"],  fields["cnpjOrdem"],  fields["cnpjDV"],  fields["identificador"])
	cnpjBasico := fields["cnpjBasico"]
	cnpjOrdem := fields["cnpjOrdem"]
	cnpjDV := fields["cnpjDV"]
	identificador := fields["identificador"]
	nomeFantasia := fields["nomeFantasia"]
	situacaoCadastral := fields["situacaoCadastral"]
	dataSituacaoCadastral := formatDate( fields["dataSituacaoCadastral"] )
	motivoSituacaoCadastral := fields["motivoSituacaoCadastral"]
	nomeCidadeExterior := fields["nomeCidadeExterior"]
	pais := fields["pais"]
	dataInicio := formatDate( fields["dataInicio"])
	cnaePrincipal := fields["cnaePrincipal"]
	cnaeSecundario := fields["cnaeSecundario"]
	tipoLogradouro := fields["tipoLogradouro"]
	logradouro := fields["logradouro"]
	numero := fields["numero"]
	complemento := fields["complemento"]
	bairro := fields["bairro"]
	cep := fields["cep"]
	uf := fields["uf"]
	municipio := fields["municipio"]
	ddd1 := fields["ddd1"]
	telefone1 := fields["telefone1"]
	ddd2 := fields["ddd2"]
	telefone2 := fields["telefone2"]
	dddFax := fields["dddFax"]
	fax := fields["fax"]
	email := fields["email"]
	situacaoEspecial := fields["situacaoEspecial"]
	dataSituacaoEspecial := formatDate( fields["dataSituacaoEspecial"] )
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

func insertError(r *repository, errors []validation.ErrorStruct, path string) error {
	var err error
	query := InsertError

	for _, errorValue := range errors { 
		_, err = r.db.Exec(
			query,
			errorValue.CNPJ,
			errorValue.Field,
			errorValue.Value,
			errorValue.Error,
			errorValue.File,
			errorValue.Date,
		)
		if err != nil {
			fmt.Printf("Error inserting error: %v", err)
		}
	}

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
