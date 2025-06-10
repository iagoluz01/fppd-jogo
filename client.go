package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

// Variáveis globais para o cliente
var (
	rpcClient      *rpc.Client
	playerID       string
	sequenceNumber int64
	localGameState GameState
	gameMutex      sync.Mutex
)

// runClient inicializa e executa o cliente do jogo.
func runClient() {
	serverAddress := "localhost:1234"
	if len(os.Args) > 2 && os.Args[1] == "-mode" && os.Args[2] == "client" && len(os.Args) > 3 {
		serverAddress = os.Args[3] // Ex: go run . -mode client localhost:1234
	}

	var err error
	rpcClient, err = rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor: %v", err)
	}

	var replyConnect ReplyConnect
	err = rpcClient.Call("GameEngine.Connect", struct{}{}, &replyConnect)
	if err != nil {
		log.Fatalf("Erro ao conectar ao jogo: %v", err)
	}

	playerID = replyConnect.PlayerID
	localGameState = replyConnect.InitialState
	log.Printf("Conectado como %s", playerID)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	go updateLoop() // Goroutine para buscar atualizações do servidor.
	inputLoop()    // Loop principal para capturar entradas do jogador.
}

// updateLoop busca o estado do jogo do servidor em intervalos regulares.
func updateLoop() {
	for {
		var newState GameState
		err := rpcClient.Call("GameEngine.GetState", struct{}{}, &newState)
		if err != nil {
			log.Printf("Erro ao obter estado do jogo: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		gameMutex.Lock()
		localGameState = newState
		gameMutex.Unlock()

		drawGame()
		time.Sleep(100 * time.Millisecond) // Frequência de atualização de 10Hz
	}
}

// inputLoop captura a entrada do jogador e a envia para o servidor.
func inputLoop() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return // Sai do loop e encerra o cliente
			}

			var tecla rune
			if ev.Ch != 0 {
				tecla = ev.Ch
			}

			if tecla == 'w' || tecla == 'a' || tecla == 's' || tecla == 'd' || tecla == 'e' {
				gameMutex.Lock()
				sequenceNumber++
				args := &Args{
					PlayerID:       playerID,
					Tecla:          tecla,
					SequenceNumber: sequenceNumber,
				}
				gameMutex.Unlock()

				var replyState GameState
				err := rpcClient.Call("GameEngine.ExecuteAction", args, &replyState)
				if err != nil {
					log.Printf("Erro ao executar ação: %v", err)
					continue
				}

				// Atualiza o estado local com a resposta e redesenha imediatamente
				gameMutex.Lock()
				localGameState = replyState
				gameMutex.Unlock()
				drawGame()
			}
		}
	}
}

// drawGame renderiza o estado atual do jogo na tela.
func drawGame() {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Desenha o mapa
	for y, linha := range localGameState.Mapa {
		for x, elem := range linha {
			termbox.SetCell(x, y, elem.Simbolo, termbox.Attribute(elem.Cor), termbox.Attribute(elem.CorFundo))
		}
	}

	// Desenha todos os jogadores
	for _, player := range localGameState.Players {
		simbolo := player.Simbolo
		if player.ID == playerID {
			simbolo = '☻' // Símbolo especial para o jogador local
		}
		termbox.SetCell(player.X, player.Y, simbolo, termbox.Attribute(player.Cor), termbox.ColorDefault)
	}

	// Desenha a barra de status
	var statusMsg string
	if p, ok := localGameState.Players[playerID]; ok {
		statusMsg = fmt.Sprintf("ID: %s | Pos: (%d, %d) | Jogadores: %d", p.ID, p.X, p.Y, len(localGameState.Players))
	}
	for i, c := range statusMsg {
		termbox.SetCell(i, len(localGameState.Mapa)+1, c, termbox.ColorWhite, termbox.ColorDefault)
	}
	msg := "Use WASD para mover, E para interagir, ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(localGameState.Mapa)+3, c, termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.Flush()
}

