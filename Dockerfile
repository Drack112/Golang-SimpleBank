FROM golang:1.18.8-alpine as builder

ENV GO111MODULE=on

RUN apk update && apk add --no-cache bash

WORKDIR /go/src/app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod /go/src/app/
COPY go.sum /go/src/app/

RUN go mod download

COPY . /go/src/app/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /go/src/app

COPY --from=builder /go/src/app/main .
COPY --from=builder /go/src/app/.env .

EXPOSE 8080

CMD ["./main"]
