PACKAGE="github.com/kos-v/dbunderfs"
BINARY=dbfs
RELEASE=`git describe --abbrev=0`
BUILD=`git rev-parse --short=8 HEAD`
BUILD_DATETIME=`date +%FT%H:%M:%S`

LDFLAGS=-ldflags "-w -s -X ${PACKAGE}/cmd.binary=${BINARY} -X ${PACKAGE}/cmd.release=${RELEASE} -X ${PACKAGE}/cmd.build=${BUILD} -X ${PACKAGE}/cmd.buildDatetime=${BUILD_DATETIME}"

build: clean
	go build ${LDFLAGS} -o ${BINARY} main.go

clean:
	rm -f ./${BINARY}

fmt:
	go fmt ./...