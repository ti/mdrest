FROM golang:1.13

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

COPY ./ /go/src/github.com/ti/mdrest
WORKDIR /go/src/github.com/ti/mdrest/mdrest
RUN go install -ldflags '-s -w'

FROM scratch
COPY --from=0 /go/bin/mdrest /mdrest
COPY --from=0 /go/src/github.com/ti/mdrest/mdrest/config.json /
CMD ["/mdrest"]

