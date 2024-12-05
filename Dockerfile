# Build stage
FROM golang AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o PM-BasicApi

# Production stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/PM-BasicApi /app/PM-BasicApi
ENTRYPOINT ["/app/PM-BasicApi"]