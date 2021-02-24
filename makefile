GO=go1.16.exe

all:
	$(GO) build -o gvm.exe main.go
	gvm

install:
	$(GO) install 
	gvm

use:
	$(GO) build -o gvm.exe main.go
	gvm use 1.15.3

list:
	$(GO) build -o gvm.exe main.go
	gvm list

listall:
	$(GO) build -o gvm.exe main.go
	gvm listall