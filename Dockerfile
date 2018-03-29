FROM alpine

ARG SOURCE_REPOSITORY=github.com/Project-Prismatica/prismatica-idam

ENV GOPATH=/go

ADD . /go/src/${SOURCE_REPOSITORY}

RUN apk add --no-cache dumb-init git go musl-dev && \
    go get ${SOURCE_REPOSITORY}/go/prismatica-idam-server && \
    apk del git go musl-dev && \
    adduser -S -D -H prismatica-idam

USER prismatica-idam
ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]
CMD [ "/go/bin/prismatica-idam-server" ]
