package processor

import (
	"os"
	"reflect"
	"testing"
)

func TestGetLineHash(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "Retuns the hash of a line empty",
			input:       "",
			expected:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			description: "Retuns the hash of a line",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    "e7e59d9c340aa2b66d44167ea1ad38a7c4b4a4e8b27f6fd3fbcc5444339ff8e0",
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)
		output := getLineHash(tc.input)
		if output != tc.expected {
			t.Errorf("getLineHash(%q) = %q; want %q", tc.input, output, tc.expected)
		}
	}
}

func TestCanWrite(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    bool
	}{
		{
			description: "Retuns true if the line is empty",
			input:       "",
			expected:    false,
		},
		{
			description: "Retuns true if the line contain the cnae 2621300",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    true,
		},
		{
			description: "Retuns true if the line contain the cnae 2622100",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2622100;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    true,
		},
		{
			description: "Retuns true if the line contain the secundary cnae 2622100",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;7739099;4541205;4541203,4541206,4619200,4753900,6204000,7020400,7112000,2622100,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    true,
		},
		{
			description: "Retuns false if the line without valid cnae",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621400;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;",
			expected:    false,
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)
		output := canWrite(tc.input)
		if output != tc.expected {
			t.Errorf("canWrite(%q) = %v; want %v", tc.input, output, tc.expected)
		}
	}
}

func TestISO88591ToUTF8(t *testing.T) {
	cases := []struct {
		description string
		input       []byte
		expected    string
	}{
		{
			description: "Retuns the string empty if the input is empty",
			input:       []byte(""),
			expected:    "",
		},
		{
			description: "Retuns the string empty if the input is nil",
			input:       nil,
			expected:    "",
		},
		{
			description: "Retuns the string in UTF-8",
			input:       []byte("\xed, \xf3, \xfa, \xe1, \xe9, \xe0, \xea, \xf4, \xe2, \xe3, \xfc, \xef, \xe7, \xf1, \xff, \xfd, \xf2, \xec, \xe8, \xf9, \xe3, \xf5, \xe2, \xea, \xee, \xf4, \xfb, \xe1, \xe9, \xed, \xf3, \xfa, \xe0, \xe8, \xec, \xf2, \xf9, \xe3, \xf5, \xe2, \xea, \xee, \xf4, \xfb, \xe4, \xeb, \xef, \xf6, \xfc, \xff"),
			expected:    "í, ó, ú, á, é, à, ê, ô, â, ã, ü, ï, ç, ñ, ÿ, ý, ò, ì, è, ù, ã, õ, â, ê, î, ô, û, á, é, í, ó, ú, à, è, ì, ò, ù, ã, õ, â, ê, î, ô, û, ä, ë, ï, ö, ü, ÿ",
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)

		output := ISO88591ToUTF8(tc.input)
		if output != tc.expected {
			t.Errorf("ISO88591ToUTF8(%q) = %q; want %q", tc.input, output, tc.expected)
		}
	}
}

func TestLineToMap(t *testing.T) {
	var fields = map[string]string{
		"cnpjBasico":         "24770086",
		"nomeCidadeExterior": "",
		"logradouro":         "PADRE CAMARGO",
		"cep":                "84130000",
		"municipio":          "7735",
		"ddd1":               "00",
		"situacaoEspecial":   "",
		"situacaoCadastral":  "02",
		"ddd2":               "", "telefone2": "",
		"hash":                    "e7e59d9c340aa2b66d44167ea1ad38a7c4b4a4e8b27f6fd3fbcc5444339ff8e0",
		"dataSituacaoCadastral":   "20160510",
		"dataInicio":              "20160510",
		"numero":                  "341",
		"dddFax":                  "",
		"email":                   "",
		"pais":                    "",
		"complemento":             "",
		"cnpjDV":                  "08",
		"telefone1":               "11111111",
		"fax":                     "",
		"cnpjOrdem":               "0001",
		"identificador":           "1",
		"motivoSituacaoCadastral": "00",
		"cnaePrincipal":           "2621300",
		"cnaeSecundario":          "4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300",
		"bairro":                  "CENTRO",
		"uf":                      "PR", "nomeFantasia": "RCPM",
		"tipoLogradouro":       "RUA",
		"dataSituacaoEspecial": "",
	}
	cases := []struct {
		description string
		input       string
		expected    map[string]string
	}{
		{
			description: "Retuns the map empty if the line is empty",
			input:       "",
			expected:    map[string]string{},
		},
		{
			description: "Retuns the map with the key and value",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    fields,
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)

		output := lineToMap(tc.input)
		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("lineToMap(%q) = %q; want %q", tc.input, output, tc.expected)
		}
	}
}

func TestWriteFile(t *testing.T) {
	cases := []struct {
		description string
		input       []string
		path        string
		file        string
		expected    int
	}{
		{
			description: "Retuns the error if the file is empty",
			input:       []string{},
			path:        "../files/txt",
			file:        "../files/txt/test.txt",
			expected:    0,
		},
		{
			description: "shoud write the 1 file",
			input: []string{
				"24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
				"24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
				"24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
				"24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			},
			path:     "../files/txt",
			file:     "../files/txt/test.txt",
			expected: 1,
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)

		output, err := writeFile(tc.input, tc.file)
		if err != nil {
			t.Errorf("writeFile(%q) = %q; want %q", tc.input, output, tc.expected)
		}

		files := getFiles(tc.path)

		if len(files) != tc.expected {
			t.Errorf("WriteFile(%q) = %q; want %q", tc.input, output, tc.expected)
		}

		removeFiles(files)

	}
}

func removeFiles(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}
}
