FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o build/server cmd/server/server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql
WORKDIR /app
COPY --from=builder /app/build/server .
EXPOSE 8080
CMD [ "./server" ]
