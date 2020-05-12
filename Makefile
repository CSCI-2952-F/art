.PHONY: all

all:
	mkdir -p bin
	cd src/client && go mod tidy && go mod vendor && go build -o ../../bin/client -mod vendor .