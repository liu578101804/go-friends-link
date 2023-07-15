FROM golang:1.20-alpine3.18 as builder

RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++

# 启用go module
ENV TZ=Asia/Shanghai \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

RUN mkdir /src
WORKDIR /src

ADD . /src

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy
RUN GOPROXY=https://goproxy.cn go build -o app .  && chmod +x app


# 第二步部生成镜像
FROM alpine:3.18

# 指定运行时环境变量
ENV GIN_MODE=release \
    TZ=Asia/Shanghai

RUN mkdir /app
WORKDIR /app

COPY --from=builder /src/app    /app/app
COPY --from=builder /src/ui    /app/ui

CMD ["/app/app"]
