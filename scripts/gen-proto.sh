#!/usr/bin/env bash
set -e

# Add standard binary search paths for Go and Protoc plugins
export PATH="$PATH:/usr/local/bin:/usr/local/go/bin:/opt/homebrew/bin:$HOME/go/bin"

if command -v protoc &> /dev/null; then
    echo "Generating Protocol Buffer code using local protoc..."
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/inventory/v1/inventory.proto

    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/order/v1/order.proto
else
    echo "Local protoc not found. Running generation inside Docker golang container..."
    docker run --rm -v "$(pwd)":/workspace -w /workspace golang:alpine sh -c \
      "apk add --no-cache protobuf protobuf-dev > /dev/null 2>&1 && \
       go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 > /dev/null 2>&1 && \
       go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 > /dev/null 2>&1 && \
       protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/inventory/v1/inventory.proto && \
       protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/order/v1/order.proto"
fi

echo "Protobuf code generation completed successfully."
