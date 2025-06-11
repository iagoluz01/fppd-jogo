package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Registrar o tipo usado em interface{}
func init() {
	// Registrar os tipos usados na interface{}
	gob.Register(&EstadoJogo{})
	gob.Register(map[int]JogadorInfo{})
	gob.Register([][]Elemento{})
}

// ServidorJogo implementa as funções RPC do servidor
type ServidorJogo struct {
	// Campos necessários para o servidor
	estado      *EstadoJogo
	mapaArquivo string
}

// ProcessarComando processa um comando recebido de um cliente
func (s *ServidorJogo) ProcessarComando(cmd *Comando, resposta *Resposta) error {
	log.Printf("Recebido comando: %v de %s", cmd.Parametro, cmd.Jogador)

	if cmd.Parametro == "mover" {
		found := false
		// Procura o jogador na lista (usando o nome)
		for id, jogador := range s.estado.Jogadores {
			if jogador.Nome == cmd.Jogador {
				found = true
				dx, dy := 0, 0
				switch cmd.TeclaRune {
				case 'w', 'W':
					dy = -1
				case 's', 'S':
					dy = 1
				case 'a', 'A':
					dx = -1
				case 'd', 'D':
					dx = 1
				}
				newX := jogador.PosX + dx
				newY := jogador.PosY + dy
				// Verifica se as coordenadas são válidas e se o elemento não bloqueia a passagem
				if newY >= 0 && newY < len(s.estado.ElementosMapa) &&
					newX >= 0 && newX < len(s.estado.ElementosMapa[newY]) &&
					!s.estado.ElementosMapa[newY][newX].Tangivel {
					jogador.PosX = newX
					jogador.PosY = newY
					s.estado.Jogadores[id] = jogador
				}
				break
			}
		}
		// Se o jogador não for encontrado, adiciona-o com a posição inicial definida pelo mapa
		if !found {
			dx, dy := 0, 0
			switch cmd.TeclaRune {
			case 'w', 'W':
				dy = -1
			case 's', 'S':
				dy = 1
			case 'a', 'A':
				dx = -1
			case 'd', 'D':
				dx = 1
			}
			newPlayer := JogadorInfo{
				ID:      len(s.estado.Jogadores) + 1,
				Nome:    cmd.Jogador,
				PosX:    s.estado.StartX + dx,
				PosY:    s.estado.StartY + dy,
				Simbolo: Personagem.Simbolo,
				Cor:     CorCinzaEscuro,
			}
			s.estado.Jogadores[newPlayer.ID] = newPlayer
		}
	}

	resposta.Sucesso = true
	resposta.Mensagem = "Comando processado com sucesso"
	resposta.EstadoJogo = s.estado
	return nil
}

// ObterEstado retorna o estado atual do jogo
func (s *ServidorJogo) ObterEstado(cmd *Comando, resposta *Resposta) error {
	resposta.Sucesso = true
	resposta.EstadoJogo = s.estado
	return nil
}

// IniciarServidor inicializa o servidor RPC
func IniciarServidor(porta string, mapaArquivo string) {
	// Criar o objeto servidor e carregar o estado inicial
	servidor := &ServidorJogo{
		estado: &EstadoJogo{
			Jogadores:     make(map[int]JogadorInfo),
			ElementosMapa: [][]Elemento{},
			Mensagens:     []string{},
		},
		mapaArquivo: mapaArquivo,
	}
	// Carregar o mapa do arquivo usando um Jogo temporário
	game := jogoNovo()
	if err := jogoCarregarMapa(mapaArquivo, &game); err != nil {
		log.Fatalf("Erro ao carregar mapa: %v", err)
	}
	servidor.estado.ElementosMapa = game.Mapa
	servidor.estado.StartX = game.PosX
	servidor.estado.StartY = game.PosY

	// Registrar o servidor RPC
	rpc.Register(servidor)

	// Iniciar o listener TCP
	enderecoCompleto := ":" + porta
	listener, err := net.Listen("tcp", enderecoCompleto)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}

	fmt.Printf("Servidor RPC inicializado com sucesso na porta %s\n", porta)

	// Aceitar conexões
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}

		fmt.Printf("Nova conexão: %s\n", conn.RemoteAddr())
		go rpc.ServeConn(conn)
	}
}
