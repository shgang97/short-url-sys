FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /workspace

# 只复制当前服务的代码和 shared 依赖
COPY generate-service/ ./generate-service/
COPY shared/ ./shared/

# 进入服务目录构建
WORKDIR /workspace/generate-service

# 设置 Go 代理（可选，国内网络推荐）
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖并构建
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o generate-service ./cmd/

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder --chown=appuser:appgroup /workspace/generate-service/generate-service .
# 设置环境变量 - 默认使用容器内的配置文件
ENV CONFIG_PATH=configs/config.dev.yaml
USER appuser
EXPOSE 8080
CMD ["./generate-service"]