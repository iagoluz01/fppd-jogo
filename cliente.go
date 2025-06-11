package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

// ClienteJogo representa um cliente do jogo
type ClienteJogo struct {
	nome      string
	clienteRPC *ClienteRPC
	ID        int
	Cor       termbox.Attribute
	// Adicione outros campos conforme necessário
}

// NovoCliente cria uma nova instância de ClienteJogo
func NovoCliente(nome string, endereco string, id int, cor termbox.Attribute) (*ClienteJogo, error) {
	rpcCliente, err := NovoClienteRPC(endereco)
	if err != nil {
		return nil, err
	}

	return &ClienteJogo{
		nome:      nome,
		clienteRPC: rpcCliente,
		ID:        id,
		Cor:       cor,
	}, nil
}

// Nome retorna o nome do cliente
func (c *ClienteJogo) Nome() string {
	return c.nome
}

// EnviarComando envia um comando para o servidor
func (c *ClienteJogo) EnviarComando(tipo string, tecla rune) (*Resposta, error) {
	cmd := &Comando{
		Jogador: c.nome,
		Parametro: tipo,
		TeclaRune: tecla,
	}
	
	var resposta Resposta
	err := c.clienteRPC.EnviarComando(cmd, &resposta)
	if err != nil {
		log.Printf("Erro ao enviar comando: %v", err)
		return nil, err
	}
	return &resposta, nil
}

// Sair encerra a conexão do cliente
func (c *ClienteJogo) Sair() {
	c.clienteRPC.Sair()
}
