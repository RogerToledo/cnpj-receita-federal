package file

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadSaveCSV(paths []string) {

	newPath := removeQuotes(paths[1])

	file, err := os.Open(newPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == csv.ErrFieldCount {
			break
		}
		if err != nil {
			panic(err)
		}

		rc := strings.Split(record[0], ";")

		println(rc[0])
		println(rc[1])
		println(rc[2])
		println(rc[3])
		println(rc[4])
		println(rc[5])
		println(rc[6])
		println(rc[7])
		println(rc[8])
		println(rc[9])
		println(rc[10])
		println(rc[11])
		println(rc[12])
		println(rc[13])
		println(rc[14])
		println(rc[15])
		println(rc[16])
		println(rc[17])
		println(rc[18])
		println(rc[19])
		println(rc[20])
		println(rc[21])
		println(rc[22])
		println(rc[23])
		println(rc[24])
		println(rc[25])
		println(rc[26])
		println(rc[27])
		println(rc[28])
		println(rc[29])

		println()

	}
}

func ReadCSVZIP(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 1024)
		record, err := reader.Read(buf)
		if err == gzip.ErrHeader {
			break
		}
		if err != nil {
			panic(err)
		}
		println(record) // [1 2 3 4]
		println()

	}
}

func removeQuotes(path string) string {
	fileCSV := strings.Replace(path, "ESTABELE", "csv", 1)
	output := strings.Replace(fileCSV, "files/files/estabele", "files/files/csv", 1)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\"", "")
	}

	if err = ioutil.WriteFile(output, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		panic(err)
	}

	return output
}

func ReadDir(dir string) []string {
	fd := make([]string, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fd = append(fd, fmt.Sprintf("%s/%s", dir, file.Name()))
	}

	return fd
}

