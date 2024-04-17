FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o palworld_exporter

FROM alpine:3.19

COPY --from=builder /app/palworld_exporter /bin/palworld_exporter
EXPOSE 18212
ENTRYPOINT ["/bin/palworld_exporter"]