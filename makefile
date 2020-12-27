BIN := go-mml.exe
SRC := $(wildcard *.go)

default: build test

build: $(BIN)

test:
	go install
	$(BIN) \
		-s \
		-m $(USERPROFILE)\tools\nsf\ppmck09a\mck \
		-n $(USERPROFILE)\tools\nsf\nsf2wav\nsf2wav \
		-f .\test\sample_auto_bank.mml

setup:
	go get -u golang.org/x/text/encoding/japanese
	go get -u golang.org/x/text/transform
	go get -u github.com/faiface/beep
	go get -u github.com/pkg/errors
	go get -u github.com/hajimehoshi/oto

$(BIN): $(SRC)
	go build -o $@ $^

.PHONY: default build test setup
