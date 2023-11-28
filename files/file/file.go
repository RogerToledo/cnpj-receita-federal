package file

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/repository"
	"github.com/me/rfb/validation"
)

const (
	pathTXT = "files/files/txt"
	pathEstab = "files/files/estabele"
	pathDone = "files/files/done"
)

func Process(db *sql.DB, paths []string) error {
	fmt.Println("Processing ...")
	files := ReadDir(pathTXT)

	if len(files) == 0 {
		for _, path := range paths {
			err := formatTXT(path)
			if err != nil {
				fmt.Printf("Error formatTXT: %s", err)
				return err
			}
		}
	}

	for _, file := range files {
		if err := ReadSaveTXT(db, file); err != nil {
			return err
		}
	}	

	fmt.Println("Ending Process ...")

	return nil
}

func ReadSaveTXT(db *sql.DB, path string) error {
	fmt.Printf("Reading and Saving path: %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := lineToMap(line)

		repository := repository.NewRepository(db)

		err := repository.Upsert(fields, path)
		if err != nil {
			fmt.Printf("error saving on database: %v", err)
		}
	}

	if err = moveFile(path); err != nil {
		fmt.Printf("error moveFile: %v", err)
	}

	return nil
}

func formatTXT(path string) error {
	fmt.Printf("Formating path: %s\n", path)

	fileTXT := strings.Replace(path, "ESTABELE", "txt", 1)
	output := strings.Replace(fileTXT, pathEstab, pathTXT, 1)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error ReadFile: %s", err)
	}

	lines := strings.Split(string(data), "\n")

	_, err = writeFile(lines, output)
	if err != nil {
		return fmt.Errorf("error writeFile: %s", err)
	}

	return nil
}

func ReadDir(dir string) []string {
	fd := make([]string, 0)

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fd = append(fd, fmt.Sprintf("%s/%s", dir, file.Name()))
	}

	return fd
}

func writeFile(lines []string, output string) (string, error) {
	const maxLines = 100000
	newLines := make([]string, 0)
	num := 0
	file := ""

	for _, line := range lines {
		line := strings.ReplaceAll(line, "\"", "")
		newLines = append(newLines, line)

		if len(newLines) == maxLines {
			numFile := fmt.Sprintf("_%d.txt", num)
			file = strings.Replace(output, ".txt", numFile, 1)

			if err := os.WriteFile(file, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
				return "", fmt.Errorf("error WriteFile: %s", err)
			}
			newLines = []string{}
			num++
		}
		continue
	}

	return file, nil
}

func moveFile(path string) error {
	file := strings.Replace(path, pathTXT, pathDone, 1)

	if err := os.Rename(path, file); err != nil {
		return fmt.Errorf("error Rename: %s", err)
	}

	return nil
}

func lineToMap(line string) map[string]string {
	m := make(map[string]string)

	slice := strings.Split(line, ";")

	validation.Fields["cnpjBasico"] = slice[0]
	validation.Fields["cnpjOrdem"] = slice[1]
	validation.Fields["cnpjDV"] = slice[2]
	validation.Fields["identificador"] = slice[3]
	validation.Fields["nomeFantasia"] = slice[4]
	validation.Fields["situacaoCadastral"] = slice[5]
	validation.Fields["dataSituacaoCadastral"] = slice[6]
	validation.Fields["motivoSituacaoCadastral"] = slice[7]
	validation.Fields["nomeCidadeExterior"] = slice[8]
	validation.Fields["pais"] = slice[9]
	validation.Fields["dataInicio"] = slice[10]
	validation.Fields["cnaePrincipal"] = slice[11]
	validation.Fields["cnaeSecundario"] = slice[12]
	validation.Fields["tipoLogradouro"] = slice[13]
	validation.Fields["logradouro"] = slice[14]
	validation.Fields["numero"] = slice[15]
	validation.Fields["complemento"] = slice[16]
	validation.Fields["bairro"] = slice[17]
	validation.Fields["cep"] = slice[18]
	validation.Fields["uf"] = slice[19]
	validation.Fields["municipio"] = slice[20]
	validation.Fields["ddd1"] = slice[21]
	validation.Fields["telefone1"] = slice[22]
	validation.Fields["ddd2"] = slice[23]
	validation.Fields["telefone2"] = slice[24]
	validation.Fields["dddFax"] = slice[25]
	validation.Fields["fax" ] = slice[26]
	validation.Fields["email"] = slice[27]
	validation.Fields["situacaoEspecial"] = slice[28]
	validation.Fields["dataSituacaoEspecial"] = slice[29]

	m = validation.Fields

	return m
}