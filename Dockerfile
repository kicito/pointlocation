FROM golang:1.11.2-alpine3.8

RUN apk add --update curl git mercurial && \
    rm -rf /var/cache/apk/*

#install dep
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/kicito/assignment-geo-2/
COPY . .

RUN dep ensure -vendor-only

RUN apk del curl git mercurial

ENTRYPOINT ["go", "run", "main.go"]