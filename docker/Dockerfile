# Multi-stage build for PayChain
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /paychain ./cmd

# Runtime image
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=builder /paychain /paychain
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/paychain"]

