#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    cd  $ROOT_DIR
    $ROOT_DIR/hack/swag-init.sh

    go build -o ./build/stayaway-$APP_COMMAND ./cmd/main.go
)