package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2{
		log.Fatal("usage: ingrecongnition <img_url>")
	}
	fmt.Printf("url: %s", os.Args[1])

	resp, err :== http.Get(os.Args[1])
	if err != nil {
		log.Fatalf("unable to get an image: %v", err)
	}
	defer resp.Body.Close()
}

