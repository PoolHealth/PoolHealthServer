FROM golang:1.23

WORKDIR /tmp

RUN apt-get update
# dnsutils is needed to have dig installed to create cluster file
RUN apt-get install -y --no-install-recommends ca-certificates dnsutils


ARG GOPROXY
ENV \
  GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/github.com/PoolHealth/PoolHealthServer/
ADD go.mod go.sum /go/src/github.com/PoolHealth/PoolHealthServer/
RUN go mod download -x

ADD . .

ARG VERSION
RUN go build -v -ldflags="-w -s -X main.version=${VERSION}" -o /bin/server cmd/server/*.go

CMD /bin/server