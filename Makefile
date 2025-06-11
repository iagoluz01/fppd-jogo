.PHONY: all build

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod
	go build -o jogo

server: build
	./jogo -servidor -porta=8080 -mapa=mapa.txt -endereco=0.0.0.0:8080

# Variável para o endereço do servidor, que pode ser sobrescrita ao chamar o make
SERVER_ADDR ?= localhost:8080
# Nome do jogador
PLAYER_NAME ?= JogadorX

client: build
	./jogo -endereco=$(SERVER_ADDR) -nome="$(PLAYER_NAME)"
	
clean:
	rm -f jogo

distclean: clean
	rm -f go.mod go.sum
