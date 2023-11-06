package main

import (
	"fmt"

	"github.com/me/rfb/db"
	"github.com/me/rfb/files/file"
)

func main() {
	files := file.ReadDir("files/files/estabele")
	
	db, err := db.NewConnect()
	if err != nil {
		fmt.Errorf("Was not possible to connect to database - %s", err)
	}
	
	file.ReadSaveCSV(db, files)
}
