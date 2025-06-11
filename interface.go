// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

// Definição de cores
type Cor termbox.Attribute

// Cores disponíveis
const (
	CorPadrao      Cor = Cor(termbox.ColorDefault)
	CorPreto       Cor = Cor(termbox.ColorBlack)
	CorVermelho    Cor = Cor(termbox.ColorRed)
	CorVerde       Cor = Cor(termbox.ColorGreen)
	CorAmarelo     Cor = Cor(termbox.ColorYellow)
	CorAzul        Cor = Cor(termbox.ColorBlue)
	CorMagenta     Cor = Cor(termbox.ColorMagenta)
	CorCiano       Cor = Cor(termbox.ColorCyan)
	CorBranco      Cor = Cor(termbox.ColorWhite)
	CorCinzaEscuro Cor = Cor(termbox.ColorDarkGray)
	CorParede      Cor = Cor(termbox.ColorGreen)
	CorFundoParede Cor = Cor(termbox.ColorBlack)
)

// Evento de teclado
type EventoTeclado struct {
	Tipo  string
	Tecla rune
}

// Inicializa a interface termbox
func interfaceIniciar() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	
	RegistrarLog("Interface termbox inicializada com sucesso")
}

// Finaliza a interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento de teclado e retorna o tipo de evento e a tecla pressionada
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	
	// Verifica se é um evento de tecla
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyEsc:
			return EventoTeclado{Tipo: "sair", Tecla: 0}
		case termbox.KeySpace:
			return EventoTeclado{Tipo: "interagir", Tecla: ' '}
		default:
			// Se for uma tecla normal (não especial)
			if ev.Ch != 0 {
				switch ev.Ch {
				case 'q', 'Q':
					return EventoTeclado{Tipo: "sair", Tecla: ev.Ch}
				case 'w', 'a', 's', 'd', 'W', 'A', 'S', 'D':
					return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
				case 'e', 'E':
					return EventoTeclado{Tipo: "interagir", Tecla: ev.Ch}
				default:
					return EventoTeclado{Tipo: "outro", Tecla: ev.Ch}
				}
			}
		}
	}
	
	// Para qualquer outro evento
	return EventoTeclado{Tipo: "outro", Tecla: 0}
}

// Desenha o jogo na interface
func interfaceDesenharJogoMultiplayer(jogo *Jogo) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	// Se o mapa não estiver carregado ainda
	if jogo.Mapa == nil || len(jogo.Mapa) == 0 {
		interfaceDesenharTexto(0, 0, "Aguardando informações do servidor...", termbox.ColorWhite, termbox.ColorDefault)
		interfaceDesenharTexto(0, 1, fmt.Sprintf("Jogador: %s", jogo.Cliente.Nome()), termbox.ColorWhite, termbox.ColorDefault)
		termbox.Flush()
		return
	}

	// Desenhando o mapa
	for y := 0; y < len(jogo.Mapa); y++ {
		for x := 0; x < len(jogo.Mapa[y]); x++ {
			elemento := jogo.Mapa[y][x]
			termbox.SetCell(x, y, elemento.Simbolo, termbox.Attribute(elemento.Cor), termbox.Attribute(elemento.CorFundo))
		}
	}

	// Desenhando outros jogadores
	for _, jogador := range jogo.OutrosJogadores {
		termbox.SetCell(jogador.PosX, jogador.PosY, jogador.Simbolo, termbox.Attribute(jogador.Cor), termbox.ColorDefault)
	}

	// Desenhando o personagem do jogador
	termbox.SetCell(jogo.PosX, jogo.PosY, '☺', termbox.ColorWhite, termbox.ColorDefault)

	// Desenhando barra de status
	statusY := len(jogo.Mapa) + 1
	interfaceDesenharTexto(0, statusY, jogo.StatusMsg, termbox.ColorWhite, termbox.ColorDefault)
	
	termbox.Flush()
}

// Desenha um texto na posição (x,y)
func interfaceDesenharTexto(x, y int, texto string, corFrente, corFundo termbox.Attribute) {
	for i, c := range texto {
		termbox.SetCell(x+i, y, c, corFrente, corFundo)
	}
}

