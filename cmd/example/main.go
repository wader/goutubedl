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
