FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o bin/batch ./cmd

FROM alpine:latest
COPY --from=builder /build/bin /bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/batch"]