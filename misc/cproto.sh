#!/bin/sh

if [ ! -d "../proto" ]; then
  echo "Folder does not exist."
  exit 1
fi
if [ ! -d "../rpclient/proto" ]; then
  mkdir -p ../rpclient/proto
fi
cp ../proto/*.proto ../rpclient/proto
cp ../proto/build.rs  ../rpclient/build.rs

if [ ! -d "../rpcserv/proto" ]; then
  mkdir -p ../rpcserv/proto
fi
cp ../proto/*.proto ../rpcserv/proto
cp ../proto/build.rs  ../rpclient/build.rs

if [ ! -d "../yewapp/proto" ]; then
  mkdir -p ../yewapp/proto
fi
cp ../proto/*.proto ../yewapp/proto
cp ../proto/build.rs  ../rpclient/build.rs

if [ ! -d "../restgw/proto" ]; then
  mkdir -p ../restgw/proto
fi
cp ../proto/*.proto ../restgw/proto
cp ../proto/build.rs  ../rpclient/build.rs
