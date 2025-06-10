// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao     Cor = termbox.ColorDefault
	CorCinzaEscuro    = termbox.ColorDarkGray
	CorVermelho       = termbox.ColorRed
	CorVerde          = termbox.ColorGreen
	CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede    = termbox.ColorDarkGray
	CorTexto          = termbox.ColorDarkGray
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// Renderiza o estado do jogo multiplayer
func interfaceDesenharJogoMultiplayer(jogo *Jogo) {
	// Atualiza o estado local com o estado do servidor
	jogoAtualizarEstadoMultiplayer(jogo)
	
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem local sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Elemento{
		Simbolo:  '☺',
		Cor:      jogo.Cliente.Cor,
		CorFundo: CorPadrao,
		Tangivel: true,
	})
	
	// Desenha os outros jogadores
	for _, jogador := range jogo.OutrosJogadores {
		interfaceDesenharElemento(jogador.PosX, jogador.PosY, Elemento{
			Simbolo:  jogador.Simbolo,
			Cor:      jogador.Cor,
			CorFundo: CorPadrao,
			Tangivel: true,
		})
	}

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)
	
	// Desenha informações de jogadores conectados
	interfaceDesenharInfoJogadores(jogo)

	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// Exibe informações sobre os jogadores conectados
func interfaceDesenharInfoJogadores(jogo *Jogo) {
	if jogo.Cliente == nil {
		return
	}
	
	// Mensagem com nome do jogador local
	msgLocal := fmt.Sprintf("Você: %s", jogo.Cliente.Nome)
	for i, c := range msgLocal {
		termbox.SetCell(i, len(jogo.Mapa)+5, c, jogo.Cliente.Cor, CorPadrao)
	}
	
	// Listagem de outros jogadores
	msgOutros := "Outros jogadores: "
	offset := len(msgOutros)
	for i, c := range msgOutros {
		termbox.SetCell(i, len(jogo.Mapa)+6, c, CorTexto, CorPadrao)
	}
	
	// Nomes dos outros jogadores
	linha := len(jogo.Mapa) + 6
	coluna := offset
	for _, jogador := range jogo.OutrosJogadores {
		info := fmt.Sprintf("%s ", jogador.Nome)
		for _, c := range info {
			termbox.SetCell(coluna, linha, c, jogador.Cor, CorPadrao)
			coluna++
		}
		
		// Avança para próxima linha se estiver muito longo
		if coluna > 70 {
			linha++
			coluna = offset
		}
	}
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.Simbolo, elem.Cor, elem.CorFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}

