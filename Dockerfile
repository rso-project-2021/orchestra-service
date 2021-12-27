# Build stage
FROM golang:1.17.3-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY config.json .

EXPOSE 8080
CMD [ "/app/main" ]
