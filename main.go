package main

import (
	"fmt"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"os"
)

func main() {
	edlContent, err := os.ReadFile("example.edl")
	if err != nil {
		panic(err)
	}

	edlList, err := edl.Parse(string(edlContent))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", edlList)
}
