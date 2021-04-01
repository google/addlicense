FROM golang:1-alpine

RUN apk add --no-cache --upgrade git openssh-client ca-certificates

COPY . /go/src/app
WORKDIR /go/src/app

RUN go build -o addlicense ./...

ENTRYPOINT ["addlicense"]
