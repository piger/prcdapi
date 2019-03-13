FROM golang:1.12-stretch

WORKDIR /go/src/github.com/piger/prcdapi
COPY . .

RUN go get -d -v ./...
RUN go build ./cmd/prcdapi

FROM debian:stable-slim
COPY --from=0 /go/src/github.com/piger/prcdapi/prcdapi /usr/local/sbin/
ENTRYPOINT ["/usr/local/sbin/prcdapi"]
EXPOSE 30666/tcp
CMD ["-address", "0.0.0.0:30666", "/data"]
