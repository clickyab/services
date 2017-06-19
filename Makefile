export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO=$(shell which go)
export GIT=$(shell which git)
export BIN=$(ROOT)/bin
export GOPATH=$(abspath $(ROOT)/../../../..)
export GOBIN?=$(BIN)
export DIFF=$(shell which diff)
export LINTER=$(BIN)/gometalinter.v1
# TODO : Ignoring services/codegen is a bad thing. try to get it back to lint
export LINTERCMD=$(LINTER) -e ".*.gen.go" -e ".*_test.go" -e "codegen/.*" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell
metalinter:
	$(GO) get -v gopkg.in/alecthomas/gometalinter.v1
	$(GO) install -v gopkg.in/alecthomas/gometalinter.v1
	$(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make -f $(ROOT)/Makefile metalinter

dependencies:
	cd $(ROOT) && $(GO) get -t -v ./...

test: dependencies
	cd $(ROOT) && $(GO) test ./...

all: dependencies
	cd $(ROOT) && $(GO) build ./...

lint: dependencies $(LINTER)
	cd $(ROOT) && $(LINTERCMD) $(ROOT)/...
