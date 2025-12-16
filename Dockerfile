# 第一阶段：构建阶段
FROM golang:1.21-alpine AS builder

# 安装 gcc 和 musl-dev，这是编译 CGO (sqlite3) 所必需的
RUN apk add --no-cache gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
# 这样做可以利用 Docker 缓存层，如果依赖没变，就不需要重新下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
# CGO_ENABLED=0 表示禁用 CGO，生成纯静态二进制文件，适合在 scratch 或 alpine 镜像中运行
# -o main 指定输出文件名
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装 ca-certificates，以便应用可以进行 HTTPS 请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 如果你的应用需要配置文件（例如 config.yaml），也需要复制进去
# COPY --from=builder /app/config.yaml .

# 暴露应用端口（根据你的应用实际端口修改，例如 8080）
EXPOSE 8080

# 运行应用
CMD ["./main"]