# Многоэтапная сборка

# Stage 1: Сборка приложения
FROM golang:1.19-alpine3.16 as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main ./cmd

# Stage 2: Старт приложения
FROM alpine:3.16
WORKDIR /
COPY --from=builder /main /main
COPY config.yaml .
COPY wait-for.sh .
COPY ./migrations ./migrations
EXPOSE 8181
CMD ["/main"]