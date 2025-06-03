#!/bin/bash

set -e 

if [[ "$OSTYPE" == "linux-gnu" ]] then
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64

    sudo apt update
    sudo apt install libpcap-dev -y

else if [[ "$OSTYPE" == "darwin"* ]]; then
    export CGO_ENABLED=0
    export GOOS=darwin
    export GOARCH=amd64

    brew install libpcap

else
    echo "Unsupported OS"
    exit 1
fi

make all