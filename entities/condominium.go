package entities

// Eu não sei o que retorna cada linha, essa struct pode ter outro nome que seja mais coerente com o que ela representa
type Condominium struct {

	// Se vc usar alguma biblioteca que lê CSV, ex: https://github.com/gocarina/gocsv, vc pode usar tags para mapear os campos
	// Ex: `csv:"cnpjBasico"` e daí não precisa fazer o parse manualmente
	// Mas como o arquivo é grande, não sei se funcionaria e seria performático, precisa testar.
	CNPJBase                 string `csv:"cnpjBasico"`
	CNPJOrder                string
	CNPJDV                   string
	Identifier               string
	FantasyName              string
	CadastralSituation       string
	CadastralSituationDate   string
	CadastralSituationReason string
	City                     string
	Country                  string
	StartDate                string
	PrincipalCNAE            string
	SecondaryCNAE            string
	StreetType               string
	Street                   string
	Number                   string
	Complement               string
	Neighborhood             string
	CEP                      string
	UF                       string
	Municipality             string
	DDD1                     string
	Phone1                   string
	DDD2                     string
	Phone2                   string
	DDDFax                   string
	Fax                      string
	Email                    string
	SpecialSituation         string
	SpecialSituationDate     string
	Hash                     string
}
