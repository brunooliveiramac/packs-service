FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o packs-service ./cmd/packs-service

FROM scratch
COPY --from=builder /app/packs-service /packs-service
EXPOSE 8082
ENTRYPOINT ["/packs-service"]


