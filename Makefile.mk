export SERVICE_ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export GO?=$(shell which go)
export UPDATE?=
export GOPATH?=$(shell mktemp -d)
export LINTER=$(GOPATH)/bin/gometalinter.v1
export LINTERCMD?=$(LINTER) -e ".*.gen.go" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=misspell

$(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(SERVICE_ROOT)/tmp/ip2l
	wget -c http://www.clickyab.com/downloads/IP-COUNTRY-REGION-CITY.BIN -O $(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN

service_bindata:
	GOPATH=$(SERVICE_ROOT)/tmp $(GO) get -v $(UPDATE) github.com/jteeuwen/go-bindata/go-bindata

$(SERVICE_ROOT)/ip2location/data.gen.go: $(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN service_bindata
	cd $(SERVICE_ROOT)/tmp/ip2l && $(SERVICE_ROOT)/tmp/bin/go-bindata -nomemcopy -o $(SERVICE_ROOT)/ip2location/data.gen.go -pkg ip2location .

all: $(SERVICE_ROOT)/ip2location/data.gen.go

convey:
	$(GO) get -v github.com/smartystreets/goconvey/...

linter:
	$(GO) get -v gopkg.in/alecthomas/gometalinter.v1
	$(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make -f Makefile.mk linter

test: convey $(LINTER)
	rm -rf $(GOPATH)/src/services
	mkdir -p $(GOPATH)/src/ && cp -r $(SERVICE_ROOT) $(GOPATH)/src/
	$(LINTERCMD) $(GOPATH)/src/services/...
	$(GO) get -v github.com/smartystreets/goconvey/...
	cd $(GOPATH)/src/services && $(GO) get -v ./...
	cd $(GOPATH)/src/services && $(GO) test -v ./...
	rm -rf $(GOPATH)/src/services

