.PHONY: build build-debug test test-with-cover test-in-docker clean clean-debug fmt

ROOT_PACKAGE=github.com/kos-v/dbunderfs
TESTS_PATH=./test
BINARY=dbfs
BINARY_DEBUG=dbfs-debug
RELEASE=`git describe --abbrev=0`
BUILD=`git rev-parse --short=8 HEAD`
BUILD_DATETIME=`date +%FT%H:%M:%S`

LDFLAGS=-w -s \
	-X main.fBinary=${BINARY} \
	-X main.fRelease=${RELEASE} \
	-X main.fBuild=${BUILD} \
	-X main.fBuildDatetime=${BUILD_DATETIME}

build: clean
	go build -ldflags "${LDFLAGS} -X main.fBinary=${BINARY} -X main.fDebug=false" -o ${BINARY} ${ROOT_PACKAGE}/cmd/dbfs

build-debug: clean-debug
	go build -ldflags "${LDFLAGS} -X main.fBinary=${BINARY_DEBUG} -X main.fDebug=true" -o ${BINARY_DEBUG} ${ROOT_PACKAGE}/cmd/dbfs

test:
	go test -v -race ${TESTS_PATH}/...

test-with-cover:
	go test -v -race -coverprofile=./coverage.txt -covermode=atomic -coverpkg=./internal/... ${TESTS_PATH}/...

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