#!/bin/bash

set -eu

cd "$(dirname "${0}")"

rm -rf dist && mkdir dist

build() {
  rm -rf build && mkdir build
  GOOS=${1} GOARCH=${2} go build -o build/gitdump ./cmd/gitdump
  tar -czvf "dist/gitdump-${1}-${2}.tar.gz" --exclude ".*" -C build gitdump
  rm -rf build
}

build linux amd64
build windows amd64
build darwin amd64
build darwin arm64

build linux arm64
build linux loong64
build windows 386
build freebsd amd64

cd dist

shasum -a 256 *.tar.gz > SHASUM256.txt
gpg-trezor-hi@guoyk.xyz -ab SHASUM256.txt
