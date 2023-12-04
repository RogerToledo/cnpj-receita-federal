# Obtain public data from the Receita Federal of Brazil
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/RogerToledo/rfb)

## Description

This project aims to obtain public data from companies from the Receita Federal of Brazil (RFB) and persist it in a database.

## Status
:construction: Project under construction :construction:

## Project features
- [x] Reads RFB files from files/files/establishe, formats, generates smaller files and saves them in files/files/txt
- [x] Reads files from files/files/txt and validates data entry
- [x] Persists data that has passed validation
- [x] Persists data that has not passed validation in the error table
- [x] Moves files that have been processed to the done folder

## Technologies
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)

## How to run
1. Clone this repository
2. Install [Go](https://golang.org/)
3. Run compose-postgresql.yml located in the tools folder
4. Run sql scripts located in the sql folder to create the tables
5. Download the files from [RFB](https://dados.gov.br/dados/conjuntos-dados/cadastro-nacional-da-pessoa-juridica---cnpj) (Empresas), unzip and move to files/files/establishe folder
6. Run the command `go run main.go` in the root folder

## Author
[Roger Toledo](https://github.com/RogerToledo)

## Contributors
[Thalita Oliveira](https://github.com/thalita)



