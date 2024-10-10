FROM golang:1.20 as builder

WORKDIR /build

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn/,direct

COPY . .
RUN go mod download
RUN go mod tidy

RUN bash ./scripts/build-all.sh
# 多阶构建。

FROM alpine:3.19

ENV TZ Asia/Shanghai

WORKDIR /project
ENV PROJECT_MODE prod

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder  /build/output .
COPY --from=builder  /build/sentinel.yml .
