FROM golang:1.22.2-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Enabling cgo (CGO_ENABLED=1) often leads to larger binaries.
# This increase in size is attributed to the overhead introduced by bridging Go and C code
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM gcr.io/distroless/base-debian12

USER nonroot

COPY --from=builder /app/app /

CMD ["/app"]
