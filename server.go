package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// Definições de cores e elementos do jogo no lado do servidor.
var (
	// Cores (valores base do termbox para referência)
	corVermelho    Cor = 2
	corVerde       Cor = 3
	corPreto       Cor = 1
	corBranco      Cor = 8
	corCinzaEscuro Cor = 242 // Um valor aproximado de cinza

	// Elementos do jogo
	elementoParede    = Elemento{'▤', corBranco, corPreto, true}
	elementoInimigo   = Elemento{'☠', corVermelho, 0, true}
	elementoVegetacao = Elemento{'♣', corVerde, 0, false}
	elementoVazio     = Elemento{' ', 0, 0, false}
)

// GameEngine gerencia o estado do jogo no servidor.
type GameEngine struct {
	mutex               sync.Mutex
	state               GameState
	lastSequenceNumbers map[string]int64
}

// runServer inicializa e executa o servidor de jogo.
func runServer() {
	log.Println("Iniciando o servidor do jogo...")
	mapaFile := "mapa.txt"
	if _, err := os.Stat(mapaFile); os.IsNotExist(err) {
		log.Fatalf("Arquivo de mapa '%s' não encontrado.", mapaFile)
	}

	gameEngine := &GameEngine{
		state: GameState{
			Players: make(map[string]*Player),
			Mapa:    [][]Elemento{},
		},
		lastSequenceNumbers: make(map[string]int64),
	}
	if err := gameEngine.loadMap(mapaFile); err != nil {
		log.Fatalf("Erro ao carregar o mapa: %v", err)
	}

	rpc.Register(gameEngine)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erro ao iniciar o listener: ", err)
	}
	log.Println("Servidor escutando na porta 1234")

	http.Serve(listener, nil)
}

// loadMap carrega o mapa do jogo a partir de um arquivo de texto.
func (g *GameEngine) loadMap(nome string) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	for scanner.Scan() {
		var linhaElems []Elemento
		for _, ch := range scanner.Text() {
			var e Elemento
			switch ch {
			case '▤':
				e = elementoParede
			case '☠':
				e = elementoInimigo
			case '♣':
				e = elementoVegetacao
			default:
				e = elementoVazio
			}
			linhaElems = append(linhaElems, e)
		}
		g.state.Mapa = append(g.state.Mapa, linhaElems)
	}
	return scanner.Err()
}

// Connect é o método RPC para um novo jogador se conectar.
func (g *GameEngine) Connect(_ struct{}, reply *ReplyConnect) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	playerID := fmt.Sprintf("player-%d", time.Now().UnixNano())
	startX, startY := g.findStartPosition()

	newPlayer := &Player{
		ID:      playerID,
		X:       startX,
		Y:       startY,
		Simbolo: '☺',
		Cor:     Cor(len(g.state.Players)%6 + 2), // Cores diferentes (2=vermelho, 3=verde, etc.)
	}

	g.state.Players[playerID] = newPlayer
	g.lastSequenceNumbers[playerID] = -1

	reply.PlayerID = playerID
	reply.InitialState = g.state

	log.Printf("Jogador %s conectado na posição (%d, %d).", playerID, startX, startY)
	return nil
}

// GetState é o método RPC para os clientes obterem o estado atual do jogo.
func (g *GameEngine) GetState(_ struct{}, reply *GameState) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	*reply = g.state
	return nil
}

// ExecuteAction é o método RPC para um jogador enviar um comando.
func (g *GameEngine) ExecuteAction(args *Args, reply *GameState) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if args.SequenceNumber <= g.lastSequenceNumbers[args.PlayerID] {
		*reply = g.state
		return nil // Ação já processada, retorna estado atual.
	}
	g.lastSequenceNumbers[args.PlayerID] = args.SequenceNumber

	player, ok := g.state.Players[args.PlayerID]
	if !ok {
		return fmt.Errorf("jogador %s não encontrado", args.PlayerID)
	}

	if args.Tecla == 'w' || args.Tecla == 'a' || args.Tecla == 's' || args.Tecla == 'd' {
		g.movePlayer(player, args.Tecla)
	}

	*reply = g.state
	return nil
}

// movePlayer processa a lógica de movimento de um jogador.
func (g *GameEngine) movePlayer(player *Player, tecla rune) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}

	nx, ny := player.X+dx, player.Y+dy

	if ny >= 0 && ny < len(g.state.Mapa) && nx >= 0 && nx < len(g.state.Mapa[ny]) && !g.state.Mapa[ny][nx].Tangivel {
		player.X = nx
		player.Y = ny
	}
}

// findStartPosition encontra um local válido para um novo jogador aparecer.
func (g *GameEngine) findStartPosition() (int, int) {
	for y, row := range g.state.Mapa {
		for x, elem := range row {
			if !elem.Tangivel {
				occupied := false
				for _, p := range g.state.Players {
					if p.X == x && p.Y == y {
						occupied = true
						break
					}
				}
				if !occupied {
					return x, y
				}
			}
		}
	}
	return 1, 1 // Posição padrão de fallback
}
