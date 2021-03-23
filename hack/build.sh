#!/usr/bin/env bash

set -e

TAG=${1:-dev}
VERSION=$(git describe --tags --abbrev=0)
COMMIT=$(git rev-parse --short HEAD)
BINARY=dirtoracle."${TAG}"

CONFIG=config."${TAG}".yaml
if [ -f "${CONFIG}" ]; then
  trap 'rm -f config_gen.go' EXIT
  if ! type "config-gen" > /dev/null 2>/dev/null; then
    env GO111MODULE=off go get -u github.com/fox-one/pkg/config/config-gen
  fi
  echo "use config ${CONFIG}"
  config-gen --config "${CONFIG}" --tag "${TAG}"
fi

echo "build ${BINARY} with version ${VERSION} & commit ${COMMIT}"
CGO_ENABLED=1 go build -a \
         --tags "${TAG}" \
         --ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
         -o builds/

trap 'rm -f config_gen.go' EXIT
