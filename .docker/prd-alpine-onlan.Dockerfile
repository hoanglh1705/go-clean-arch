FROM k3d-pmregistry.localhost:12346/pm-gomod-1.21-alpine-3.18-pst:1.21.6 AS builder

ARG GITHUB_ACCESS_KEY
ARG GITLAB_TOKEN
ENV GO111MODULE=on
ENV APP_ENV production
ENV GOPRIVATE=gitlab.com/pmtrade/*
ENV GOPROXY=https://proxy.golang.org

# RUN apk add bash ca-certificates curl git gcc g++ libc-dev make

RUN mkdir -p /go/src/gitlab.com/pmtrade/pm-stock-trader
WORKDIR /go/src/gitlab.com/pmtrade/pm-stock-trader

ADD . .

RUN git config --global url."https://deployer:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/"

RUN go build -o build/pst ./cmd/server/main.go

FROM k3d-pmregistry.localhost:12346/pm-alpine-3.18-runner:3.18

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
