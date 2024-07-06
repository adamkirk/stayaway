#!/bin/bash

# If a command was passed to the container run the binary 
if [ "$#" != "0" ]; then
    go run ./cmd/$APP_COMMAND "$@"
    exit "$?"
fi

DEBUG_OPT=""

if [ "$AIR_DEBUG" == "true" ]; then
    DEBUG_OPT="-d"
fi

air $DEBUG_OPT -build.bin ./build/stayaway-$APP_COMMAND -build.cmd "go build -o ./build/stayaway-$APP_COMMAND ./cmd/$APP_COMMAND"