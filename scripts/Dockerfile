##################
FROM golang:1.19.3-alpine as builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 更新一下仓库
# RUN apk update

RUN apk --no-cache add git ca-certificates gcc g++
# ENV GO111MODULE=on \
#     GOPROXY=https://goproxy.cn,direct
# WORKDIR /go/src/github.com/go-sonic/sonic
# COPY ./scripts/go.mod go.mod
# COPY ./scripts/go.sum go.sum
# RUN go mod download    

COPY . /go/src/github.com/go-sonic/sonic/
WORKDIR /go/src/github.com/go-sonic/sonic

ARG BUILD_COMMIT
ARG BUILD_TIME
ARG SONIC_VERSION

ENV GO111MODULE=off \
    GOPATH=/go
    
RUN CGO_ENABLED=1 GOOS=linux && \
go build -o sonic -ldflags="-s -w -X github.com/go-sonic/sonic/consts.SonicVersion=${SONIC_VERSION} -X github.com/go-sonic/sonic/consts.BuildCommit=${BUILD_COMMIT} -X github.com/go-sonic/sonic/consts.BuildTime=${BUILD_TIME}" -trimpath .

RUN mkdir -p /app/conf && \
    mkdir /app/resources && \
    cp -r /go/src/github.com/go-sonic/sonic/sonic /app/ && \
    cp -r /go/src/github.com/go-sonic/sonic/conf/config.yaml /app/conf/config.yaml && \
    cp -r /go/src/github.com/go-sonic/sonic/resources /app/ && \
    cp /go/src/github.com/go-sonic/sonic/scripts/docker_init.sh /app/

##################
FROM alpine:3.20 as prod
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 更新一下仓库
# RUN apk update

COPY --from=builder /app/ /app/
# 检查拷贝内容
COPY . /go/src/github.com/go-sonic/sonic/ 

RUN  apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

VOLUME /sonic
EXPOSE 8080

WORKDIR /sonic
CMD /app/docker_init.sh && /app/sonic -config /sonic/conf/config.yaml