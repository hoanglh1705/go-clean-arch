FROM golang:1.23-bullseye

ENV GO111MODULE on
ENV APP_ENV development

# RUN apk add bash ca-certificates curl git gcc g++ libc-dev autoconf automake libtool make librdkafka-dev pkgconf
RUN mkdir -p /go/src/gitlab.com/hoanglh7/tele-money

WORKDIR /go/src/gitlab.com/hoanglh7/tele-money

ADD . .

# RUN git config --global url."https://developer:token@gitlab.com/".insteadOf "https://gitlab.com/"

RUN go install github.com/githubnemo/CompileDaemon@latest
# RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8193
ENTRYPOINT CompileDaemon -build="go build -o ./cmd/server/build/tele-money ./cmd/server/main.go" -command="./cmd/server/build/tele-money"
