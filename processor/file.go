// Esse arquivo não precisa estar dentro da pasta file, pode botar na raiz da pasta processor
// Faz sentido porque tem poucos arquivos, se fossem vários e de vários tipos, nesse caso poderia estar separados em pastas
package processor

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/me/rfb/base"
	"github.com/me/rfb/repository"
	"os"
	"strings"
)

const (
	//Acho que faz sentido manter os arquivos em uma pasta separada, ex: processor/files
	/*
		pathTXT   = "processor/files/txt"
		pathEstab = "processor/files/estabele"
		pathDone  = "processor/files/processed"
	*/

	pathTXT   = "processor/txt"
	pathEstab = "processor/estabele"
	pathDone  = "processor/done" // Sugestão: processed
)

// Sugestão, cria uma struct processor com o db e o repositório, assim vc consegue fazer testes unitários
// Ou passa o repositório como parâmetro para a função ReadTXTSave

// Sugestão para tratar os errors:
// https://github.com/uber-go/guide/blob/master/style.md#errors
// Criar eles aqui como variáveis globais

func Process(db *sql.DB, paths []string) error {
	// Pode usar o log pra logar, ex: log.Info("Processing ...")
	fmt.Println("Processing ...")

	files := ReadDir(pathTXT)

	if len(files) == 0 {
		for _, path := range paths {
			err := formatTXT(path)
			if err != nil {
				// Pode usar o log, mas se esse erro for logado em um nível acima, na função que chama ele, então não precisa logar aqui, só retorne o erro
				fmt.Printf("Error formatTXT: %s", err)
				return err
			}
		}

		files = ReadDir(pathTXT)
	}

	for _, file := range files {
		if err := ReadTXTSave(db, file); err != nil {
			return err
		}
	}

	// Pode usar o log pra logar, ex: log.Info("Ending Process ...")
	fmt.Println("Ending Process ...")

	return nil
}

// Essa função lê o txt, salva no banco e move o arquivo, será que não é melhor separar essas responsabilidades?
// Ex: read, save, move e na função Process chamar elas, ou criar uma função intermediária entre Process e ReadTXTSave que chama essas funções
func ReadTXTSave(db *sql.DB, path string) error {
	fmt.Printf("Reading and Saving path: %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := lineToMap(line)

		// Aqui toda hora tá instanciando um novo repositório, seria melhor instanciar só uma vez
		// Sugestão de já passar o repository como parâmetro da função ou criar uma struct processor com o repositório e colocar essa função lá
		// ex:
		/*
			 type Processor struct {
				repository *repository.Repository
			}

			func (p *Processor) ReadTXTSave(path string) error {
				...
			}
		*/
		repository := repository.NewRepository(db)

		err := repository.Save(fields, path)
		if err != nil {
			// Pode usar o log, mas se esse erro for logado em um nível acima, na função que chama ele, então não precisa logar aqui, só retorne o erro
			fmt.Printf("error saving on database: %v", err)

			return err
		}
	}

	if err = moveFile(path); err != nil {
		// Pode usar o log, mas se esse erro for logado em um nível acima, na função que chama ele, então não precisa logar aqui, só retorne o erro
		fmt.Printf("error moveFile: %v", err)

		return err
	}

	return nil
}

// Pelo que eu entendi, essa função transforma o arquivo ESTABELE em txt com com enconding UTF-8, certo?
// Nesse caso, o nome da função não tá muito condizente com o que ela faz
func formatTXT(path string) error {
	fmt.Printf("Formating path: %s\n", path)

	// Esse processo é para pegar o nome do arquivo e trocar o ESTABELE por txt?
	// Acho que assim talvez fique mais performático do que usar o replace 2x (precisa fazer o teste de benchmark):
	/*
		fileName := filepath.Base(path)
		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
		outputPath := fmt.Sprintf("%s/%s.txt", pathTXT, fileName)
	*/

	fileTXT := strings.Replace(path, "ESTABELE", "txt", 1)
	output := strings.Replace(fileTXT, pathEstab, pathTXT, 1)

	data, err := os.ReadFile(path)
	if err != nil {
		// Caso não for usar o erro como variável global, não precisa logar a palavra error (eu sei que já fiz muito isso, eu era moleque)
		// Pode usar tipo: could not read file: %s
		//https://github.com/uber-go/guide/blob/master/style.md#errors
		return fmt.Errorf("error ReadFile: %s", err)
	}

	lines := strings.Split(string(data), "\n")

	_, err = writeFile(lines, output)
	if err != nil {
		// Mesmo caso do erro acima
		return fmt.Errorf("error writeFile: %s", err)
	}

	return nil
}

// Sugestão: ReadDir -> getFiles
func ReadDir(dir string) []string {
	// Esse nome não tá fácil para entender, sugestão: fileNames
	fd := make([]string, 0)

	files, err := os.ReadDir(dir)
	// Se pode dar erro, então é melhor retornar ele do que panicar
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fd = append(fd, fmt.Sprintf("%s/%s", dir, file.Name()))
	}

	return fd
}

// Na verdade essa função não escreve o arquivo, ela escreve vários arquivos, certo?
// Nesse caso pode ser melhor mudar o nome da função para writeFiles
// Depois podemos pensar em usar go routines aqui
func writeFile(lines []string, output string) (string, error) {
	const maxLines = 100000
	newLines := make([]string, 0)
	countLines := 0
	num := 0
	file := ""

	//https://github.com/uber-go/guide/blob/master/style.md#group-similar-declarations
	// "... Groups are not limited in where they can be used. For example, you can use them inside of functions. ..."
	// Dá pra declarar essas variáveis em bloco, ex:
	/*
		var (
			newLines = make([]string, 0)
			countLines = 0
			num = 0
			file = ""
		)
	*/

	fmt.Println("Writing file ...")

	for _, line := range lines {

		// Se usar o i no for, não precisa dessa variável
		// ex:
		/*
			for i, line := range lines {...}
		*/
		countLines++

		// Melhor não usar o mesmo nome da variável do for, pode confundir
		line := strings.ReplaceAll(line, "\"", "")

		canWrite := canWrite(line)
		// se não puder escrever, não deveria ter um continue?
		// ex:
		/*
			if !canWrite {
				continue
			}

			lineUTF8 := ISO88591ToUTF8([]byte(line))
			newLines = append(newLines, lineUTF8)
			...
		*/
		if canWrite {
			lineUTF8 := ISO88591ToUTF8([]byte(line))
			newLines = append(newLines, lineUTF8)
		}

		if len(newLines) == maxLines || countLines == len(lines) {
			// Dá pra usar o i do for pra calcular o numFile, daí não precisa da variável num
			// ex:
			/*
				for i, line := range lines {
					...

					numFile := math.Ceil(i/len(lines))
					file := fmt.Sprintf("%s_%d.txt", output, numFile)

					...
				}
			*/
			numFile := fmt.Sprintf("_%d.txt", num)

			// Acho que fica mais claro se passar o nome do arquivo como parâmetro e não o caminho inteiro, pra evitar ficar dando replace, ex:
			/*
				func writeFile(lines []string, outputPath string, fileName string) (string, error) {
				...
				 file := fmt.Sprintf("%s/%s_d%s", outputPath, fileName)
				}
			*/
			// ou então: (assim não muda a extensão)
			/*
				func writeFile(lines []string, fileName string) (string, error) {
				...
				 file = fmt.Sprintf("%s_%s", numFile, output)
				}
			*/

			file = strings.Replace(output, ".txt", numFile, 1)

			// é melhor sepagar o strings.Join(newLines, "\n") em uma variável, assim fica mais legível e fácil de debugar
			// ex:
			/*
				...
				var ErrWriteFile = "could not write file"
				...

				newLinesStr := strings.Join(newLines, "\n")
				if err := os.WriteFile(file, []byte(newLinesStr), 0644); err != nil {
					return "", fmt.Errorf("%s: %v", ErrWriteFile, err)
				}
			*/

			if err := os.WriteFile(file, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
				return "", fmt.Errorf("error WriteFile: %s", err)
			}
			newLines = []string{}

			num++
		}

		//Aqui já tá no final do for, não precisa de continue
		continue
	}

	return file, nil
}

func moveFile(path string) error {
	file := strings.Replace(path, pathTXT, pathDone, 1)

	if err := os.Rename(path, file); err != nil {
		//Sugestão:
		//return fmt.Errorf("could not rename file %s: %s", path, err)
		return fmt.Errorf("error Rename: %s", err)
	}

	return nil
}

/*
	// Ou to + outro nome que signifique o que essa linha representa, cnae?
    // Ou build+nome da struct, ex: buildCondominium
	// Se for performático usar uma biblioteca para fazer o parse do CSV, então não precisa dessa função
	func toCondominium(line string) *entities.Condominium {
		if line == "" {
			return nil
		}

		fields := strings.Split(line, ";")
		condominium := entities.Condominium{
			CNPJBase: fields[0],
			...
		}

		return &condominium
	}

*/
func lineToMap(line string) map[string]string {
	if line == "" {
		return nil
	}

	fields := strings.Split(line, ";")

	// A ideia de usar uma struct é porque ela representa um objeto, então é mais fácil de entender o que ela representa
	// E além disso, a cada linha vai modificar o mesmo Fields, então não faz sentido usar um map
	// Também tem o caso que se for usar go routines depois, usando assim não é thread safe, pode dar inconsistência dos dados
	base.Fields["cnpjBasico"] = fields[0]
	base.Fields["cnpjOrdem"] = fields[1]
	base.Fields["cnpjDV"] = fields[2]
	base.Fields["identificador"] = fields[3]
	base.Fields["nomeFantasia"] = fields[4]
	base.Fields["situacaoCadastral"] = fields[5]
	base.Fields["dataSituacaoCadastral"] = fields[6]
	base.Fields["motivoSituacaoCadastral"] = fields[7]
	base.Fields["nomeCidadeExterior"] = fields[8]
	base.Fields["pais"] = fields[9]
	base.Fields["dataInicio"] = fields[10]
	base.Fields["cnaePrincipal"] = fields[11]
	base.Fields["cnaeSecundario"] = fields[12]
	base.Fields["tipoLogradouro"] = fields[13]
	base.Fields["logradouro"] = fields[14]
	base.Fields["numero"] = fields[15]
	base.Fields["complemento"] = fields[16]
	base.Fields["bairro"] = fields[17]
	base.Fields["cep"] = fields[18]
	base.Fields["uf"] = fields[19]
	base.Fields["municipio"] = fields[20]
	base.Fields["ddd1"] = fields[21]
	base.Fields["telefone1"] = fields[22]
	base.Fields["ddd2"] = fields[23]
	base.Fields["telefone2"] = fields[24]
	base.Fields["dddFax"] = fields[25]
	base.Fields["fax"] = fields[26]
	base.Fields["email"] = fields[27]
	base.Fields["situacaoEspecial"] = fields[28]
	base.Fields["dataSituacaoEspecial"] = fields[29]
	base.Fields["hash"] = getLineHash(line)

	m := base.Fields

	return m
}

// Apesar da função chamar ISO88591ToUTF8 e o parâmetro ser iso88591, o parâmetro é um slice de bytes, então o parâmetro poderia chamar bytes e a função ToUTF8
// Isso pq o parâmetro já é do tipo bytes, então não precisa colocar o tipo na frente da função. Ex:
// func toUTF8(bytes []byte) string { ... } ou convertToUTF8(bytes []byte) string { ... }
func ISO88591ToUTF8(iso88591 []byte) string {
	buf := make([]rune, len(iso88591))
	for i, b := range iso88591 {
		buf[i] = rune(b)
	}

	return string(buf)
}

// Vale um comentário do que essa função faz, ex:
// getLineHash returns the hash that is built from the line content. Ex: "bla;ble;bli;blo;blu" -> "0x123456789"
func getLineHash(line string) string {
	sha := sha256.New()

	sha.Write([]byte(line))

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func canWrite(line string) bool {
	if line == "" {
		return false
	}

	// Porque precisa dos CNAES repetidos?
	// não é possível usar o strings contains só com o número do CNAE? ex:
	/*
		cnaes := []string{"2621300", "2622100"}
		for _, cnae := range cnaes {

		if strings.Contains(line, cnae) {
			return true
		}

		return false
	*/

	cnaes := []string{
		";2621300;", ";2621300,", ",2621300,", ",2621300;",
		";2622100;", ";2622100,", ",2622100,", ";2622100,",
	}

	for _, cnae := range cnaes {
		if strings.Contains(line, cnae) {
			return true
		}
	}

	return false
}
