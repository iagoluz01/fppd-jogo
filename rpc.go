package main

// Estruturas para RPC (Remote Procedure Call)

// Args para requisição de um jogador se conectar ao jogo
type EntrarArgs struct {
	Nome    string
	Simbolo rune
	Cor     Cor
}

// Resposta do servidor para um jogador que deseja se conectar
type EntrarReply struct {
	JogadorID int
	Sucesso   bool
	Mensagem  string
	Estado    EstadoJogo
}

// Args para enviar um comando ao servidor
type EnviarComandoArgs struct {
	JogadorID int
	Tipo      string // "mover" ou "interagir"
	Tecla     rune   // Para comandos de movimento
}

// Resposta do servidor para um comando enviado
type EnviarComandoReply struct {
	Sucesso  bool
	Mensagem string
}

// Args para obter o estado atual do jogo
type ObterEstadoArgs struct {
	JogadorID int
}

// Resposta do servidor com o estado atual do jogo
type ObterEstadoReply struct {
	Estado   EstadoJogo
	Sucesso  bool
	Mensagem string
}

// Args para um jogador sair do jogo
type SairArgs struct {
	JogadorID int
}

// Resposta do servidor para um jogador que deseja sair
type SairReply struct {
	Sucesso  bool
	Mensagem string
}
