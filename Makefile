PACKAGE  = gitlab.sas.com/async-event-infrastructure/server
BINARY   = bin/server
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
M = $(shell printf "\033[34;1mserver ▶\033[0m")

.PHONY: all
all: megalint test $(BINARY) $(BINARY)-arm64 $(BINARY)-darwin   ## Build all the binary types

.PHONY: linux
linux: test $(BINARY) ## build linux binary

.PHONY: darwin
darwin: test $(BINARY)-darwin  ## build darwin binary


SOURCES = $(shell find -name vendor -prune -o -name \*.go -print)

$(BINARY): $(SOURCES); $(info $(M) building linux executable...) @ ## Build program binary
	$Q GOOS=linux GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-arm64: $(SOURCES); $(info $(M) building arm64 executable...) @ ## Build program binary for arm64
	$Q GOOS=linux GOARCH=arm64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-darwin: $(SOURCES); $(info $(M) building darwin executable...) @ ## Build program binary
	$Q GOOS=darwin GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

GOIMPORTS = $(TOOLS)/goimports
$(GOIMPORTS): ; $(info $(M) building goimports...)
	$Q go build -o $@ golang.org/x/tools/cmd/goimports

GOCOVMERGE = $(TOOLS)/gocovmerge
$(GOCOVMERGE): ; $(info $(M) building gocovmerge...)
	$Q go build -o $@ github.com/wadey/gocovmerge

GOCOV = $(TOOLS)/gocov
$(GOCOV): ; $(info $(M) building gocov...)
	$Q go build -o $@ github.com/axw/gocov/gocov

GOCOVXML = $(TOOLS)/gocov-xml
$(GOCOVXML): ; $(info $(M) building gocov-xml...)
	$Q go build -o $@ github.com/AlekSi/gocov-xml

GO2XUNIT = $(TOOLS)/go2xunit
$(GO2XUNIT): ; $(info $(M) building go2xunit...)
	$Q go build -o $@ github.com/tebeka/go2xunit

GOBINDATA = $(TOOLS)/go-bindata
$(GOBINDATA): ; $(info $(M) building go-bindata...)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/go-bindata/go-bindata/v3/go-bindata

GOVERSIONINFO = $(TOOLS)/goversioninfo
$(GOVERSIONINFO): ; $(info $(M) building goversioninfo...)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/josephspurrier/goversioninfo/cmd/goversioninfo

$(TOOLS)/protoc-gen-go: ; $(info $(M) building protoc-gen-go...)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/golang/protobuf/protoc-gen-go

.PHONY: docker-image
docker-image:;$(info $(M) running docker build...) @ ## Builds a local docker image tagged "server:local"
	$Q docker build -t server:local -f Dockerfile .

.PHONY: install
install: linux ;$(info $(M) installing server...) @ ## Installs server binary into DESTDIR
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./bin/server $(DESTDIR)$(PREFIX)/bin/server

.PHONY: install-darwin
install-darwin: darwin ;$(info $(M) installing server...) @ ## Installs server-darwin binary into DESTDIR aserver
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./bin/server $(DESTDIR)$(PREFIX)/bin/server

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
