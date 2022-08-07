ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR:=${CURDIR}/..
export GOBIN :=${BIN_DIR}/bin

all: manager
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
	GOBIN=$(shell go env GOPATH)/bin
else
	GOBIN=$(shell go env GOBIN)
endif
# Build manager binary
manager: fmt clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=auto go install -ldflags="-s -w"

# Run go fmt against code
fmt: tidy
	go fmt ./...

# Run go vet against code
tidy:
	go mod tidy
	go mod download

clean:
	go clean
	rm -f ${GOBIN}/*