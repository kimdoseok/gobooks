#!/bin/sh

set -eu

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
cd "$script_dir"

if [ ! -d "$script_dir/../third_party" ]; then
    mkdir -p "$script_dir/../third_party"
fi

googleapis_dir="$script_dir/../third_party/googleapis"
if [ ! -d "$googleapis_dir" ]; then
    git clone --depth 1 https://github.com/googleapis/googleapis.git "$googleapis_dir"
fi

protoc -I=. \
    -I $(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2 \
    -I=../third_party/googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
    *.proto


clientdir="../../grpclient/pb"
serverdir="../../grpcserv/pb"
goappdir="../../goapp/pb"

dirlist="$clientdir $serverdir $goappdir"
# Loop through elements
for item in $dirlist; do
    echo "cp *.go $item"
    if [ ! -d "$item" ]; then
        mkdir -p "$item"
    fi
    cp *.go "$item"
done
rm *.go