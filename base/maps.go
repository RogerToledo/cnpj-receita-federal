package base

// Essa variável é declarada aqui, mas é usada só no processor/file.go, não seria melhor declarar lá?
// Ou então criar uma struct com os campos? Coloquei um exemplo criando uma entity Condominium
// Mapa é mais performático, mas nesse caso não faz diferença e acho que fica mais semântico usar struct
var Fields = map[string]string{
	"cnpjBasico":              "",
	"cnpjOrdem":               "",
	"cnpjDV":                  "",
	"identificador":           "",
	"nomeFantasia":            "",
	"situacaoCadastral":       "",
	"dataSituacaoCadastral":   "",
	"motivoSituacaoCadastral": "",
	"nomeCidadeExterior":      "",
	"pais":                    "",
	"dataInicio":              "",
	"cnaePrincipal":           "",
	"cnaeSecundario":          "",
	"tipoLogradouro":          "",
	"logradouro":              "",
	"numero":                  "",
	"complemento":             "",
	"bairro":                  "",
	"cep":                     "",
	"uf":                      "",
	"municipio":               "",
	"ddd1":                    "",
	"telefone1":               "",
	"ddd2":                    "",
	"telefone2":               "",
	"dddFax":                  "",
	"fax":                     "",
	"email":                   "",
	"situacaoEspecial":        "",
	"dataSituacaoEspecial":    "",
	"hash":                    "",
}
