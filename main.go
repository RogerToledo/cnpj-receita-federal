package main

import (
	"fmt"

	"github.com/me/rfb/db"
	"github.com/me/rfb/processor"
	"github.com/me/rfb/repository"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		fmt.Println(err)
		panic("Was not possible to connect to database")
	}
	defer db.Close()

	repository := repository.NewRepository(db)
	processor := processor.NewProcessor(repository)
	
	if err != nil {
		fmt.Println(err)
		panic("Was not possible to connect to database")
	}
	defer db.Close()

	processor.Process()
}
