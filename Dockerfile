FROM golang:1.11.1-alpine as builder

RUN apk --no-cache add git

COPY . $GOPATH/src/github.com/raqbit/catbot/
WORKDIR $GOPATH/src/github.com/raqbit/catbot/

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/catbot

FROM alpine:3.8

RUN apk --no-cache add ca-certificates tini net-tools

COPY --from=builder /go/bin/catbot /go/bin/catbot

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/bin/catbot"]
