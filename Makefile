.PHONY: build build-debug test test-with-cover test-in-docker clean clean-debug fmt

SRC_PACKAGE="github.com/kos-v/dbunderfs/src"
TESTS_PATH=./tests
BINARY=dbfs
BINARY_DEBUG=dbfs-debug
RELEASE=`git describe --abbrev=0`
BUILD=`git rev-parse --short=8 HEAD`
BUILD_DATETIME=`date +%FT%H:%M:%S`

LDFLAGS=-w -s \
	-X ${SRC_PACKAGE}/cmd.fBinary=${BINARY} \
	-X ${SRC_PACKAGE}/cmd.fRelease=${RELEASE} \
	-X ${SRC_PACKAGE}/cmd.fBuild=${BUILD} \
	-X ${SRC_PACKAGE}/cmd.fBuildDatetime=${BUILD_DATETIME}

build: clean
	go build -ldflags "${LDFLAGS} -X ${SRC_PACKAGE}/cmd.fBinary=${BINARY} -X ${SRC_PACKAGE}/cmd.fDebug=false" -o ${BINARY} main.go

build-debug: clean-debug
	go build -ldflags "${LDFLAGS} -X ${SRC_PACKAGE}/cmd.fBinary=${BINARY_DEBUG} -X ${SRC_PACKAGE}/cmd.fDebug=true" -o ${BINARY_DEBUG} main.go

test:
	go test -v -race ${TESTS_PATH}/...

test-with-cover:
	go test -v -race -coverprofile=./coverage.txt -covermode=atomic -coverpkg=${SRC_PACKAGE}/... ${TESTS_PATH}/...

test-in-docker:
	rm -f ${TESTS_PATH}/.env
	cp ${TESTS_PATH}/.env.example ${TESTS_PATH}/.env
	docker-compose -f ${TESTS_PATH}/docker-compose.yml up -d --build --force-recreate
	docker-compose -f ${TESTS_PATH}/docker-compose.yml run app bash -c "make test-with-cover"

clean:
	rm -f ./${BINARY}

clean-debug:
	rm -f ./${BINARY_DEBUG}

fmt:
	go fmt ./...