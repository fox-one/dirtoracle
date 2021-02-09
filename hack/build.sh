#!/usr/bin/env bash

set -e

TAG=${1:-dev}
VERSION=$(git describe --tags)
COMMIT=$(git rev-parse --short HEAD)
BINARY=dirtoracle."${TAG}"

CONFIG=config."${TAG}".yaml
if [ -f "${CONFIG}" ]; then
  trap 'rm -f config_gen.go' EXIT
  # go get -u github.com/fox-one/pkg/config/config-gen
  echo "use config ${CONFIG}"
  config-gen --config "${CONFIG}" --tag "${TAG}"
fi

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "build ${BINARY} with version ${VERSION} & commit ${COMMIT}"
go build -a -installsuffix cgo \
         --tags "${TAG}" \
         --ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
         -o "${BINARY}"
