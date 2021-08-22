PACKAGE="github.com/kos-v/dbunderfs"
BINARY=dbfs
RELEASE=`git describe --abbrev=0`
BUILD=`git rev-parse --short=8 HEAD`
BUILD_DATETIME=`date +%FT%H:%M:%S`

LDFLAGS=-w -s \
	-X ${PACKAGE}/cmd.fBinary=${BINARY} \
	-X ${PACKAGE}/cmd.fRelease=${RELEASE} \
	-X ${PACKAGE}/cmd.fBuild=${BUILD} \
	-X ${PACKAGE}/cmd.fBuildDatetime=${BUILD_DATETIME}

build: clean
	go build -ldflags "${LDFLAGS} -X ${PACKAGE}/cmd.fDebug=false" -o ${BINARY} main.go

build_debug: clean
	go build -ldflags "${LDFLAGS} -X ${PACKAGE}/cmd.fDebug=true" -o ${BINARY} main.go

clean:
	rm -f ./${BINARY}

test:
	go test -v -race ./...

test_with_cover:
	rm -f ./coverage.txt
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	go fmt ./...