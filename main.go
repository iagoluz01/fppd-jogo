// main.go - Loop principal do jogo
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Definir flags para modo cliente e servidor
	modoServidor := flag.Bool("servidor", false, "Iniciar como servidor")
	porta := flag.String("porta", "8080", "Porta para o servidor")
	endereco := flag.String("endereco", "localhost:8080", "Endereço do servidor para conexão do cliente")
	nome := flag.String("nome", "Jogador", "Nome do jogador")
	mapaFile := flag.String("mapa", "mapa.txt", "Arquivo de mapa")
	
	flag.Parse()

	// Verificar o modo de execução
	if *modoServidor {
		// Modo servidor - inicia o servidor RPC
		fmt.Println("Iniciando servidor na porta:", *porta)
		fmt.Println("Usando mapa:", *mapaFile)
		
		// Iniciar o servidor
		IniciarServidor(*porta, *mapaFile)
	} else {
		// Modo cliente - inicia o cliente do jogo
		fmt.Println("Conectando ao servidor:", *endereco)
		fmt.Println("Nome do jogador:", *nome)
		
		// Inicializa a interface (termbox)
		interfaceIniciar()
		defer interfaceFinalizar()
		
		// Conectar ao servidor
		cliente, err := NovoCliente(*endereco, *nome, '☺', CorCinzaEscuro)
		if err != nil {
			fmt.Printf("Erro ao conectar: %v\n", err)
			return
		}
		defer cliente.Sair()
		
		// Criar jogo local
		jogo := jogoNovoMultiplayer(cliente)
		
		// Desenha o estado inicial do jogo
		interfaceDesenharJogoMultiplayer(&jogo)
		
		// Loop principal de entrada
		for {
			evento := interfaceLerEventoTeclado()
			if continuar := personagemExecutarAcaoMultiplayer(evento, &jogo); !continuar {
				break
			}
			interfaceDesenharJogoMultiplayer(&jogo)
		}
	}
}