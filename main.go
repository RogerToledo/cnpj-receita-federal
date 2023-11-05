package main

import (
	"github.com/me/rfb/files/file"
)

func main() {
	files := file.ReadDir("files/files/estabele")
	file.ReadSaveCSV(files)
}
