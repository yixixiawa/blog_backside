# Dockerfile - 精简版本
FROM alpine:latest

WORKDIR /app

# 复制编译好的二进制文件
COPY sqlite_test-linux-amd64 /app/sqlite_test

# 复制配置文件和其他必要文件
COPY ./yaml ./yaml/
COPY ./data ./data/
COPY ./img ./img/

# 设置执行权限
RUN chmod +x /app/sqlite_test

# 暴露端口
EXPOSE 8080

# 运行程序
CMD ["./sqlite_test"]