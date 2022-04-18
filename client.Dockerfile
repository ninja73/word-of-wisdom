FROM golang:1.18-alpine AS builder

WORKDIR /build

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN pwd
RUN go build -ldflags="-s -w" -o ./client ./cmd/client/client.go

FROM scratch

COPY --from=builder ["/build/client", "/build/client-config.toml", "/"]

ENTRYPOINT ["/client"]