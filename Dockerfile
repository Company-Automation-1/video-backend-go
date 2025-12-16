# 第一阶段：构建
FROM golang:1.25-alpine AS builder

WORKDIR /build

# 安装必要的工具
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（编译为Linux二进制）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -ldflags="-w -s" .

# 第二阶段：运行
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/app .

# 复制配置文件
COPY config.yaml .

# 暴露端口
EXPOSE 8888

# 运行应用
CMD ["./app"]