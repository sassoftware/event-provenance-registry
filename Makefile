PACKAGE  = github.com/sassoftware/event-provenance-registry
BINARY   = bin/epr-server
COMMIT  ?= $(shell git rev-parse --short=16 HEAD)
gitversion := $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo 0.1.0-0)
VERSION ?= $(gitversion)
PREFIX = /usr/local

TOOLS    = $(CURDIR)/tools
PKGS     = $(or $(PKG),$(shell $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
GOLDFLAGS = "-X $(PACKAGE)/pkg/config.Version=$(VERSION) -X $(PACKAGE)/pkg/config.Commit=$(COMMIT) -s -w"

# Allow tags to be set on command-line, but don't set them
# by default
override TAGS := $(and $(TAGS),-tags $(TAGS))

GO      = go
GOBUILD = CGO_ENABLED=0 go build -v -tags netgo -installsuffix netgo
GOVET   = go vet
GODOC   = godoc
GOFMT   = gofmt
GOGENERATE = go generate
TIMEOUT = 15

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1mEvent Provenance Registry ▶\033[0m")

.PHONY: all
all: megalint test $(BINARY) $(BINARY)-arm64 $(BINARY)-darwin   ## Build all the binary types

.PHONY: linux
linux: test $(BINARY) ## build linux binary

.PHONY: darwin
darwin: test $(BINARY)-darwin  ## build darwin binary

.PHONY: darwin-arm64
darwin-arm64: test $(BINARY)-darwin-arm64 ## build darwin aarch64 binary

SOURCES = $(shell find -name vendor -prune -o -name \*.go -print)

$(BINARY): $(SOURCES); $(info $(M) building linux executable...) @ ## Build program binary for amd64
	$Q GOOS=linux GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-arm64: $(SOURCES); $(info $(M) building arm64 executable...) @ ## Build program binary for arm64
	$Q GOOS=linux GOARCH=arm64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-darwin: $(SOURCES); $(info $(M) building darwin executable...) @ ## Build program binary for darwin
	$Q GOOS=darwin GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-darwin-arm64: $(SOURCES); $(info $(M) building darwin executable...) @ ## Build program binary for darwin aarch64 
	$Q GOOS=darwin GOARCH=arm64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

.PHONY: docker-image
docker-image:;$(info $(M) running docker build...) @ ## Builds a local docker image tagged "epr-server:local"
	$Q docker build -t epr-server:local -f Dockerfile .

.PHONY: install
install: linux ;$(info $(M) installing epr-server...) @ ## Installs epr-server binary into DESTDIR
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./bin/epr-server $(DESTDIR)$(PREFIX)/bin/epr-server

.PHONY: install-darwin
install-darwin: darwin ;$(info $(M) installing epr-server...) @ ## Installs epr-server-darwin binary into DESTDIR epr-server
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./bin/epr-server-darwin $(DESTDIR)$(PREFIX)/bin/epr-server

.PHONY: install-darwin-arm64
install-darwin: darwin-arm64 ;$(info $(M) installing epr-server arm64...) @ ## Installs epr-server-darwin-arm64 binary into DESTDIR epr-server
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./bin/epr-server-darwin-arm64 $(DESTDIR)$(PREFIX)/bin/epr-server

.PHONY: list-updates
list-updates: ;$(info $(M) listing available go library updates...) @ ## List available go library updates
	$Q	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

.PHONY: megalint
megalint: ; $(info $(M) running golangci-lint...) @ ## Runs golangci-lint with a lot of switches
	$Q golangci-lint run --config ./golangci-megalint-config.yaml

.PHONY: test
test: ; $(info $(M) running tests...) @ ## Runs go test ./...
	$Q go test -v ./...

.PHONY: tidy
tidy: ; $(info $(M) running go mod tidy...) @ ## Run go mod tidy
	$Q go mod tidy

.PHONY: clean
clean: ; $(info $(M) cleaning...)	@ ## Cleanup everything
	@rm -rvf bin tools vendor build
	@rm -rvf tests/tests.* tests/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
