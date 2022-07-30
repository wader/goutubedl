package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/wader/goutubedl"
)

var dumpFlag = flag.Bool("J", false, "Dump JSON")
var typeFlag = flag.String("t", "any", "Type")

func main() {
	goutubedl.Path = "yt-dlp"

	log.SetFlags(0)
	flag.Parse()

	optType := goutubedl.TypeFromString[*typeFlag]
	result, err := goutubedl.New(
		context.Background(),
		flag.Arg(0),
		goutubedl.Options{Type: optType, DebugLog: log.Default()},
	)
	if err != nil {
		log.Fatal(err)
	}

	if *dumpFlag {
		json.NewEncoder(os.Stdout).Encode(result.Info)
		return
	}

	filter := flag.Arg(1)
	if filter == "" {
		filter = "best"
	}

	dr, err := result.Download(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filter)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(f, dr); err != nil {
		log.Fatal(err)
	}
	dr.Close()
}
