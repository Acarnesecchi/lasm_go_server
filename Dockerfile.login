FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o login-server .

FROM alpine:latest

WORKDIR /app

ENV APP=login-server

COPY --from=builder /app/login-server .

COPY login.html .

ENTRYPOINT ["./login-server"]
