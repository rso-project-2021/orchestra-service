# Build stage
FROM golang:1.17.3-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY config.json .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration
RUN apk add --no-cache jq
RUN chmod +x ./start.sh
RUN chmod +x ./wait-for.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
