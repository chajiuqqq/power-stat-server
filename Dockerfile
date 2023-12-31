# 使用golang作为基础镜像
FROM golang:1.21

ENV BROKER_IP=127.0.0.1
ENV PORT=1883
ENV PROFILE=dev
ENV CRON="0 0 1 * * *"
# 将源代码复制到容器中
COPY . /app/

# 设置工作目录
WORKDIR /app/cmd

# 构建应用程序
RUN go build -o app

# 设置容器启动命令
CMD ["./app"]