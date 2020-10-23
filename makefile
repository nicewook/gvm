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