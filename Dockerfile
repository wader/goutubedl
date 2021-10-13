# bump: golang /GOLANG_VERSION=([\d.]+)/ docker:golang|^1
# bump: golang link "Release notes" https://golang.org/doc/devel/release.html
ARG GOLANG_VERSION=1.17.2
# bump: youtube-dl /YDL_VERSION=([\d.]+)/ https://github.com/ytdl-org/youtube-dl.git|/^\d/|sort
# bump: youtube-dl link "Release notes" https://github.com/ytdl-org/youtube-dl/releases/tag/$LATEST
ARG YDL_VERSION=2021.06.06

FROM golang:$GOLANG_VERSION AS base
ARG YDL_VERSION

RUN \
  apt-get update -q && \
  apt-get install -y -q python-is-python3 && \
  curl -L -o /usr/local/bin/youtube-dl https://yt-dl.org/downloads/$YDL_VERSION/youtube-dl && \
  chmod a+x /usr/local/bin/youtube-dl

FROM base AS dev

FROM base
WORKDIR /src
COPY go.* *.go ./
COPY cmd cmd
RUN \
  go mod download && \
  go build ./cmd/goutubedl && \
  go test -v -race -cover
