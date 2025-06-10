package main

import (
	"flag"
	"log"
)

func main() {
	// Define uma flag de linha de comando para escolher o modo de execução.
	// O valor padrão é "client".
	mode := flag.String("mode", "client", "run in 'client' or 'server' mode")
	
	// Analisa as flags fornecidas na linha de comando.
	flag.Parse()

	// Inicia o programa no modo apropriado com base na flag.
	switch *mode {
	case "server":
		runServer() // Função definida em server.go
	case "client":
		runClient() // Função definida em client.go
	default:
		// Se um modo inválido for fornecido, exibe um erro e encerra.
		log.Fatalf("Modo desconhecido: %s. Use 'client' ou 'server'.", *mode)
	}
}
