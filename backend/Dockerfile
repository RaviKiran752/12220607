# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o urlshortener

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/urlshortener .
EXPOSE 3001
CMD ["./urlshortener"] 