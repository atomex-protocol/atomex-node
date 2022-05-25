# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.18-alpine as builder

ENV AP_ENV=production

RUN apk --no-cache add ca-certificates

RUN mkdir -p $GOPATH/src/github.com/atomex-protocol/atomex-node/

COPY ./go.* $GOPATH/src/github.com/atomex-protocol/atomex-node/
WORKDIR $GOPATH/src/github.com/atomex-protocol/atomex-node/
RUN go mod download

COPY cmd/watch_tower cmd/watch_tower
COPY internal internal

WORKDIR $GOPATH/src/github.com/atomex-protocol/atomex-node/cmd/watch_tower/
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -a -o /go/bin/watch_tower .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/watch_tower
ENV AP_ENV=$AP_ENV

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/watch_tower /go/bin/watch_tower
COPY assets/abi abi
COPY configs configs

ENTRYPOINT ["/go/bin/watch_tower", "-c", "/app/watch_tower/configs"]