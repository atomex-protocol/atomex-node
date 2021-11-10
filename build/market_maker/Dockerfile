# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.16-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV AP_ENV=production

RUN apk --no-cache add ca-certificates
RUN apk add --update git

RUN mkdir -p $GOPATH/src/github.com/atomex-protocol/watch_tower/

COPY ./go.* $GOPATH/src/github.com/atomex-protocol/watch_tower/
WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/
RUN go mod download

COPY cmd/market_maker cmd/market_maker
COPY internal internal

WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/cmd/market_maker/
RUN go build -a -installsuffix cgo -o /go/bin/market_maker .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/market_maker
ENV GOPATH=$GOPATH
ENV AP_ENV=${AP_ENV}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/market_maker /go/bin/market_maker
COPY assets/abi/ abi/
COPY configs/${AP_ENV} configs/

ENTRYPOINT ["/go/bin/market_maker"]