FROM golang:1.20 as builder
RUN mkdir /src
# 启用go module
ENV TZ=Asia/Shanghai \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
ADD . /src
WORKDIR /src
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy
RUN GOPROXY=https://goproxy.cn go build -o app .  && chmod +x app


# 第二步部生成镜像
FROM golang:1.20
#ENV ZONEINFO=/app/zoneinfo.zip
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/app    /app
# 指定运行时环境变量
ENV GIN_MODE=release \
    TZ=Asia/Shanghai

CMD ["./app"]
