#!/bin/bash

set -eu

cd $(dirname $0)

rm -rf dist
mkdir -p dist

build() {
  export GOOS=$1
  export GOARCH=$2
  export CGO_ENABLED=0
  FILENAME="gitdump-${GOOS}-${GOARCH}"
  go build -o "dist/${FILENAME}" ./cmd/gitdump
  cd dist
  tar czvf "${FILENAME}.tar.gz" "${FILENAME}"
  rm -f "${FILENAME}"
  cd ..
}

build linux amd64
build windows amd64
build darwin amd64
build darwin arm64

build linux arm64
build linux loong64
build windows 386
build freebsd amd64
