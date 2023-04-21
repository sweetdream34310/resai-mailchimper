############################
# STEP 1 build executable binary
############################

FROM golang:1.15-alpine AS builder

RUN apk --update --no-cache add \
    openssl \
    git \
    tzdata \
    curl \
    make \
    ca-certificates \
    && update-ca-certificates \
    && go get -v github.com/mercari/gaurun | true

WORKDIR /go/src/github.com/mercari/gaurun


RUN GO111MODULE=on go build -o bin/gaurun cmd/gaurun/gaurun.go
RUN GO111MODULE=on go build -o bin/gaurun_recover cmd/gaurun_recover/gaurun_recover.go

############################
# STEP 2 build a small image
############################
FROM alpine:3.14

COPY --from=builder /go/src/github.com/mercari/gaurun/bin/gaurun /gaurun

RUN apk --update --no-cache add \
        openssl \
        tzdata \
        ca-certificates \
        && update-ca-certificates
COPY ./ ./
CMD ["/gaurun", "-c", "gaurun.toml"]
