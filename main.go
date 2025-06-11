// main.go - Loop principal do jogo
package main

import (
	"flag"
	"fmt"
	"log"
	"net" // Adicionar o pacote net
	"time"

	"github.com/nsf/termbox-go"
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

		// Inicializar o cliente RPC
		if err := InicializarClienteRPC(*endereco); err != nil {
			log.Fatalf("Erro ao conectar no servidor: %v", err)
			return
		}

		RegistrarLog("Conexão estabelecida com o servidor %s", *endereco)

		// Inicializa a interface (termbox)
		interfaceIniciar()
		defer interfaceFinalizar()

		// Conectar ao servidor e criar o cliente
		cliente, err := NovoCliente(*nome, *endereco, 1, termbox.ColorWhite) // ID temporário 1, cor branca
		if err != nil {
			RegistrarLog("Erro ao criar cliente: %v", err)
			log.Fatalf("Erro ao conectar: %v\n", err)
			return
		}
		defer cliente.Sair()

		RegistrarLog("Cliente criado com sucesso, nome: %s", *nome)

		// Criar jogo local
		jogo := jogoNovoMultiplayer(cliente)

		// Atualizar jogo com o estado inicial do servidor
		if estado, err := clienteRPC.Estado(); err != nil {
			RegistrarLog("Erro ao obter estado inicial: %v", err)
		} else {
			jogo.Mapa = estado.ElementosMapa
			jogo.OutrosJogadores = estado.Jogadores
			jogo.StatusMsg = "Jogo carregado do servidor"
		}

		RegistrarLog("Jogo multiplayer criado")

		// Desenha o estado inicial do jogo
		interfaceDesenharJogoMultiplayer(&jogo)
		RegistrarLog("Interface inicializada e primeiro quadro desenhado")

		// Novo: Loop de atualização do estado do jogo
		go func() {
			// Atualiza a cada ~16ms (~60 FPS)
			ticker := time.NewTicker(16 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				estado, err := clienteRPC.Estado()
				if err != nil {
					RegistrarLog("Erro ao obter estado: %v", err)
					continue
				}
				// Atualiza o mapa e os jogadores
				jogo.Mapa = estado.ElementosMapa
				jogo.OutrosJogadores = estado.Jogadores
				// Atualiza a posição do jogador local, se encontrado
				for _, jogador := range estado.Jogadores {
					if jogador.Nome == cliente.Nome() {
						jogo.PosX = jogador.PosX
						jogo.PosY = jogador.PosY
						break
					}
				}
				interfaceDesenharJogoMultiplayer(&jogo)
			}
		}()

		// Loop principal de entrada
		for {
			evento := interfaceLerEventoTeclado()
			RegistrarLog("Evento de teclado recebido: %v", evento)

			if continuar := personagemExecutarAcaoMultiplayer(evento, &jogo); !continuar {
				RegistrarLog("Saindo do jogo")
				break
			}
			interfaceDesenharJogoMultiplayer(&jogo)
		}
	}
}

// Função exemplo para tratar a conexão do cliente
func tratarConexao(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Cliente conectado:", conn.RemoteAddr().String())

	// Aqui você pode adicionar o código para comunicação com o cliente
}
