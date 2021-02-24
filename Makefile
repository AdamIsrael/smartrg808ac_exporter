GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


all: test build

build:
	$(GOBUILD) -o smartrg808ac_exporter

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -cover ./...

clean:
	$(GOCLEAN)
	rm -f smartrg808ac_exporter
