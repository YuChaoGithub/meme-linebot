# Build stage.
FROM golang:1.14.4 AS builder
WORKDIR /src
COPY . .
COPY /ui /bin/ui
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/server .

# Final stage.
FROM alpine:3.12.0
WORKDIR /app

COPY --from=builder /bin/ .
CMD ["./server"]