package main

import (
	"fmt"

	"github.com/me/rfb/db"
	"github.com/me/rfb/processor"
)

func main() {
	db, err := db.NewDB()
	defer db.Close()
	
	if err != nil {
		fmt.Println(err)
		panic("Was not possible to connect to database")
	}
	defer db.Close()

	processor.Process(db)
}
