all: help

BIN := bin/
GO_LDFLAGS = -ldflags "-s -w"

ALWAYS:

server-build: ALWAYS
	go build $(GO_LDFLAGS) -o bin/server.bc cmd/server.go

server-run: ALWAYS
	./bin/server.bc $(ARGS)

wallet-build: ALWAYS
	go build $(GO_LDFLAGS) -o bin/apiWallet.bc cmd/apiWallet.go

wallet-run: ALWAYS
	./bin/apiWallet.bc $(ARGS)

clean: ALWAYS
	rm -f $(BIN)*.bc

help:
	@echo '	Usage:'
	@echo '	  server-build	Compile blockchain server.'
	@echo '	  server-run	Start blockchain server on port 5000 by default. To start on other port (e.g. 4000) execute:'
	@echo '				  make server-run ARGS="-port 4000"'
	@echo '	  wallet-build	Compile blockchain wallet.'
	@echo '	  wallet-run	Start wallet server on port 8080 and connecting to blockchain node on http://127.0.0.1:5000 by default.'
	@echo '			To start with other values execute:'
	@echo '				  make wallet-run ARGS="-port 8000 -gateway=http://node.corp:5003"'

author:
	@echo "Project by feliux"
	@echo "https://github.com/feliux"
