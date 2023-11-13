package file

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/me/rfb/repository"
)

func Process(db *sql.DB, paths []string) error {
	fmt.Println("Processing ...")

	for _, path := range paths {
		csvPath, err:= formatTXT(path)
		if err != nil {
			fmt.Printf("Error formatCSV: %s", err)
			return err
		}
	
		if err := ReadSaveTXT(db, csvPath); err != nil {
			return err
		}
	}

	fmt.Println("Ending Process ...")

	return nil
}

func ReadSaveTXT(db *sql.DB, path string) error {
	fmt.Printf("Reading path: %s", path)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	countSave := 0
	countSkip := 0

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := strings.Split(line, ";")

		repository := repository.NewRepository(db)

		err := repository.Upsert(fields)
		if err != nil {
			countSkip++
			fmt.Printf("Error saving on database: %v", err)
			continue
		}

		countSave++
		
		fmt.Printf("Path: %s, Saved: %d, Skipped: %d", path, countSave, countSkip)

	}

	return nil
}

func formatTXT(path string) (string, error) {
	fileTXT := strings.Replace(path, "ESTABELE", "txt", 1)
	output := strings.Replace(fileTXT, "files/files/estabele", "files/files/txt", 1)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Error ReadFile: ", err)
	}
	lines := strings.Split(string(data), "\n")

	newLines := make([]string, 0)

	for _, line := range lines {
		line := strings.ReplaceAll(line, "\"", "")
		newLines = append(newLines, line)
	}

	if err = os.WriteFile(output, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", fmt.Errorf("Error WriteFile: ", err)
	}

	return output, nil
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
