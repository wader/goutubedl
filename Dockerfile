ARG GO_VERSION=1.12
ARG YDL_VERSION=2019.07.16

FROM golang:$GO_VERSION
ARG YDL_VERSION

RUN \
  curl -L -o /usr/local/bin/youtube-dl https://yt-dl.org/downloads/$YDL_VERSION/youtube-dl && \
  chmod a+x /usr/local/bin/youtube-dl

WORKDIR /src
COPY go.* *.go ./
COPY cmd cmd
RUN \
  go mod download && \
  go build ./cmd/goutubedl && \
  go test -v -race -cover
