FROM golang:1.13.0 AS Builder

WORKDIR /go/src/app
COPY . /go/src/app

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags '-w -extldflags "-static"' -o awsctl

FROM alpine:3.8 AS Runner

RUN apk add --update ca-certificates
COPY  --from=Builder /go/src/app/awsctl /usr/local/bin/awsctl

CMD [ "awsctl" ]
