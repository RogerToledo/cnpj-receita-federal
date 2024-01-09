package processor

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/entity"
	"github.com/me/rfb/repository"
)

const (
	pathTXT    = "processor/files/txt"
	pathEstab  = "processor/files/estabele"
	pathProces = "processor/files/processed"
)

type Processor struct {
	repository repository.Repository
}

func NewProcessor(repository repository.Repository) *Processor {
	return &Processor{repository}
}

func (p Processor) Process() error {
	fmt.Println("Processing ...")

	paths := getFiles(pathEstab)
	files := getFiles(pathTXT)

	if len(files) == 0 {
		for _, path := range paths {
			err := formatTXT(path)
			if err != nil {
				fmt.Printf("Error formatTXT: %s", err)
				return err
			}
		}

		files = getFiles(pathTXT)
	}

	for _, file := range files {
		if err := p.ReadTXTSave(file); err != nil {
			return err
		}
	}

	fmt.Println("Ending Process ...")

	return nil
}

func (p Processor) ReadTXTSave(path string) error {
	fmt.Printf("Reading and Saving path: %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		company := lineToCompany(line)

		err := p.repository.Upsert(company, path)
		if err != nil {
			fmt.Printf("error saving on database: %v\n", err)
		}
	}

	if err = moveFile(path); err != nil {
		fmt.Printf("error moveFile: %v\n", err)
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

	_, err = writeFiles(lines, output)
	if err != nil {
		return fmt.Errorf("error writeFile: %s\n", err)
	}

	return nil
}

func getFiles(dir string) []string {
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

func writeFiles(lines []string, output string) (string, error) {
	const maxLines = 100000
	var (
		newLines   = make([]string, 0)
		countLines = 0
		num        = 0
		file       = ""
	)

	fmt.Println("Writing file ...")

	for _, line := range lines {
		countLines++
		line := strings.ReplaceAll(line, "\"", "")

		couldWrite := canWrite(line)
		if couldWrite {
			lineUTF8 := ISO88591ToUTF8([]byte(line))
			newLines = append(newLines, lineUTF8)
		}

		if len(newLines) == maxLines || countLines == len(lines) {
			numFile := fmt.Sprintf("_%d.txt", num)
			file = strings.Replace(output, ".txt", numFile, 1)

			if err := os.WriteFile(file, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
				return "", fmt.Errorf("error WriteFile: %s\n", err)
			}
			newLines = []string{}
			num++
		}
	}

	return file, nil
}

func moveFile(path string) error {
	file := strings.Replace(path, pathTXT, pathProces, 1)

	if err := os.Rename(path, file); err != nil {
		return fmt.Errorf("error moving file: %s\n", err)
	}

	return nil
}

func lineToCompany(line string) entity.Company {
	if line == "" {
		return entity.Company{}
	}

	slice := strings.Split(line, ";")

	company := entity.Company{
		CNPJ:                     slice[0] + slice[1] + slice[2] + slice[3],
		CNPJBase:                 slice[0],
		CNPJOrder:                slice[1],
		CNPJDV:                   slice[2],
		Identifier:               slice[3],
		FantasyName:              slice[4],
		CadastralSituation:       slice[5],
		CadastralSituationDate:   slice[6],
		CadastralSituationReason: slice[7],
		CityNameExterior:         slice[8],
		Country:                  slice[9],
		StartDate:                slice[10],
		PrincipalCNAE:            slice[11],
		SecondaryCNAE:            slice[12],
		StreetType:               slice[13],
		Street:                   slice[14],
		Number:                   slice[15],
		Complement:               slice[16],
		Neighborhood:             slice[17],
		CEP:                      slice[18],
		UF:                       slice[19],
		Municipality:             slice[20],
		DDD1:                     slice[21],
		Phone1:                   slice[22],
		DDD2:                     slice[23],
		Phone2:                   slice[24],
		DDDFax:                   slice[25],
		Fax:                      slice[26],
		Email:                    slice[27],
		SpecialSituation:         slice[28],
		SpecialSituationDate:     slice[29],
		Hash:                     getLineHash(line),
	}

	return company
}

func ISO88591ToUTF8(iso88591 []byte) string {
	buf := make([]rune, len(iso88591))
	for i, b := range iso88591 {
		buf[i] = rune(b)
	}

	return string(buf)
}

// getLineHash returns the hash of the line with 64 caracters
// Example: "24591259;0001;12;1;LOJA DA INFO;02;20160414;00;;;;;;;;;;;;"
// Return: 37a63bad8b1e059bf4eb8f6885bdad01fbb8342851448d4114288a5304728118
func getLineHash(line string) string {
	sha := sha256.New()

	sha.Write([]byte(line))

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func canWrite(line string) bool {
	lineSplited := strings.Split(line, ";")

	cnaes := []string{
		"2621300", "2622100",
	}

	if line == "" {
		return false
	}

	for _, cnae := range cnaes {
		if strings.Contains(lineSplited[11], cnae) || strings.Contains(lineSplited[12], cnae) {
			return true
		}
	}

	return false
}
