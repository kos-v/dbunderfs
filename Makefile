SRC_PACKAGE="github.com/kos-v/dbunderfs/src"
TESTS_PATH=./tests
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

test: clean
	go test -v -race ./tests/...

test_with_cover: clean
	go test -v -race -coverprofile=./coverage.txt -covermode=atomic ${TESTS_PATH}/...

test-in-docker: prepare-docker-test-env
	docker-compose -f ${TESTS_PATH}/docker-compose.yml run app bash -c "make test_with_cover"

clean:
	rm -f ./${BINARY}
	rm -f ./coverage.txt

fmt:
	go fmt ./...

prepare-docker-test-env:
	rm -f ${TESTS_PATH}/.env
	cp ${TESTS_PATH}/.env.example ${TESTS_PATH}/.env
	echo ${GOLANG_VERSION}
	docker-compose -f ${TESTS_PATH}/docker-compose.yml up -d --build --force-recreate