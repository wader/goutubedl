# bump: golang /GOLANG_VERSION=([\d.]+)/ docker:golang|^1
# bump: golang link "Release notes" https://golang.org/doc/devel/release.html
ARG GOLANG_VERSION=1.22.0
# bump: yt-dlp /YT_DLP=([\d.-]+)/ https://github.com/yt-dlp/yt-dlp.git|/^\d/|sort
# bump: yt-dlp link "Release notes" https://github.com/yt-dlp/yt-dlp/releases/tag/$LATEST
ARG YT_DLP=2023.12.30

FROM golang:$GOLANG_VERSION AS base
ARG YT_DLP

RUN \
  apt-get update -q && \
  apt-get install -y -q python-is-python3 && \
  curl -L https://github.com/yt-dlp/yt-dlp/releases/download/$YT_DLP/yt-dlp -o /usr/local/bin/yt-dlp && \
  chmod a+x /usr/local/bin/yt-dlp && \
  apt-get install -y ffmpeg

FROM base AS dev

FROM base
WORKDIR /src
COPY go.* *.go ./
COPY cmd cmd
RUN \
  go mod download && \
  go build ./cmd/goutubedl && \
  go test -v -race -cover
