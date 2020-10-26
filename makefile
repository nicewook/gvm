GO=go1.15.3.exe

all:
	$(GO) mod tidy
	$(GO) mod vendor 
	$(GO) build -o gvm.exe main.go
	gvm

use:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) build -o gvm.exe main.go
	gvm use 1.15.3

list:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) build -o gvm.exe main.go
	gvm list

listall:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) build -o gvm.exe main.go
	gvm listall