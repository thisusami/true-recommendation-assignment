FROM golang:1.26.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./src

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]