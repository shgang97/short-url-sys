#!/bin/bash

set -e

echo "🔧 Generating Protobuf code..."

# 创建输出目录
mkdir -p shared/proto/generate

# 生成生成服务的代码
protoc --go_out=shared/proto/generate \
       --go-grpc_out=shared/proto/generate \
       -Ishared/proto \
       shared/proto/generate/generate_service.proto

echo "✅ Protobuf code generated successfully!"