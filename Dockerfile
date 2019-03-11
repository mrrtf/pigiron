FROM golang:alpine as builder

RUN apk add git

WORKDIR /go/src/github.com/aphecetche/

COPY . pigiron

RUN cd /go/src/github.com/aphecetche && go get ./...

RUN cd /go/src/github.com/aphecetche && go install ./...

FROM alpine:3.9

COPY --from=builder /go/bin/mch-mapping-api /mch-mapping-api

CMD ["/mch-mapping-api"]



