package main

import "github.com/VyacheArt/edl-to-youtube/converter"

func main() {
	a := converter.NewApplication()
	if err := a.Run(); err != nil {
		panic(err)
	}
}
