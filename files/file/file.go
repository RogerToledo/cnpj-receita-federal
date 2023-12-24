package file

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/base"
	"github.com/me/rfb/repository"
)

const (
	pathTXT   = "files/files/txt"
	pathEstab = "files/files/estabele"
	pathDone  = "files/files/done"
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
		if err := ReadTXTSave(db, file); err != nil {
			return err
		}
	}

	fmt.Println("Ending Process ...")

	return nil
}

func ReadTXTSave(db *sql.DB, path string) error {
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

		err := repository.Save(fields, path)
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

	UTF8 := ISO88591ToUTF8(data)

	lines := strings.Split(string(UTF8), "\n")

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

		if len(newLines) == maxLines || len(newLines) == len(lines) {
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
	slice := strings.Split(line, ";")

	base.Fields["cnpjBasico"] = slice[0]
	base.Fields["cnpjOrdem"] = slice[1]
	base.Fields["cnpjDV"] = slice[2]
	base.Fields["identificador"] = slice[3]
	base.Fields["nomeFantasia"] = slice[4]
	base.Fields["situacaoCadastral"] = slice[5]
	base.Fields["dataSituacaoCadastral"] = slice[6]
	base.Fields["motivoSituacaoCadastral"] = slice[7]
	base.Fields["nomeCidadeExterior"] = slice[8]
	base.Fields["pais"] = slice[9]
	base.Fields["dataInicio"] = slice[10]
	base.Fields["cnaePrincipal"] = slice[11]
	base.Fields["cnaeSecundario"] = slice[12]
	base.Fields["tipoLogradouro"] = slice[13]
	base.Fields["logradouro"] = slice[14]
	base.Fields["numero"] = slice[15]
	base.Fields["complemento"] = slice[16]
	base.Fields["bairro"] = slice[17]
	base.Fields["cep"] = slice[18]
	base.Fields["uf"] = slice[19]
	base.Fields["municipio"] = slice[20]
	base.Fields["ddd1"] = slice[21]
	base.Fields["telefone1"] = slice[22]
	base.Fields["ddd2"] = slice[23]
	base.Fields["telefone2"] = slice[24]
	base.Fields["dddFax"] = slice[25]
	base.Fields["fax"] = slice[26]
	base.Fields["email"] = slice[27]
	base.Fields["situacaoEspecial"] = slice[28]
	base.Fields["dataSituacaoEspecial"] = slice[29]
	base.Fields["hash"] = getLineHash(line)

	m := base.Fields

	return m
}

func ISO88591ToUTF8(iso88591 []byte) string {
	buf := make([]rune, len(iso88591))
	for i, b := range iso88591 {
		buf[i] = rune(b)
	}

	return string(buf)
}

func getLineHash(line string) string {
	sha := sha256.New()

	sha.Write([]byte(line))

	return fmt.Sprintf("%x", sha.Sum(nil))
}
