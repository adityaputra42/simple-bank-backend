# Build stage
FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
 

# Run Stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY start.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]