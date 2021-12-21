## goutubedl

Go wrapper for [youtube-dl](https://github.com/ytdl-org/youtube-dl) and [yt-dlp](https://github.com/yt-dlp/yt-dlp), currently tested and
developed using yt-dlp.
API documentation can be found at [godoc.org](https://pkg.go.dev/github.com/wader/goutubedl?tab=doc).

See [yt-dlp documentation](https://github.com/yt-dlp/yt-dlp) for how to
install and what is recommended to install in addition to youtube-dl.

goutubedl default uses `PATH` to find youtube-dl but it can be configured with the `goutubedl.Path`
variable. Default is currently `youtube-dl` for backwards compability. If your using yt-dlp you
probably want to set it to `yt-dlp`.

Due to the nature and frequent updates of youtube-dl only the latest version
is tested. But it seems to work well with older versions also.

### Usage

From [cmd/example/main.go](cmd/example/main.go)
```go
package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/wader/goutubedl"
)

func main() {
	result, err := goutubedl.New(context.Background(), "https://www.youtube.com/watch?v=jgVhBThJdXc", goutubedl.Options{})
	if err != nil {
		log.Fatal(err)
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Fatal(err)
	}
	defer downloadResult.Close()
	f, err := os.Create("output")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, downloadResult)
}

```

See [goutubedl cmd tool](cmd/goutubedl/main.go) or [ydls](https://github.com/wader/ydls)
for usage examples.

### Development

```sh
docker build -t goutubedl-dev .
docker run --rm -ti -v "$PWD:$PWD" -w "$PWD" goutubedl-dev
go test -v -race -cover
```
