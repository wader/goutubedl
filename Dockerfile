# bump: golang /GOLANG_VERSION=([\d.]+)/ docker:golang|^1
# bump: golang link "Release notes" https://golang.org/doc/devel/release.html
ARG GOLANG_VERSION=1.15.6
# bump: youtube-dl /YDL_VERSION=([\d.]+)/ https://github.com/ytdl-org/youtube-dl.git|/^\d/|sort
# bump: youtube-dl link "Release notes" https://github.com/ytdl-org/youtube-dl/releases/tag/$LATEST
ARG YDL_VERSION=2021.01.08

FROM golang:$GOLANG_VERSION
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
