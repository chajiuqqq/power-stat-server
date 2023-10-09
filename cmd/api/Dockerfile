# 使用golang作为基础镜像
FROM golang:1.21

ENV BROKER_IP=127.0.0.1
ENV PORT=1883
ENV PROFILE=dev
ENV CRON="0 0 1 * * *"
ENV WX_TOKEN=""
ENV WX_RECEIVER_ID=""
ENV WX_ENCODING_AES_KEY=""
# 将源代码复制到容器中
COPY ../.. /app/

# 设置工作目录
WORKDIR /app/cmd/api

# 构建应用程序
RUN go build -o api

# 设置容器启动命令
CMD ["./api"]