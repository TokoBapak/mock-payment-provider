FROM golang:1.20-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o mock-payment-provider .

FROM debian:bullseye AS runtime
RUN apt-get update && apt-get upgrade -y && apt-get install -y curl ca-certificates sqlite3
WORKDIR /app
COPY --from=builder /app/mock-payment-provider .
ENV HTTP_PORT=3000
EXPOSE ${HTTP_PORT}
CMD [ "/app/mock-payment-provider" ]
