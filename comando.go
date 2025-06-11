package main

// Comando representa uma ação do jogador enviada ao servidor
type Comando struct {
	Tipo      TipoComando
	Jogador   string
	Direcao   int     // Para comandos de movimento
	Posicao   [2]int  // Para posições x,y
	Parametro string  // Tipo de comando em forma de string
	TeclaRune rune    // Tecla pressionada
}

// TipoComando define os tipos de comandos possíveis
type TipoComando int

const (
	CmdMover TipoComando = iota
	CmdAtacar
	CmdSair
	CmdEstado // Adicionar este comando se ainda não existir
	// Adicione outros tipos conforme necessário
)

// Resposta representa o resultado de um comando processado pelo servidor
type Resposta struct {
	Sucesso    bool
	Mensagem   string
	EstadoJogo interface{} // O estado atualizado do jogo ou parte dele
}
