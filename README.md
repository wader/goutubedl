## goutubedl

Go wrapper for [youtube-dl](https://github.com/ytdl-org/youtube-dl).

Online [API documentation](https://godoc.org/github.com/wader/goutubedl) can be found at godoc.org.

See [youtube-dl documentation](https://github.com/ytdl-org/youtube-dl#do-i-need-any-other-programs)
for what is recommended to install in addition to youtube-dl.

### Usage

```go
result, err := goutubedl.New(context.Background(), URL, goutubedl.Options{})
downloadResult, err := result.Download(context.Background(), FormatID)
io.Copy(ioutil.Discard, downloadResult)
downloadResult.Close()
```

See [goutubedl cmd tool](cmd/goutubedl/main.go) or [ydls](https://github.com/wader/ydls)
for usage examples.
