# Protoc

## install protoc for golang
brew install protoc

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"

echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
protoc-gen-go --version
protoc-gen-go-grpc --version

protoc --go-grpc_out=../grpcserv/proto system.proto
protoc --go_out=. --go-grpc_out=. system.proto

### for all proto files
find proto -name "*.proto" -exec protoc \
  --go_out=. \
  --go-grpc_out=. \
  {} \;