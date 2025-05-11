FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o app ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
