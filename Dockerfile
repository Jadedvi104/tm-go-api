# 1. Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the binary named "gateway"
RUN go build -o gateway .

# 2. Run Stage (Small Image)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .
# Install certificates (needed for external API calls)
RUN apk --no-cache add ca-certificates
EXPOSE 80
CMD ["./gateway"]