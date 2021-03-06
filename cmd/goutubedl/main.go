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
	log.SetFlags(0)
	flag.Parse()

	optType, _ := goutubedl.TypeFromString[*typeFlag]
	result, err := goutubedl.New(context.Background(), flag.Arg(0), goutubedl.Options{Type: optType})
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
	io.Copy(f, dr)
	dr.Close()
}
