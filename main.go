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
		fmt.Println(err)
		panic("Was not possible to connect to database")
	}
	defer db.Close()
	
	file.Process(db, files)
}
