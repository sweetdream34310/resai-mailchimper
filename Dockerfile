############################
# STEP 1 build executable binary
############################
FROM golang:1.17.3-alpine AS builder

ARG ENV

RUN apk --update --no-cache add \
    openssl \
    git \
    curl \
    tzdata \
    ca-certificates \
    && update-ca-certificates

COPY . /go/src/github.com/cloudsrc/api.awaymail.v1.go

WORKDIR /go/src/github.com/cloudsrc/api.awaymail.v1.go

RUN mkdir -p public

RUN go mod vendor && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o am_srv

############################
# STEP 2 build a small image
############################
FROM scratch

ARG ENV
ENV GOLANG_ENV=${ENV}

COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/config /config
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/am_srv /am_srv
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/serviceAccount-upload.json /serviceAccount-upload.json
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/awaymail-84795cd1daa6.json /awaymail-84795cd1daa6.json
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/public /public
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/mail/email.txt /resources/email.txt
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/mail/email.html /resources/email.html
COPY --from=builder /go/src/github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/mail/email.amp.html /resources/email.amp.html

CMD ["/am_srv"]
