BIN=$(CURDIR)/bin
BUILD=$(CURDIR)
GOLINT=$(BIN)/golint
PROG=$(CURDIR)/gocrypter

SHELL:=/bin/bash

GO=go
GOBUILD=$(GO) build -v -x
GOTEST=$(GO) test -v
GOINSTALL=$(GO) install -v -x
GOGET=$(GO) get -v
GOFMT=$(GO) fmt
GOCLEAN=$(GO) clean
GOBIN=$$GOBIN

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
	@env $(GOCLEAN)
	@rm -rf $(BIN)

lint: | $(GOLINT)
	@env $(GOLINT) ./...

get:
	@env $(GOGET)

fmt:
	@env $(GOFMT) ./...

install: completion
	@env GOBIN=$(GOBIN) $(GOINSTALL) gocrypter.go

completion:
ifeq "$(shell basename $$SHELL)" "bash"
	$(info installing command line completion for bash...)
	@if [ $(shell cat ~/.bashrc | grep "source <(gocrypter completion bash)" | wc -l) -ne 0 ]; then \
	echo "bash completion has already been installed" ; \
	else echo 'source <(gocrypter completion bash)' >> ~/.bashrc && echo "Please reload your shell!" ; \
	fi;
else ifeq "$(shell basename $$SHELL)" "zsh"
	$(info installing command line completion for zsh...)
	@if [ $(shell cat ~/.zshrc | grep "source <(gocrypter completion zsh)" | wc -l) -ne 0 ]; then \
	echo "zsh completion has already been installed" ; \
	else echo 'source <(gocrypter completion zsh)' >> ~/.zshrc && echo "Please reload your shell!" ; \
	fi;
else
	@echo WARNING: cannot install command line completion for: \"$(shell echo $$SHELL)\": unsupported shell
endif

remove-completion:
ifeq "$(shell basename $$SHELL)" "bash"
	$(info removing command line completion for bash...)
	@if [ $(shell cat ~/.bashrc | grep "source <(gocrypter completion bash)" | wc -l) -gt 0 ]; then \
	sed -i '/source <(gocrypter completion bash)/ d' ~/.bashrc && echo "removed bash completion" ; \
	else echo "bash completion is not installed" ; \
	fi;
else ifeq "$(shell basename $$SHELL)" "zsh"
	$(info removing command line completion for zsh...)
	@if [ $(shell cat ~/.bashrc | grep "source <(gocrypter completion bash)" | wc -l) -gt 0 ]; then \
	sed -i '/source <(gocrypter completion zsh)/ d' ~/.zshrc && echo "removed zsh completion" ; \
	else echo "zsh completion is not installed" ; \
	fi;
else
	@echo WARNING: cannot remove command line completion for: \"$(shell echo $$SHELL)\": unsupported shell
endif