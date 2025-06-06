# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git && \
    go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY src/ ./src/
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o storj-dashboard ./src

# Final stage
FROM alpine:latest
ENV PORT=80
RUN adduser -D -h /app appuser
WORKDIR /app
COPY --from=builder /app/storj-dashboard /app/
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 80
CMD ["/app/storj-dashboard"]