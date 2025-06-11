package main

// JogadorInfo contém informações sobre um jogador conectado
type JogadorInfo struct {
	ID      int
	PosX    int
	PosY    int
	Simbolo rune
	Cor     Cor
	Nome    string
}

// EstadoJogo representa o estado global do jogo no servidor
type EstadoJogo struct {
	Jogadores     map[int]JogadorInfo
	ElementosMapa [][]Elemento
	Mensagens     []string
	StartX        int // posição inicial para novos jogadores
	StartY        int // posição inicial para novos jogadores
}

// Elementos visuais do jogo (com campos exportados)
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)
