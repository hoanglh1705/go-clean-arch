FROM golang:1.23-alpine as builder

# Install required dependencies
RUN apk --no-cache add make git bash openssh openssl-dev

# Update CA Certificates and install timezone data
RUN apk add --update --no-cache git build-base

# Update all to latest
RUN apk update

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy

# Copy the source code to Work DIR
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-X 'main.Version=`git rev-parse --short HEAD`'" -o /go/bin/gocleanarch ./cmd/server

FROM alpine:latest
# COPY ./.env ./.env
WORKDIR /go/bin

COPY --from=builder /go/bin/gocleanarch /go/bin/gocleanarch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy all swagger
RUN mkdir swaggerui
COPY ./swaggerui ./swaggerui

ENTRYPOINT ["./gocleanarch", "-c", "/config/.env"]
