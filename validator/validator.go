package validator

import (
	"errors"
	"fmt"
	"time"
)

//Dá pra usar o pacote https://github.com/go-playground/validator pra fazer uma validação mais robusta
// Ex:
/*
type Condominium struct {
 	CNPJBase                 string `validate:"len=8"`
	email                    string `validate:"email"`
	...
}

func (c *Condominium) IsValid() error {
	return validate.Struct(c)
}

*/

type ErrorStruct = struct {
	CNPJ  string
	Field string
	Value string
	Error string
	File  string
	Date  string
}

//Alguns desses tamanhos não faz sentido serem fixos, ex: nomeFantasia, logradouro, bairro, municipio, nomeCidadeExterior, pais, etc
// Nesses casos, o ideal é que o tamanho esteja definido no banco de dados, e caso for tentar dar upsert e o tamanho for maior que o do banco, o próprio banco vai retornar o erro
var fieldsValid = map[string]int{
	"cnpjBasico":              8,
	"cnpjOrdem":               4,
	"cnpjDV":                  2,
	"identificador":           2,
	"nomeFantasia":            250,
	"situacaoCadastral":       2,
	"dataSituacaoCadastral":   8,
	"motivoSituacaoCadastral": 2,
	"nomeCidadeExterior":      250,
	"pais":                    100,
	"dataInicio":              8,
	"cnaePrincipal":           7,
	"cnaeSecundario":          3000,
	"tipoLogradouro":          25,
	"logradouro":              250,
	"numero":                  8,
	"complemento":             100,
	"bairro":                  100,
	"cep":                     8,
	"uf":                      2,
	"municipio":               100,
	"ddd1":                    3,
	"telefone1":               9,
	"ddd2":                    3,
	"telefone2":               9,
	"dddFax":                  3,
	"fax":                     9,
	"email":                   100,
	"situacaoEspecial":        250,
	"dataSituacaoEspecial":    8,
	"hash":                    64,
}

func ValidateData(fields map[string]string, path string) (map[string]string, []ErrorStruct) {
	errors := make([]ErrorStruct, 0)
	cnpj := fmt.Sprintf("%s%s%s%s", fields["cnpjBasico"], fields["cnpjOrdem"], fields["cnpjDV"], fields["identificador"])
	date := time.Now()
	dateNow := date.Format("2006-01-02 03:04:05")

	for k, v := range fields {
		err := validate(k, v)
		if err != nil {
			errors = append(errors, ErrorStruct{
				CNPJ:  cnpj,
				Field: k,
				Value: v,
				Error: err.Error(),
				File:  path,
				Date:  dateNow,
			})
		}
	}

	return fields, errors
}

func validate(field string, value string) error {
	if len(value) > fieldsValid[field] {
		msg := fmt.Sprintf("length %d is greater than %d for field %s", len(value), fieldsValid[field], field)
		return errors.New(msg)
	}

	return nil
}
