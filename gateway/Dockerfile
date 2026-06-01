FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gateway main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 80
CMD ["./gateway"]