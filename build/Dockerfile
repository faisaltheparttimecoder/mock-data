############################
# STEP 0 get dependencies
############################
FROM golang:1.16.5 AS dependencies
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
RUN go mod download
############################
# STEP 1 build executable binary
############################
FROM dependencies AS builder
##
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src
COPY . .
##
RUN go build -o /bin/mock .
############################
# STEP 2 build a small image
############################
FROM alpine:3.12

RUN addgroup -S mock; \
    adduser -S mock -G mock -D -u 10000 -h /home/mock -s /bin/nologin; \
    chown -R mock:mock /home/mock
COPY --from=builder /bin/mock /usr/bin/mock
VOLUME [ "/home/mock" ]

USER 10000
WORKDIR /home/mock

ENTRYPOINT [ "mock" ]
CMD ["--help"]
