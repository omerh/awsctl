FROM golang:1.11.2 AS Builder

WORKDIR /go/src/app
COPY main.go /go/src/app

RUN go get github.com/aws/aws-sdk-go

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags '-w -extldflags "-static"'

FROM alpine:3.8 AS Runner

RUN apk add --update ca-certificates
COPY  --from=Builder /go/src/app/app /usr/local/bin/app

CMD [ "app" ]
