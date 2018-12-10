SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

TARGETS = linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-386 windows-amd64

BINARY=ixkit

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "" 
# "-X github.com/icoco/ixkit-cli/core.Version=${VERSION} -X github.com/icoco/ixkit-cli/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} console.go main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi