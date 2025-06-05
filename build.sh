#!/bin/bash

set -e 

command_exist() {
    command -v "$1" > /dev/null 2>&1
}

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64

    if ! dpkg -l | grep -q libpcap-dev; then
        sudo apt update
        sudo apt install libpcap-dev -y
    fi

elif [[ "$OSTYPE" == "darwin"* ]]; then
    export CGO_ENABLED=0
    export GOOS=darwin
    export GOARCH=amd64

    if ! brew list | grep -q libpcap; then
        brew install libpcap
    fi

else
    echo "WARNNING: Unsupported OS: $OSTYPE. We still working to support this OS."
    exit 1
fi

##make all