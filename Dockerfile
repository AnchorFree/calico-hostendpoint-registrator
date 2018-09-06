FROM golang:1.10-alpine3.7 as builder
LABEL maintainer="v.zorin@anchorfree.com"

RUN apk add --no-cache git bash
COPY cmd /go/src/github.com/anchorfree/hostendpoint/cmd
COPY Gopkg.toml /go/src/github.com/anchorfree/hostendpoint/
COPY Gopkg.lock /go/src/github.com/anchorfree/hostendpoint/

RUN cd /go && go get -u github.com/golang/dep/cmd/dep
RUN cd /go/src/github.com/anchorfree/hostendpoint/ && dep ensure
RUN cd /go && go build github.com/anchorfree/hostendpoint/cmd/hostendpoint

FROM alpine:3.7
LABEL maintainer="v.zorin@anchorfree.com"

COPY --from=builder /go/hostendpoint /usr/local/bin/hostendpoint

ENTRYPOINT ["/usr/local/bin/hostendpoint"]
