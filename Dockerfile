FROM golang:1.17-alpine

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

COPY ./ /go/src/github.com/ti/mdrest
WORKDIR /go/src/github.com/ti/mdrest/mdrest
RUN go install -ldflags '-s -w'

FROM alpine:3.14
COPY --from=0 /go/bin/mdrest /app/mdrest
COPY --from=0 /go/src/github.com/ti/mdrest/mdrest/config.json /app/
CMD  ["/app/mdrest", "--help"]