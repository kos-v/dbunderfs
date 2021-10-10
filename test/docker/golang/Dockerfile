ARG GOLANG_VERSION=1.16
FROM golang:${GOLANG_VERSION}

ARG HOST_USER_UID=1000
ARG HOST_USER_GID=1000
RUN set -ex                                                   && \
                                                                 \
    echo 'Creating tester user and group from host'           && \
    groupadd -g $HOST_USER_GID tester                         && \
    useradd -lm -u $HOST_USER_UID -g $HOST_USER_GID tester

WORKDIR /go/src/app