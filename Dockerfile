FROM golang:1.18 AS builder

WORKDIR /go/src/websocket-loadtester

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /loadtester .

FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /loadtester /loadtester

ENTRYPOINT ["./loadtester"]

