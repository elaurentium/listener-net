FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache make
RUN make all

FROM alpine:latest
RUN apk --no-cache add ca-certificates mysql-client
WORKDIR /app
COPY --from=builder /app/build .
EXPOSE 8080
CMD [ "./server" ]
