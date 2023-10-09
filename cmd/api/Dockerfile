# 使用golang作为基础镜像
FROM golang:1.21

# 将源代码复制到容器中
COPY ../.. /app/

# 设置工作目录
WORKDIR /app/cmd/api

# 构建应用程序
RUN go build -o api

# 设置容器启动命令
CMD ["./api"]