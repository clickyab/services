export SERVICE_ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export GO?=$(shell which go)
export UPDATE?=

$(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(SERVICE_ROOT)/tmp/ip2l
	wget -c http://www.clickyab.com/downloads/IP-COUNTRY-REGION-CITY.BIN -O $(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN

service_bindata:
	GOPATH=$(SERVICE_ROOT)/tmp $(GO) get -v $(UPDATE) github.com/jteeuwen/go-bindata/go-bindata

$(SERVICE_ROOT)/ip2location/data.gen.go: $(SERVICE_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN service_bindata
	cd $(SERVICE_ROOT)/tmp/ip2l && $(SERVICE_ROOT)/tmp/bin/go-bindata -nomemcopy -o $(SERVICE_ROOT)/ip2location/data.gen.go -pkg ip2location .

all: $(SERVICE_ROOT)/ip2location/data.gen.go
