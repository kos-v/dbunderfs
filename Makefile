SRC_PACKAGE="github.com/kos-v/dbunderfs/src"
BINARY=dbfs
RELEASE=`git describe --abbrev=0`
BUILD=`git rev-parse --short=8 HEAD`
BUILD_DATETIME=`date +%FT%H:%M:%S`

LDFLAGS=-w -s \
	-X ${SRC_PACKAGE}/cmd.fBinary=${BINARY} \
	-X ${SRC_PACKAGE}/cmd.fRelease=${RELEASE} \
	-X ${SRC_PACKAGE}/cmd.fBuild=${BUILD} \
	-X ${SRC_PACKAGE}/cmd.fBuildDatetime=${BUILD_DATETIME}

build: clean
	go build -ldflags "${LDFLAGS} -X ${SRC_PACKAGE}/cmd.fDebug=false" -o ${BINARY} main.go

build_debug: clean
	go build -ldflags "${LDFLAGS} -X ${SRC_PACKAGE}/cmd.fDebug=true" -o ${BINARY} main.go

clean:
	rm -f ./${BINARY}

test:
	go test -v -race ./tests/...

test_with_cover:
	rm -f ./coverage.txt
	go test -v -race -coverprofile=./tests/coverage.txt -covermode=atomic ./tests/...

fmt:
	go fmt ./...