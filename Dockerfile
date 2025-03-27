# docker build  -t dmc-task .
# docker run -it --rm -p 8080:8080 dmc-task

ARG base_image=golang:1.24-alpine
ARG runner_image=debian:stable-slim

# Build phase
FROM ${base_image} AS builder

ARG go_proxy=https://proxy.golang.org,direct

# 创建配置目录并设置工作目录
WORKDIR /app

# 将整个代码复制到 /app 目录中
COPY . .

ENV GOPROXY=${go_proxy}

# 设置 Go 环境变量
RUN go env -w GO111MODULE=on \
    && go env -w GOOS=linux \
    && go env -w GOARCH=amd64 \
    && go env \
    && go mod tidy

# 编译 Go 程序
RUN go build -o build/dmc-task ./cmd/dmctask
RUN go build -o build/app ./cmd/app

# 查看编译的文件
RUN ls -l /app/build

# Final image phase
FROM ${runner_image}

# 安装证书包，确保 SSL/TLS 证书有效
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Shanghai

# 设置工作目录为 /app
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件到 /app 目录
COPY --from=builder /app/build/dmc-task ./dmc-task
COPY --from=builder /app/build/app ./app

# 例如将本地的配置文件复制到容器的 /app/conf 目录
#COPY --from=builder /app/cmd/dmctask/conf/conf.yaml /app/conf.yaml
#COPY --from=builder /app/cmd/app/conf/conf.yaml /app/biz_conf.yaml

# 查看容器的文件结构，确保配置文件和二进制文件存在
RUN ls -la *

# 暴露 http-7888 和 grpc-7889端口
EXPOSE 7888
EXPOSE 7889

# 设置容器启动时的命令
# 通过 -c 参数指定配置文件的路径
ENTRYPOINT ["/app/dmc-task","server", "--cfg", "/app/conf.yaml"]