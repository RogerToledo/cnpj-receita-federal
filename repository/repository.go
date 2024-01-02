package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/me/rfb/validator"
	"time"
)

type repository struct {
	db *sql.DB
}

// Como já tá dentro do package repository, não precisa do nome do package no nome da função
// Ex: NewRepository -> New
// https://google.github.io/styleguide/go/best-practices#avoid-repetition
func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

// Não tá sendo usado
func (r *repository) Transaction() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil

}

// Esse método chama Save, mas faz o upsert. O ideal é que ele chamasse Upsert
func (r *repository) Save(fields map[string]string, path string) error {
	skipUpsert, err := skipUpsert(r, fields["hash"])
	if err != nil {
		return err
	}

	if !skipUpsert {
		if err := upsert(r, fields, path); err != nil {
			return err
		}
	}

	return nil
}

// Poderia colocar esse método como parte da struct, assim não precisa passar o repository como parâmetro
func upsert(r *repository, fields map[string]string, path string) error {
	// Essa função de validação não seria melhor se fosse chamada antes de chamar o método Save?
	// Assim o método ficaria com unicamente a responsabilidade de salvar no banco
	fields, errs := validator.ValidateData(fields, path)
	if len(errs) != 0 {
		err := insertError(r, errs, path)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Got %d validation errors\n", len(errs))

		return errors.New(msg)
	}

	// dá pra declarar essas variáveis em bloco, assim fica mais legível
	// porém, se usar uma struct como a condominium de exemplo, fica melhor
	query := Insert
	cnpj := fmt.Sprintf("%s%s%s%s", fields["cnpjBasico"], fields["cnpjOrdem"], fields["cnpjDV"], fields["identificador"])
	cnpjBasico := fields["cnpjBasico"]
	cnpjOrdem := fields["cnpjOrdem"]
	cnpjDV := fields["cnpjDV"]
	identificador := fields["identificador"]
	nomeFantasia := fields["nomeFantasia"]
	situacaoCadastral := fields["situacaoCadastral"]
	dataSituacaoCadastral := formatDate(fields["dataSituacaoCadastral"])
	motivoSituacaoCadastral := fields["motivoSituacaoCadastral"]
	nomeCidadeExterior := fields["nomeCidadeExterior"]
	pais := fields["pais"]
	dataInicio := formatDate(fields["dataInicio"])
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
	dataSituacaoEspecial := formatDate(fields["dataSituacaoEspecial"])

	// pq já não usa a variável date como string?
	// ex: date := time.Now().Format("2006-01-02 03:04:05")
	date := time.Now()
	dateNow := date.Format("2006-01-02 03:04:05")

	hash := fields["hash"]

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

// O parâmetro path não está sendo usado
// Poderia colocar esse método como parte da struct, assim não precisa passar o repository como parâmetro
// Podemos pensar em outra forma de armazenar esses erros, talvez em um arquivo de log
// Quando o projeto for para produção, os erros vão para o kibana, daí fica mais fácil trocar
func insertError(r *repository, errors []validator.ErrorStruct, path string) error {
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

// Poderia colocar esse método como parte da struct, assim não precisa passar o repository como parâmetro
func skipUpsert(r *repository, hash string) (bool, error) {
	// Dá pra usar direto o sql.ErrNoRows e daí não precisa dessa constante
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

	// sugestão:
	/*
		err := r.db.QueryRow(query, hash).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			}

			return false, err
		}

		return true, nil
	*/
}

// O nome do parâmetro s não é muito descritivo, poderia ser date
func formatDate(s string) string {
	if s != "" {
		date, _ := time.Parse("20060102", s)
		df := date.Format("2006-01-02")
		return df
	}

	// Aqui se o parâmetro for vazio, vai retornar 0001-01-01, não seria melhor retornar vazio?
	return time.Time{}.Format("2006-01-02")
}
