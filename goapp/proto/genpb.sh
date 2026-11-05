#!/bin/sh

set -eu

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
cd "$script_dir"

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. \
    *.proto
clientdir="../grpclient/pb"
serverdir="../grpcserv/pb"
goappdir="../goapp/pb"

dirlist="$clientdir $serverdir $goappdir"
# Loop through elements
for item in $dirlist; do
    echo "cp *.go $item"
    if [ ! -d "$item" ]; then
        mkdir -p "$item"
    fi
    cp *.go "$item"
done
#rm *.go