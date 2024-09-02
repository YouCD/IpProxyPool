FROM golang:1.22.3-alpine AS builder
WORKDIR /builder
# 设置环境变量, 指定编码
ENV CGO_ENABLED=0 \
    GOPATH=/root/gopath \
    GOPROXY=https://goproxy.cn,direct \
    GO111MODULE='on' \
    GIT_TERMINAL_PROMPT=1 \
    LANG="en_US.UTF-8"

COPY . .
RUN go mod tidy && go build -o app .

FROM backplane/upx:latest AS upx
WORKDIR /app
COPY --from=builder /builder/app /app/app
RUN upx /app/app

FROM alpine

# 指定时区
ENV TIMEZONE=Asia/Shanghai

# 指定工作目录
WORKDIR /app

# 执行的命令
RUN set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk upgrade \
    && apk update \
    && apk add --no-cache ca-certificates upx --no-progress bash tzdata busybox-extras \
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo ${TIMEZONE} > /etc/timezone \
    && rm -rf /var/cache/apk/*

COPY --from=upx /app/app /app/IpProxyPool
COPY ./conf/config.yaml /app/conf/config.yaml

# 映射一个端口
EXPOSE 3000

ENTRYPOINT ["/app/IpProxyPool"]
