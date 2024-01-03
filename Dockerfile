FROM golang:1.21 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /gitdump ./cmd/gitdump

FROM yankeguo/minit:1.13.0 AS minit

FROM alpine:3.18

# install packages
RUN apk add --no-cache tzdata ca-certificates git

# configure git
RUN git config --global --add safe.directory '*' && \
    git config --global init.defaultBranch main

# install minit
RUN mkdir -p /opt/bin
ENV PATH "/opt/bin:${PATH}"
COPY --from=minit /minit /opt/bin/minit
ENV MINIT_LOG_DIR none
ENTRYPOINT ["/opt/bin/minit"]

# install gitdump
COPY --from=builder /gitdump /gitdump
ENV MINIT_MAIN          /gitdump
ENV MINIT_MAIN_DIR      /data
ENV MINIT_MAIN_NAME     gitdump
ENV MINIT_MAIN_KIND     cron
ENV MINIT_MAIN_CRON     "@every 6h"

WORKDIR /data
