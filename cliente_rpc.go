package main

import (
	"fmt"
	"net/rpc"
)

// ClienteRPC representa a conexão RPC com o servidor
type ClienteRPC struct {
	cliente *rpc.Client
	addr    string
}

// NovoClienteRPC cria uma nova conexão com o servidor
func NovoClienteRPC(endereco string) (*ClienteRPC, error) {
	cliente, err := rpc.Dial("tcp", endereco)
	if err != nil {
		return nil, err
	}

	return &ClienteRPC{
		cliente: cliente,
		addr:    endereco,
	}, nil
}

// EnviarComando envia um comando para o servidor via RPC
func (c *ClienteRPC) EnviarComando(cmd *Comando, resposta *Resposta) error {
	return c.cliente.Call("ServidorJogo.ProcessarComando", cmd, resposta)
}

// Estado obtém o estado atual do jogo do servidor
func (c *ClienteRPC) Estado() (*EstadoJogo, error) {
	cmd := &Comando{
		Tipo: CmdEstado,
	}

	var resposta Resposta
	err := c.cliente.Call("ServidorJogo.ObterEstado", cmd, &resposta)
	if err != nil {
		return nil, err
	}

	// Verificar se a resposta foi bem-sucedida
	if !resposta.Sucesso {
		return nil, fmt.Errorf("erro no servidor: %s", resposta.Mensagem)
	}

	// Tentar converter o EstadoJogo da resposta
	estado, ok := resposta.EstadoJogo.(*EstadoJogo)
	if !ok {
		// Se não conseguir converter, criar um estado vazio
		return &EstadoJogo{}, nil
	}

	return estado, nil
}

// Sair fecha a conexão com o servidor
func (c *ClienteRPC) Sair() {
	if c.cliente != nil {
		c.cliente.Close()
	}
}

// Esta é a referência global que parece ser usada em jogo.go
var clienteRPC *ClienteRPC

// Remova a função init() que está causando a conexão automática
// e substitua por uma função de inicialização que pode ser chamada
// com o endereço correto
func InicializarClienteRPC(endereco string) error {
	var err error
	clienteRPC, err = NovoClienteRPC(endereco)
	if err != nil {
		return fmt.Errorf("erro ao conectar no servidor: %v", err)
	}
	return nil
}
