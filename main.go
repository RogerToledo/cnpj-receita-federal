package main

import (
	"fmt"
	"github.com/me/rfb/db"
	"github.com/me/rfb/processor"
	"os"
)

const dir = "processor/estabele"

// 1 - Acho que o nome do projeto poderia mudar para algo mais semântico, é difícil saber o que é rfb
// 2 - Como só tem 1 docker-compose, pode botar ele na raiz do projeto
func main() {

	// Esta parte não deveria estar dentro do Process?
	files := processor.ReadDir(dir)

	db, err := db.NewConnect()
	if err != nil {
		fmt.Println(err)
		panic("Was not possible to connect to database")
	}
	defer db.Close()

	if err := processor.Process(db, files); err != nil {
		// logar erro

		// Sai do programa com erro
		os.Exit(1)
	}
}
