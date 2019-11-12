BIN=$(CURDIR)/bin
BUILD=$(CURDIR)
GOLINT=$(BIN)/golint
PROG=$(CURDIR)/gocrypter


GO=go
GOBUILD=$(GO) build -v -x
GOTEST=$(GO) test -v
GOINSTALL=$(GO) install -v -x
GOGET=$(GO) get -v
GOFMT=$(GO) fmt
GOCLEAN=$(GO) clean

all: $(BUILD)/gocrypter

$(BIN):
	@mkdir -p "$@"

$(BIN)/%: | $(BIN)
	@env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) go get -v $(PACKAGE) || ret=$$?; rm -rf $$tmp ; exit $$ret

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

$(BUILD):
	@mkdir -p "$@"

$(BUILD)/%: | $(BUILD) $(BIN)
	@sh -c "cd $(BUILD) && $(GOBUILD) -o $(PROG) $(SRCFILE)"

$(BUILD)/gocrypter: SRCFILE=gocrypter.go

clean:
	@$(GOCLEAN)
	@rm -rf $(BIN)



lint: | $(GOLINT)
	@$(GOLINT) ./...

get:
	@$(GOGET)

fmt:
	@$(GOFMT) ./...

install:
	$(GOINSTALL) gocrypter.go

completion:
ifeq ("$$SHELL", "/bin/bash")
	@echo "bash"
else ifeq ("$$SHELL", "zsh")
	@echo "zsh"
else
	@echo "no shell support for: $$SHELL"
endif
