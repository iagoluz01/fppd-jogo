package main

import (
	"fmt"
	"net/rpc"
	"time"
)

// Cliente RPC para comunicação com o servidor
var clienteRPC *ClienteRPC

// ClienteRPC encapsula a comunicação RPC
type ClienteRPC struct {
	Client  *rpc.Client
	Estado  EstadoJogo
}

// NovoCliente estabelece uma conexão com o servidor
func NovoCliente(endereco, nome string, simbolo rune, cor Cor) (*ClienteJogo, error) {
	// Tentar estabelecer conexão RPC
	client, err := rpc.Dial("tcp", endereco)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor: %v", err)
	}

	// Criar cliente RPC global
	clienteRPC = &ClienteRPC{
		Client: client,
		Estado: EstadoJogo{
			Jogadores: make(map[int]JogadorInfo),
		},
	}

	// Tentar entrar no jogo
	args := EntrarArgs{
		Nome:    nome,
		Simbolo: simbolo,
		Cor:     cor,
	}
	reply := EntrarReply{}
	
	err = client.Call("ServidorJogo.Entrar", &args, &reply)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("erro ao entrar no jogo: %v", err)
	}

	if !reply.Sucesso {
		client.Close()
		return nil, fmt.Errorf("não foi possível entrar no jogo: %s", reply.Mensagem)
	}

	// Criar objeto de cliente local
	c := &ClienteJogo{
		ID:      reply.JogadorID,
		Nome:    nome,
		Simbolo: simbolo,
		Cor:     cor,
	}

	// Atualizar estado
	clienteRPC.Estado = reply.Estado
	
	// Iniciar goroutine para atualizações periódicas
	go atualizarEstadoPeriodicamente()

	return c, nil
}

// EnviarComando envia um comando para o servidor
func (c *ClienteJogo) EnviarComando(tipo string, tecla rune) error {
	if clienteRPC == nil || clienteRPC.Client == nil {
		return fmt.Errorf("cliente não está conectado")
	}

	args := EnviarComandoArgs{
		JogadorID: c.ID,
		Tipo:      tipo,
		Tecla:     tecla,
	}
	reply := EnviarComandoReply{}

	err := clienteRPC.Client.Call("ServidorJogo.EnviarComando", &args, &reply)
	if err != nil {
		return fmt.Errorf("erro ao enviar comando: %v", err)
	}

	if !reply.Sucesso {
		return fmt.Errorf("erro no servidor: %s", reply.Mensagem)
	}

	return nil
}

// Sair desconecta o cliente do servidor
func (c *ClienteJogo) Sair() error {
	if clienteRPC == nil || clienteRPC.Client == nil {
		return nil
	}

	args := SairArgs{
		JogadorID: c.ID,
	}
	reply := SairReply{}

	err := clienteRPC.Client.Call("ServidorJogo.Sair", &args, &reply)
	if err != nil {
		fmt.Printf("Aviso: erro ao sair do servidor: %v\n", err)
	}
	
	// Fechar conexão
	clienteRPC.Client.Close()
	clienteRPC = nil
	
	return nil
}

// atualizarEstadoPeriodicamente obtém o estado do jogo a cada intervalo
func atualizarEstadoPeriodicamente() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if clienteRPC == nil || clienteRPC.Client == nil {
			return
		}

		args := ObterEstadoArgs{
			JogadorID: clienteRPC.Estado.Jogadores[0].ID, // Usar primeiro jogador como ID
		}
		reply := ObterEstadoReply{}

		err := clienteRPC.Client.Call("ServidorJogo.ObterEstado", &args, &reply)
		if err != nil {
			fmt.Printf("Erro ao atualizar estado: %v\n", err)
			continue
		}

		if reply.Sucesso {
			clienteRPC.Estado = reply.Estado
		}
	}
}
