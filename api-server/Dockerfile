FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /api-server /api-server

EXPOSE 8080

CMD ["/api-server"]
