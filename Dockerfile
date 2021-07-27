FROM golang:1.16.3-alpine as builder
ENV GOPROXY https://goproxy.cn
ENV GO111MODULE=on
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download
WORKDIR /app
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o mock main.go

FROM alpine
#ENV GIN_MODE=release
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/config.yml /config.yml
COPY --from=builder /app/mock /mock
CMD ["/mock"]