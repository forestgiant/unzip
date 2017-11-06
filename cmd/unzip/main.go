package main

import (
	"fmt"
	"log"
	"os"

	"github.com/forestgiant/unzip"
)

func main() {
	// Read command arguments
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Must supply source and destination as arguments")
	}

	src := args[0]
	dest := args[1]

	// Try to download the src
	fmt.Printf("Downloading: %s\n", args[0])
	f, err := unzip.DownloadFile(src, "")
	if err == nil {
		// If there wasn't an error then the src was a url
		// so set the downloaded filepath to the src
		src = f
		defer os.Remove(f)
	}

	fmt.Printf("Download Finished. Starting unzip.\n")

	_, err = unzip.Unzip(src, dest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully unzipped: %s into: %s\n", args[0], dest)
}
