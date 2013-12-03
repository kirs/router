export GOPATH := $(CURDIR)

all: build

.PHONY: build

build: install-dependencies bin/router bin/server

clean:
	-rm -rf bin/

bin/router: router/main.go router/views
	go build -o $@ router/*.go
	cp -R router/views bin/

bin/server: server/main.go
	go build -o $@ server/*.go

install-dependencies:
	go get github.com/BurntSushi/toml
	go get github.com/howeyc/fsnotify
