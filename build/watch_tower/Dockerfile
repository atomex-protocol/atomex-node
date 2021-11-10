# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.16-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV ATOMEX_PROTOCOL_ENV=production
ARG config_file

RUN apk --no-cache add ca-certificates
RUN apk add --update git

RUN mkdir -p $GOPATH/src/github.com/atomex-protocol/watch_tower/

COPY ./go.* $GOPATH/src/github.com/atomex-protocol/watch_tower/
WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/
RUN go mod download

COPY cmd/watch_tower cmd/watch_tower
COPY internal internal

WORKDIR $GOPATH/src/github.com/atomex-protocol/watch_tower/cmd/watch_tower/
RUN go build -a -installsuffix cgo -o /go/bin/watch_tower .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/watch_tower
ENV GOPATH=$GOPATH
ENV ATOMEX_PROTOCOL_ENV=$ATOMEX_PROTOCOL_ENV
ARG config_file

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/watch_tower /go/bin/watch_tower
COPY cmd/watch_tower/$config_file config.yml
COPY assets/abi/ abi/
COPY configs/${ATOMEX_PROTOCOL_ENV} configs/

ENTRYPOINT ["/go/bin/watch_tower"]