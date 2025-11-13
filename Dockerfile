FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN go env -w GOEXPERIMENT=czero,unixdefn && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /uptimer ./cmd

FROM scratch
COPY --from=builder /uptimer /uptimer
ENTRYPOINT ["/uptimer"]
