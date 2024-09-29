#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    cd $ROOT_DIR
    swag init -o ./internal/openapidoc -g ./internal/api/server.go
)