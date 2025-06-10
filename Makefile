.PHONY: all build

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod
	go build -o jogo

server: build
	./jogo -servidor -porta=8080 -mapa=mapa.txt

client: build
	./jogo -endereco=localhost:8080 -nome="JogadorX"
	
clean:
	rm -f jogo

distclean: clean
	rm -f go.mod go.sum
