package repository

import (
	"database/sql"

	"fmt"
	"time"

	"github.com/me/rfb/entity"
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

func (r *repository) Save(company entity.Company, path string) error {
	skipUpsert, err := skipUpsert(r, company.Hash)
	if err != nil {
		return err
	}

	if !skipUpsert {
		if err := upsert(r, company, path); err != nil {
			return err
		}
	}

	return nil
}

func upsert(r *repository, company entity.Company, path string) error {
	company.CNPJ = fmt.Sprintf("%s%s%s%s", company.CNPJBase, company.CNPJOrder, company.CNPJDV, company.Identifier)

	fieldsFail := company.Validate()
	if len(fieldsFail) != 0 {
		err := insertError(r, fieldsFail, path)
		if err != nil {
			return err
		}

		return nil
	}

	query := Insert
	cnpj := fmt.Sprintf("%s%s%s%s", company.CNPJBase, company.CNPJOrder, company.CNPJDV, company.Identifier)
	cnpjBasico := company.CNPJBase
	cnpjOrdem := company.CNPJOrder
	cnpjDV := company.CNPJDV
	identificador := company.Identifier
	nomeFantasia := company.FantasyName
	situacaoCadastral := company.CadastralSituation
	dataSituacaoCadastral := formatDate(company.CadastralSituationDate)
	motivoSituacaoCadastral := company.CadastralSituationReason
	nomeCidadeExterior := company.CityNameExterior
	pais := company.Country
	dataInicio := formatDate(company.StartDate)
	cnaePrincipal := company.PrincipalCNAE
	cnaeSecundario := company.SecondaryCNAE
	tipoLogradouro := company.StreetType
	logradouro := company.Street
	numero := company.Number
	complemento := company.Complement
	bairro := company.Neighborhood
	cep := company.CEP
	uf := company.UF
	municipio := company.Municipality
	ddd1 := company.DDD1
	telefone1 := company.Phone1
	ddd2 := company.DDD2
	telefone2 := company.Phone2
	dddFax := company.DDDFax
	fax := company.Fax
	email := company.Email
	situacaoEspecial := company.SpecialSituation
	dataSituacaoEspecial := formatDate(company.SpecialSituationDate)
	date := time.Now()
	dateNow := date.Format("2006-01-02 03:04:05")
	hash := company.Hash

	_, err := r.db.Exec(
		query,
		cnpj,                    // 1
		cnpjBasico,              // 2
		cnpjOrdem,               // 3
		cnpjDV,                  // 4
		identificador,           // 5
		nomeFantasia,            // 6
		situacaoCadastral,       // 7
		dataSituacaoCadastral,   // 8
		motivoSituacaoCadastral, // 9
		nomeCidadeExterior,      // 10
		pais,                    // 11
		dataInicio,              // 12
		cnaePrincipal,           // 13
		cnaeSecundario,          // 14
		tipoLogradouro,          // 15
		logradouro,              // 16
		numero,                  // 17
		complemento,             // 18
		bairro,                  // 19
		cep,                     // 20
		uf,                      // 21
		municipio,               // 22
		ddd1,                    // 23
		telefone1,               // 24
		ddd2,                    // 25
		telefone2,               // 26
		dddFax,                  // 27
		fax,                     // 28
		email,                   // 29
		situacaoEspecial,        // 30
		dataSituacaoEspecial,    // 31
		dateNow,                 // 32
		hash,                    // 33
	)

	return err
}

func insertError(r *repository, errors []entity.CompanyError, path string) error {
	var err error
	query := InsertError

	for _, errorValue := range errors {
		_, err = r.db.Exec(
			query,
			errorValue.CNPJ,
			errorValue.Field,
			errorValue.Value,
			errorValue.Tag,
			errorValue.File,
			errorValue.Date,
		)
		if err != nil {
			fmt.Printf("Error inserting error: %v", err)
		}
	}

	return err
}

func skipUpsert(r *repository, hash string) (bool, error) {
	const noRows = "sql: no rows in result set"
	var id int
	query := SelectHash

	err := r.db.QueryRow(query, hash).Scan(&id)
	if err != nil && err.Error() != noRows {
		return false, err
	}

	if err != nil && err.Error() == noRows {
		return false, nil
	}

	return true, nil
}

func formatDate(s string) string {
	if s != "" {
		date, _ := time.Parse("20060102", s)
		df := date.Format("2006-01-02")
		return df
	}

	return time.Time{}.Format("2006-01-02")
}
