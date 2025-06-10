package main

// Cor é um alias para um tipo de cor. Usamos uint16 para que o tipo seja
// facilmente serializável pela biblioteca RPC (gob) e independente do termbox.
type Cor uint16

// Elemento representa qualquer objeto do mapa (parede, personagem, etc.).
type Elemento struct {
	Simbolo  rune
	Cor      Cor
	CorFundo Cor
	Tangivel bool // Indica se o elemento bloqueia a passagem.
}

// Player representa o estado de um único jogador no servidor.
type Player struct {
	ID      string
	X, Y    int
	Simbolo rune
	Cor     Cor
}

// GameState representa o estado completo do jogo a ser sincronizado.
// Contém o mapa e o estado de todos os jogadores conectados.
type GameState struct {
	Mapa    [][]Elemento
	Players map[string]*Player
}

// Args representa os argumentos para uma chamada RPC de ação do cliente.
// Inclui um SequenceNumber para garantir a idempotência (exactly-once).
type Args struct {
	PlayerID       string
	Tecla          rune // Para movimentos ('w', 'a', 's', 'd') ou interação ('e').
	SequenceNumber int64
}

// ReplyConnect contém a resposta inicial do servidor quando um cliente se conecta.
type ReplyConnect struct {
	PlayerID     string
	InitialState GameState
}
