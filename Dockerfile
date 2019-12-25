# bump: golang /GOLANG_VERSION=([\d.]+)/ docker:golang|^1
ARG GOLANG_VERSION=1.13.5
# bump: youtube-dl /YDL_VERSION=([\d.]+)/ https://github.com/ytdl-org/youtube-dl.git|/^\d/|sort
ARG YDL_VERSION=2019.12.25

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
