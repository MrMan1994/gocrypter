BIN=$(CURDIR)/bin
BUILD=$(CURDIR)/build
GOLINT=$(BIN)/golint
EXECUTABLE=$(BUILD)/main

GO=go
GOBUILD=$(GO) build -v -x
GOTEST=$(GO) test -v
GOINSTALL=$(GO) install -v
GOGET=$(GO) get -v
GOFMT=gofmt
GOCLEAN=$(GO) clean
SRCFILES=main.go go.mod log hash cmd

all: $(EXECUTABLE)

$(BIN):
	@mkdir -p "$@"

$(BIN)/%: | $(BIN)
	@env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) go get -v $(PACKAGE) || ret=$$?; rm -rf $$tmp ; exit $$ret

$(BUILD):
	@mkdir -p "$@"
	@cp -r $(SRCFILES) "$@"

$(BUILD)/%: | $(BUILD) $(BIN)
	@sh -c "cd $(BUILD) && $(GOBUILD) -o main main.go && cp main $(BIN) && cp main .."

clean:
	@$(GOCLEAN)
	@rm -rf $(BUILD) $(BIN)

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint
$(BUILD)/main: GOOS=linux

lint: | $(GOLINT)
	@$(GOLINT) ./...

get:
	@$(GOGET)

fmt:
	@$(GOFMT) ./...