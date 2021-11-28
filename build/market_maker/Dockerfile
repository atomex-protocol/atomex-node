# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.16-alpine as builder

ENV AP_ENV=production

RUN apk --no-cache add ca-certificates g++

RUN mkdir -p $GOPATH/src/github.com/atomex-protocol/watch_tower/

COPY ./go.* $GOPATH/src/github.com/atomex-protocol/watch_tower/
WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/
RUN go mod download

COPY cmd/market_maker cmd/market_maker
COPY internal internal

WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/cmd/market_maker/

# enable CGO for cryptography
RUN CGO_ENABLED=1 GO111MODULE=on GOOS=linux go build -ldflags "-linkmode external -extldflags '-static' -s -w"  -a -o /go/bin/market_maker .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/market_maker
ENV AP_ENV=$AP_ENV

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/market_maker /go/bin/market_maker
COPY assets/abi abi
COPY configs configs

ENTRYPOINT ["/go/bin/market_maker", "-c", "/app/market_maker/configs"]