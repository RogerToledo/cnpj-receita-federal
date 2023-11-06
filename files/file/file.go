package file

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/repository"
)

func ReadSaveCSV(db *sql.DB, paths []string) error {

	newPath := removeQuotes(paths[10])

	file, err := os.Open(newPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		result, err := reader.Read()

		if err == csv.ErrFieldCount {
			break
		}
		if err != nil {
			panic(err)
		}

		rs := strings.Split(result[0], ";")

		repository := repository.NewRepository(db)

		if err := repository.UpsertTx(rs); err != nil {
			fmt.Errorf("Error saving on database: %v", err)
		}
	}

	return nil
}

func removeQuotes(path string) string {
	fileCSV := strings.Replace(path, "ESTABELE", "csv", 1)
	output := strings.Replace(fileCSV, "files/files/estabele", "files/files/csv", 1)

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\"", "")
	}

	if err = os.WriteFile(output, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		panic(err)
	}

	return output
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
