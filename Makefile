.PHONY: all

all:
	mkdir -p bin
	go build -a -o bin/snappy cmd/main.go
	sudo mv bin/snappy /usr/local/bin/snappy