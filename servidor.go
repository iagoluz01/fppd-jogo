package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// ServidorJogo implementa o servidor RPC do jogo
type ServidorJogo struct {
	estado    EstadoJogo
	mutex     sync.RWMutex
	nextID    int
}

// NovoServidor cria uma nova instância do servidor
func NovoServidor(mapaFile string) (*ServidorJogo, error) {
	servidor := &ServidorJogo{
		estado: EstadoJogo{
			Jogadores: make(map[int]JogadorInfo),
			Mensagens: []string{"Servidor iniciado. Bem-vindo!"},
		},
	}

	// Carregar mapa
	jogoTemp := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogoTemp); err != nil {
		return nil, err
	}

	servidor.estado.ElementosMapa = jogoTemp.Mapa
	
	return servidor, nil
}

// Entrar permite que um novo jogador entre no jogo
func (s *ServidorJogo) Entrar(args *EntrarArgs, reply *EntrarReply) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Gerar ID para o novo jogador
	id := s.nextID
	s.nextID++

	// Encontrar posição livre para o jogador
	posX, posY := s.encontrarPosicaoInicial()
	if posX < 0 || posY < 0 {
		reply.Sucesso = false
		reply.Mensagem = "Não foi possível encontrar posição inicial"
		return nil
	}

	// Criar jogador
	jogador := JogadorInfo{
		ID:      id,
		PosX:    posX,
		PosY:    posY,
		Simbolo: args.Simbolo,
		Cor:     args.Cor,
		Nome:    args.Nome,
	}

	// Adicionar ao estado
	s.estado.Jogadores[id] = jogador
	s.estado.Mensagens = append(s.estado.Mensagens, fmt.Sprintf("Jogador %s entrou no jogo", args.Nome))

	// Preparar resposta
	reply.JogadorID = id
	reply.Sucesso = true
	reply.Mensagem = "Bem-vindo ao jogo!"
	reply.Estado = s.estado

	fmt.Printf("Jogador %s (ID: %d) entrou no jogo\n", args.Nome, id)
	return nil
}

// EnviarComando processa um comando de um jogador
func (s *ServidorJogo) EnviarComando(args *EnviarComandoArgs, reply *EnviarComandoReply) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Verificar se o jogador existe
	jogador, existe := s.estado.Jogadores[args.JogadorID]
	if !existe {
		reply.Sucesso = false
		reply.Mensagem = "Jogador não encontrado"
		return nil
	}

	// Processar o comando
	switch args.Tipo {
	case "mover":
		dx, dy := 0, 0
		switch args.Tecla {
		case 'w': dy = -1 // Move para cima
		case 'a': dx = -1 // Move para a esquerda
		case 's': dy = 1  // Move para baixo
		case 'd': dx = 1  // Move para a direita
		}

		nx, ny := jogador.PosX+dx, jogador.PosY+dy
		// Verificar se o movimento é permitido
		if s.podeMoverPara(nx, ny) {
			jogador.PosX, jogador.PosY = nx, ny
			s.estado.Jogadores[args.JogadorID] = jogador
		}

	case "interagir":
		s.estado.Mensagens = append(s.estado.Mensagens, 
			fmt.Sprintf("%s está interagindo em (%d, %d)", 
				jogador.Nome, jogador.PosX, jogador.PosY))
	}

	reply.Sucesso = true
	reply.Mensagem = "Comando processado com sucesso"
	return nil
}

// ObterEstado retorna o estado atual do jogo
func (s *ServidorJogo) ObterEstado(args *ObterEstadoArgs, reply *ObterEstadoReply) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reply.Estado = s.estado
	reply.Sucesso = true
	return nil
}

// Sair remove um jogador do jogo
func (s *ServidorJogo) Sair(args *SairArgs, reply *SairReply) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	jogador, existe := s.estado.Jogadores[args.JogadorID]
	if !existe {
		reply.Sucesso = false
		reply.Mensagem = "Jogador não encontrado"
		return nil
	}

	// Remover jogador
	delete(s.estado.Jogadores, args.JogadorID)
	s.estado.Mensagens = append(s.estado.Mensagens, fmt.Sprintf("Jogador %s saiu do jogo", jogador.Nome))

	reply.Sucesso = true
	reply.Mensagem = "Você saiu do jogo"
	
	fmt.Printf("Jogador %s (ID: %d) saiu do jogo\n", jogador.Nome, jogador.ID)
	return nil
}

// Funções auxiliares
func (s *ServidorJogo) encontrarPosicaoInicial() (int, int) {
	// Procurar posição livre
	for y := range s.estado.ElementosMapa {
		for x := range s.estado.ElementosMapa[y] {
			if s.podeMoverPara(x, y) {
				return x, y
			}
		}
	}
	return -1, -1 // Não encontrou posição válida
}

func (s *ServidorJogo) podeMoverPara(x, y int) bool {
	// Verificar limites do mapa
	if y < 0 || y >= len(s.estado.ElementosMapa) {
		return false
	}
	if x < 0 || x >= len(s.estado.ElementosMapa[y]) {
		return false
	}

	// Verificar se o elemento é tangível
	if s.estado.ElementosMapa[y][x].Tangivel {
		return false
	}

	// Verificar se há outro jogador na posição
	for _, j := range s.estado.Jogadores {
		if j.PosX == x && j.PosY == y {
			return false
		}
	}

	return true
}

// IniciarServidor inicia o servidor RPC
func IniciarServidor(porta string, mapaFile string) {
	servidor, err := NovoServidor(mapaFile)
	if err != nil {
		log.Fatalf("Erro ao criar servidor: %v", err)
	}

	// Registrar o servidor RPC
	rpc.Register(servidor)
	
	// Configurar o listener TCP
	l, err := net.Listen("tcp", ":"+porta)
	if err != nil {
		log.Fatalf("Erro ao ouvir na porta %s: %v", porta, err)
	}
	
	fmt.Printf("Servidor iniciado na porta %s\n", porta)
	
	// Aceitar conexões
	rpc.Accept(l)
}
