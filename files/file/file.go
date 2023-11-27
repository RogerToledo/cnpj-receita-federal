package file

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/repository"
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
		fields := strings.Split(line, ";")

		repository := repository.NewRepository(db)

		err := repository.Upsert(fields)
		if err != nil {
			fmt.Printf("error saving on database: %v", err)
			continue
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
