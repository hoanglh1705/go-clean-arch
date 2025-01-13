# FROM golang:1.21-alpine3.18 AS builder
FROM hoangnguyen1247/golang-alpine-build-tools:1.22-alpine AS builder

ARG APP_BRANCH
ARG GITHUB_ACCESS_KEY
ARG GITLAB_TOKEN

ENV GO111MODULE=on
ENV APP_ENV production
ENV GOPRIVATE=gitlab.com/pmtrade/*
ENV GOPROXY=https://proxy.golang.org

# RUN apk add bash ca-certificates curl git gcc g++ libc-dev make

RUN git config --global url."https://deployer:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/"

RUN mkdir -p /go/src/gitlab.com/pmtrade
WORKDIR /go/src/gitlab.com/pmtrade

RUN git clone --single-branch --branch master https://gitlab.com/pmtrade/pm-helper-libs.git
RUN git clone --single-branch --branch $APP_BRANCH https://gitlab.com/pmtrade/pm-identity-management.git
RUN git clone --single-branch --branch $APP_BRANCH https://gitlab.com/pmtrade/pm-notification-hub.git
RUN git clone --single-branch --branch $APP_BRANCH https://gitlab.com/pmtrade/pm-stock-agent.git

RUN mkdir -p /go/src/gitlab.com/pmtrade/pm-stock-trader
WORKDIR /go/src/gitlab.com/pmtrade/pm-stock-trader

ADD . .

RUN go build -o build/pst ./cmd/server/main.go

# FROM alpine:3.18
FROM hoangnguyen1247/alpine-runner:3.18

# RUN apk add ca-certificates tzdata

RUN mkdir /usr/local/lib/pst
WORKDIR /usr/local/lib/pst

ENV GO111MODULE=on
ENV APP_ENV production

COPY --from=builder /go/src/gitlab.com/pmtrade/pm-stock-trader/build/pst /usr/local/lib/pst/pst
COPY --from=builder /go/src/gitlab.com/pmtrade/pm-stock-trader/cmd/server/statik/statik.go /usr/local/lib/pst/statik.go

# EXPOSE 10003 10004
EXPOSE 10000 80
CMD ["/usr/local/lib/pst/pst"]
