FROM golang:1.11.1-alpine as builder

RUN apk --no-cache add git

COPY . /go/src/github.com/raqbit/catbot/
WORKDIR /go/src/github.com/raqbit/catbot/

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/catbot

FROM alpine:3.8

RUN apk add --no-cache --update wget && \
    apk add --no-cache ca-certificates tini net-tools && \
    wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x wait-for-it.sh && \
    apk del wget

COPY --from=builder /go/bin/catbot /catbot
COPY --from=builder /go/src/github.com/raqbit/catbot/migrations /migrations

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/catbot"]
