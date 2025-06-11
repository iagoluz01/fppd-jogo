// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"log"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	Simbolo  rune // Letra maiúscula para exportar o campo
	Cor      Cor
	CorFundo Cor
	Tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento        // grade 2D representando o mapa
	PosX, PosY      int                 // posição atual do personagem
	UltimoVisitado  Elemento            // elemento que estava na posição do personagem antes de mover
	StatusMsg       string              // mensagem para a barra de status
	Cliente         *ClienteJogo        // referência ao cliente para modo multiplayer
	OutrosJogadores map[int]JogadorInfo // informações sobre outros jogadores
}

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{
		UltimoVisitado:  Vazio,
		OutrosJogadores: make(map[int]JogadorInfo),
	}
}

// Cria uma nova instância do jogo para modo multiplayer
func jogoNovoMultiplayer(cliente *ClienteJogo) Jogo {
	jogo := Jogo{
		UltimoVisitado:  Vazio,
		Cliente:         cliente,
		OutrosJogadores: make(map[int]JogadorInfo),
	}
	return jogo
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.Simbolo:
				e = Parede
			case Inimigo.Simbolo:
				e = Inimigo
			case Vegetacao.Simbolo:
				e = Vegetacao
			case Personagem.Simbolo:
				// Registra a posição inicial e remove o símbolo do jogador do mapa
				jogo.PosX, jogo.PosY = x, y
				e = Vazio
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	return scanner.Err()
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].Tangivel {
		return false
	}

	// Em modo multiplayer, verificar se há outro jogador na posição
	if jogo.OutrosJogadores != nil {
		for _, outroJogador := range jogo.OutrosJogadores {
			if outroJogador.PosX == x && outroJogador.PosY == y {
				return false
			}
		}
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

// Atualiza o estado do jogo com base no estado recebido do servidor
func jogoAtualizarEstadoMultiplayer(jogo *Jogo) {
	if jogo.Cliente == nil || clienteRPC == nil {
		return
	}

	// Obter o estado atual do servidor (já atualizado pela goroutine)
	estado, err := clienteRPC.Estado()
	if err != nil {
		log.Printf("Erro ao obter estado do jogo: %v", err)
		return
	}

	// Agora podemos acessar diretamente os campos do EstadoJogo
	if len(estado.Jogadores) > 0 {
		// Processar jogadores...
	}

	// Para ElementosMapa
	for i := 0; i < len(estado.ElementosMapa); i++ {
		// Processar elementos...
	}

	// Para Jogadores novamente
	if len(estado.Jogadores) > 0 {
		// Processar jogadores...
	}

	// Para Mensagens
	for i := 0; i < len(estado.Mensagens); i++ {
		// Processar mensagens...
	}
}
