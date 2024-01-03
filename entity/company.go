package entity

import (
	"time"

	validate "github.com/go-playground/validator/v10"
)

type Company struct {
	CNPJ                     string
	CNPJBase                 string `validate:"required,max=8"`
	CNPJOrder                string `validate:"required,max=4"`
	CNPJDV                   string `validate:"required,max=2"`
	Identifier               string `validate:"max=2"`
	FantasyName              string `validate:"max=250"`
	CadastralSituation       string `validate:"max=2"`
	CadastralSituationDate   string `validate:"max=8"`
	CadastralSituationReason string `validate:"max=2"`
	CityNameExterior         string `validate:"max=250"`
	Country                  string `validate:"max=100"`
	StartDate                string `validate:"max=8"`
	PrincipalCNAE            string `validate:"max=7"`
	SecondaryCNAE            string `validate:"max=3000"`
	StreetType               string `validate:"max=25"`
	Street                   string `validate:"max=250"`
	Number                   string `validate:"max=8"`
	Complement               string `validate:"max=100"`
	Neighborhood             string `validate:"max=100"`
	CEP                      string `validate:"max=8"`
	UF                       string `validate:"max=2"`
	Municipality             string `validate:"max=100"`
	DDD1                     string `validate:"max=3"`
	Phone1                   string `validate:"max=9"`
	DDD2                     string `validate:"max=3"`
	Phone2                   string `validate:"max=9"`
	DDDFax                   string `validate:"max=3"`
	Fax                      string `validate:"max=9"`
	Email                    string `validate:"max=100"`
	SpecialSituation         string `validate:"max=250"`
	SpecialSituationDate     string `validate:"max=8"`
	Hash                     string `validate:"len=64"`
}

type CompanyError = struct {
	CNPJ  string
	Field string
	Value string
	Tag string
	File  string
	Date  string
}

func (c *Company) Validate() []CompanyError {
	date := time.Now()
	dateNow := date.Format("2006-01-02 03:04:05")
	companyErros := make([]CompanyError, 0)

	valid := validate.New()

	if errs := valid.Struct(c) ; errs != nil {
		ve := errs.(validate.ValidationErrors)
		for _, err := range ve {
			companyError := CompanyError{
				CNPJ: c.CNPJ,
				Field: err.Field(),
				Value: err.Value().(string),
				Tag: err.Tag(),
				Date: dateNow,
			}
			companyErros = append(companyErros, companyError)
		}
	}

	return companyErros
}
